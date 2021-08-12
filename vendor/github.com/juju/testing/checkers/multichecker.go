// Copyright 2020 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package checkers

import (
	"fmt"
	"go/ast"
	"go/parser"
	"regexp"
	"strings"
	"sync"

	gc "gopkg.in/check.v1"
)

// MultiChecker is a deep checker that by default matches for equality.
// But checks can be overriden based on path (either explicit match or regexp)
type MultiChecker struct {
	*gc.CheckerInfo
	checks      map[string]checkerWithArgs
	matchChecks []matchCheck
}

type checkerWithArgs interface {
	gc.Checker
	Args() []interface{}
}

type matchCheck interface {
	checkerWithArgs
	MatchString(string) bool
	WantTopLevel() bool
}

type multiCheck struct {
	gc.Checker
	args []interface{}
}

func (m *multiCheck) Args() []interface{} {
	return m.args
}

type regexCheck struct {
	multiCheck
	*regexp.Regexp
}

func (r *regexCheck) WantTopLevel() bool {
	return false
}

type astCheck struct {
	multiCheck
	astExpr ast.Expr
}

func (a *astCheck) WantTopLevel() bool {
	return true
}

// NewMultiChecker creates a MultiChecker which is a deep checker that by default matches for equality.
// But checks can be overriden based on path (either explicit match or regexp)
func NewMultiChecker() *MultiChecker {
	return &MultiChecker{
		CheckerInfo: &gc.CheckerInfo{Name: "MultiChecker", Params: []string{"obtained", "expected"}},
		checks:      make(map[string]checkerWithArgs),
	}
}

// Add an explict checker by path.
func (checker *MultiChecker) Add(path string, c gc.Checker, args ...interface{}) *MultiChecker {
	checker.checks[path] = &multiCheck{
		Checker: c,
		args:    args,
	}
	return checker
}

// AddRegex exception which matches path with regex.
func (checker *MultiChecker) AddRegex(pathRegex string, c gc.Checker, args ...interface{}) *MultiChecker {
	checker.matchChecks = append(checker.matchChecks, &regexCheck{
		multiCheck: multiCheck{
			Checker: c,
			args:    args,
		},
		Regexp: regexp.MustCompile("^" + pathRegex + "$"),
	})
	return checker
}

// AddExpr exception which matches path with go expression. Use _ for wildcard.
// The top level or root value must be a _ when using expression.
func (checker *MultiChecker) AddExpr(expr string, c gc.Checker, args ...interface{}) *MultiChecker {
	astExpr, err := parser.ParseExpr(expr)
	if err != nil {
		panic(err)
	}

	checker.matchChecks = append(checker.matchChecks, &astCheck{
		multiCheck: multiCheck{
			Checker: c,
			args:    args,
		},
		astExpr: astExpr,
	})
	return checker
}

// topLevel is a substitute for the top level or root object.
// We use an unlikely value to provide backwards compatability with previous deep equals
// behaviour. It is stripped out before any errors are printed.
const topLevel = "🔝"

// Check for go check Checker interface.
func (checker *MultiChecker) Check(params []interface{}, names []string) (result bool, errStr string) {
	customCheckFunc := func(path string, a1 interface{}, a2 interface{}) (useDefault bool, equal bool, err error) {
		var mc checkerWithArgs
		if c, ok := checker.checks[strings.Replace(path, topLevel, "", 1)]; ok {
			mc = c
		} else {
			for _, v := range checker.matchChecks {
				pathCopy := path
				if !v.WantTopLevel() {
					pathCopy = strings.Replace(pathCopy, topLevel, "", 1)
				}
				if v.MatchString(pathCopy) {
					mc = v
					break
				}
			}
		}
		if mc == nil {
			return true, false, nil
		}

		params := append([]interface{}{a1}, mc.Args()...)
		info := mc.Info()
		if len(params) < len(info.Params) {
			return false, false, fmt.Errorf("Wrong number of parameters for %s: want %d, got %d", info.Name, len(info.Params), len(params)+1)
		}
		// Copy since it may be mutated by Check.
		names := append([]string{}, info.Params...)

		// Trim to the expected params len.
		params = params[:len(info.Params)]

		// Perform substitution
		for i, v := range params {
			if v == ExpectedValue {
				params[i] = a2
			}
		}

		result, errStr := mc.Check(params, names)
		if result {
			return false, true, nil
		}
		path = strings.Replace(path, topLevel, "", 1)
		if path == "" {
			path = "top level"
		}
		return false, false, fmt.Errorf("mismatch at %s: %s", path, errStr)
	}
	if ok, err := DeepEqualWithCustomCheck(params[0], params[1], customCheckFunc); !ok {
		return false, err.Error()
	}
	return true, ""
}

// ExpectedValue if passed to MultiChecker.Add or MultiChecker.AddRegex, will be substituded with the expected value.
var ExpectedValue = &struct{}{}

var (
	astCache     = make(map[string]ast.Expr)
	astCacheLock = sync.Mutex{}
)

func (a *astCheck) MatchString(expr string) bool {
	expr = strings.Replace(expr, topLevel, "_", 1)
	astCacheLock.Lock()
	astExpr, ok := astCache[expr]
	astCacheLock.Unlock()
	if !ok {
		var err error
		astExpr, err = parser.ParseExpr(expr)
		if err != nil {
			panic(err)
		}
		astCacheLock.Lock()
		astCache[expr] = astExpr
		astCacheLock.Unlock()
	}

	return matchAstExpr(a.astExpr, astExpr)
}

func matchAstExpr(expected, obtained ast.Expr) bool {
	switch expr := expected.(type) {
	case *ast.IndexExpr:
		x, ok := obtained.(*ast.IndexExpr)
		if !ok {
			return false
		}
		if !matchAstExpr(expr.Index, x.Index) {
			return false
		}
		if !matchAstExpr(expr.X, x.X) {
			return false
		}
	case *ast.ParenExpr:
		x, ok := obtained.(*ast.ParenExpr)
		if !ok {
			return false
		}
		if !matchAstExpr(expr.X, x.X) {
			return false
		}
	case *ast.StarExpr:
		x, ok := obtained.(*ast.StarExpr)
		if !ok {
			return false
		}
		if !matchAstExpr(expr.X, x.X) {
			return false
		}
	case *ast.SelectorExpr:
		x, ok := obtained.(*ast.SelectorExpr)
		if !ok {
			return false
		}
		if !matchAstExpr(expr.X, x.X) {
			return false
		}
		if !matchAstExpr(expr.Sel, x.Sel) {
			return false
		}
	case *ast.Ident:
		if expr.Name == "_" {
			// Wildcard
			return true
		}
		x, ok := obtained.(*ast.Ident)
		if !ok {
			return false
		}
		if expr.Name != x.Name {
			return false
		}
	case *ast.BasicLit:
		x, ok := obtained.(*ast.BasicLit)
		if !ok {
			return false
		}
		if expr.Value != x.Value {
			return false
		}
	default:
		panic(fmt.Sprintf("unknown type %#v", expected))
	}
	return true
}
