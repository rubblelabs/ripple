package data

import (
	"fmt"
	"github.com/donovanhide/ripple/crypto"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

const (
	minOffset        int64  = -96
	maxOffset        int64  = 80
	minValue         uint64 = 1000000000000000
	maxValue         uint64 = 9999999999999999
	maxNative        uint64 = 9000000000000000000
	maxNativeNetwork uint64 = 100000000000000000
	notNative        uint64 = 0x8000000000000000
	positive         uint64 = 0x4000000000000000
	maxNativeSqrt    uint64 = 3000000000
	maxNativeDiv     uint64 = 2095475792 // MaxNative / 2^32
	tenTo14          uint64 = 100000000000000
	tenTo14m1        uint64 = tenTo14 - 1
	tenTo17          uint64 = tenTo14 * 1000
	tenTo17m1        uint64 = tenTo17 - 1
)

var (
	bigTen     = big.NewInt(10)
	bigTenTo14 = big.NewInt(0).SetUint64(tenTo14)
	bigTenTo17 = big.NewInt(0).SetUint64(tenTo17)
)

type Value struct {
	Native   bool
	Negative bool
	Num      uint64
	Offset   int64
}

type Amount struct {
	*Value
	Currency Currency
	Issuer   Account
}

func (v *Value) Debug() string {
	return fmt.Sprintf("Native: %t Negative: %t Value: %d Offset: %d", v.Native, v.Negative, v.Num, v.Offset)
}

func NewValue(native, negative bool, num uint64, offset int64) *Value {
	return &Value{
		Native:   native,
		Negative: negative,
		Num:      num,
		Offset:   offset,
	}
}

// Match fields:
// 0 = whole input
// 1 = sign
// 2 = integer portion
// 3 = whole fraction (with '.')
// 4 = fraction (without '.')
// 5 = whole exponent (with 'e')
// 6 = exponent sign
// 7 = exponent number
var valueRegex = regexp.MustCompile("([+-]?)(\\d*)(\\.(\\d*))?([eE]([+-]?)(\\d+))?")

func (v *Value) Parse(s string) error {
	var err error
	matches := valueRegex.FindStringSubmatch(s)
	if matches == nil {
		return fmt.Errorf("Invalid Number: %s", s)
	}
	if len(matches[2])+len(matches[4]) > 32 {
		return fmt.Errorf("Overlong Number: %s", s)
	}
	if matches[1] == "-" {
		v.Negative = true
	}
	if len(matches[4]) == 0 {
		if v.Num, err = strconv.ParseUint(matches[2], 10, 64); err != nil {
			return fmt.Errorf("Invalid Number: %s", s)
		}
		v.Offset = 0
	} else {
		if v.Num, err = strconv.ParseUint(matches[2]+matches[4], 10, 64); err != nil {
			return fmt.Errorf("Invalid Number: %s", s)
		}
		v.Offset = -int64(len(matches[4]))
	}
	if len(matches[5]) > 0 {
		exp, err := strconv.ParseInt(matches[7], 10, 64)
		if err != nil {
			return fmt.Errorf("Invalid Number: %s", s)
		}
		if matches[6] == "-" {
			v.Offset -= exp
		} else {
			v.Offset += exp
		}
	}
	if v.Native {
		if len(matches[3]) > 0 {
			v.Offset += 6
		}
	}
	return v.canonicalise()
}

func (v *Value) canonicalise() error {
	if v.Native {
		if v.Num == 0 {
			v.Offset = 0
			v.Negative = false
		} else {
			for v.Offset < 0 {
				v.Num /= 10
				v.Offset++
			}
			for v.Offset > 0 {
				v.Num *= 10
				v.Offset--
			}
			if v.Num > maxNative {
				return fmt.Errorf("Native amount out of range: %s", v.Debug())
			}
		}
	} else {
		if v.Num == 0 {
			v.Offset = -100
			v.Negative = false
		} else {
			for v.Num < minValue && v.Offset > minOffset {
				v.Num *= 10
				v.Offset--
			}
			for v.Num > maxValue {
				if v.Offset >= maxOffset {
					return fmt.Errorf("Value overflow: %s", v.Debug())
				}
				v.Num /= 10
				v.Offset++
			}
			if v.Offset < minOffset || v.Num < minValue {
				v.Num = 0
				v.Offset = 0
				v.Negative = false
			}
			if v.Offset > maxOffset {
				return fmt.Errorf("Value overflow: %s", v.Debug())
			}
		}
	}
	return nil
}

func (v *Value) Clone() *Value {
	return NewValue(v.Native, v.Negative, v.Num, v.Offset)
}

func (v *Value) Equals(other *Value) bool {
	return v.Native == other.Native &&
		v.Negative == other.Negative &&
		v.Num == other.Num &&
		v.Offset == other.Offset
}

func newAmount(value *Value, currency Currency, issuer Account) *Amount {
	return &Amount{
		Value:    value,
		Currency: currency,
		Issuer:   issuer,
	}
}

func NewAmount(v interface{}) (*Amount, error) {
	switch n := v.(type) {
	case int64:
		return &Amount{
			Value: NewValue(true, n < 0, abs(n), 0),
		}, nil
	case string:
		amount := &Amount{
			Value: &Value{},
		}
		var err error
		parts := strings.Split(n, "/")
		if len(parts) == 1 {
			amount.Native = true
		}
		if len(parts) > 1 && parts[1] == "XRP" {
			amount.Native = true
			if !strings.Contains(parts[0], ".") {
				parts[0] = parts[0] + "."
			}
		}
		if amount.Parse(parts[0]); err != nil {
			return nil, err
		}
		if len(parts) > 1 {
			if amount.Currency, err = NewCurrency(parts[1]); err != nil {
				return nil, err
			}
		}
		if len(parts) > 2 {
			if issuer, err := crypto.NewRippleHash(parts[2]); err != nil {
				return nil, err
			} else {
				copy(amount.Issuer[:], issuer.Payload())
			}
		}
		return amount, nil
	default:
		return nil, fmt.Errorf("Bad type: %+v", v)
	}
}

func (a *Amount) Equals(b *Amount) bool {
	return a.Value.Equals(b.Value) &&
		a.Currency == b.Currency &&
		a.Issuer == b.Issuer
}

func (a *Amount) SameValue(b *Amount) bool {
	return a.Value.Equals(b.Value)
}

func (a *Amount) Clone() *Amount {
	return newAmount(a.Value.Clone(), a.Currency, a.Issuer)
}

func (a *Amount) ZeroClone() *Amount {
	zero := &Value{Native: a.Native}
	return newAmount(zero, a.Currency, a.Issuer)
}

func (a *Amount) IsPositive() bool {
	return !a.Negative
}

func (a *Amount) Negate() *Amount {
	clone := a.Clone()
	clone.Negative = !clone.Negative
	return clone
}

func (a *Amount) Abs() *Amount {
	clone := a.Clone()
	clone.Negative = false
	return clone
}

func (a *Value) factor(b *Value) (int64, int64, int64) {
	ao, bo := a.Offset, b.Offset
	av, bv := int64(a.Num), int64(b.Num)
	if a.Negative {
		av = -av
	}
	if b.Negative {
		bv = -bv
	}
	for ; ao < bo; ao++ {
		av /= 10
	}
	for ; bo < ao; bo++ {
		bv /= 10
	}
	return av, bv, ao
}

func (a *Amount) Add(b *Amount) (*Amount, error) {
	switch {
	case b.IsZero():
		return a.Clone(), nil
	case a.IsZero():
		return newAmount(b.Value.Clone(), a.Currency, a.Issuer), nil
	case a.Native:
		return NewAmount(int64(a.Num + b.Num))
	default:
		av, bv, ao := a.factor(b.Value)
		v := NewValue(false, (av+bv) < 0, abs(av+bv), ao)
		c := newAmount(v, a.Currency, a.Issuer)
		return c, c.canonicalise()
	}
}

func (a *Amount) Subtract(b *Amount) (*Amount, error) {
	switch {
	case b.IsZero():
		return a.Clone(), nil
	case a.IsZero():
		return newAmount(b.Value.Clone(), a.Currency, a.Issuer), nil
	case a.Native:
		return NewAmount(int64(a.Num - b.Num))
	default:
		av, bv, ao := a.factor(b.Value)
		v := NewValue(false, (av-bv) < 0, abs(av-bv), ao)
		c := newAmount(v, a.Currency, a.Issuer)
		return c, c.canonicalise()
	}
}

func (a *Amount) Multiply(b *Amount) (*Amount, error) {
	if a.IsZero() || b.IsZero() {
		return a.ZeroClone(), nil
	}
	if a.Native && b.Native && a.Currency.IsNative() {
		min := min64(a.Num, b.Num)
		max := max64(a.Num, b.Num)
		if min > maxNativeSqrt || (((max >> 32) * min) > maxNativeDiv) {
			return nil, fmt.Errorf("Native value overflow: %s*%s", a.Debug(), b.Debug())
		}
		return NewAmount(int64(min * max))
	}
	av, bv := a.Num, b.Num
	ao, bo := a.Offset, b.Offset
	if a.Native {
		for ; av < minValue; av *= 10 {
			ao--
		}
	}
	if b.Native {
		for ; bv < minValue; bv *= 10 {
			bo--
		}
	}
	// Compute (numerator * denominator) / 10^14 with rounding
	// 10^16 <= result <= 10^18
	m := big.NewInt(0).SetUint64(av)
	m.Mul(m, big.NewInt(0).SetUint64(bv))
	m.Div(m, bigTenTo14)
	// 10^16 <= product <= 10^18
	if len(m.Bytes()) > 64 {
		return nil, fmt.Errorf("Multiply: %s*%s", a.Debug(), b.Debug())
	}
	v := NewValue(a.Native, a.Negative != b.Negative, m.Uint64()+7, ao+bo+14)
	c := newAmount(v, a.Currency, a.Issuer)
	return c, c.canonicalise()
}

func (num *Amount) Divide(den *Amount) (*Amount, error) {
	if den.IsZero() {
		return nil, fmt.Errorf("Division by zero")
	}
	if num.IsZero() {
		return num.ZeroClone(), nil
	}
	av, bv := num.Num, den.Num
	ao, bo := num.Offset, den.Offset
	if num.Native {
		for ; av < minValue; av *= 10 {
			ao--
		}
	}
	if den.Native {
		for ; bv < minValue; bv *= 10 {
			bo--
		}
	}
	// Compute (numerator * 10^17) / denominator
	d := big.NewInt(0).SetUint64(av)
	d.Mul(d, bigTenTo17)
	d.Div(d, big.NewInt(0).SetUint64(bv))
	// 10^16 <= quotient <= 10^18
	if len(d.Bytes()) > 64 {
		return nil, fmt.Errorf("Divide: %s/%s", num.Debug(), den.Debug())
	}
	v := NewValue(num.Native, num.Negative != den.Negative, d.Uint64()+5, ao-bo-17)
	c := newAmount(v, num.Currency, num.Issuer)
	return c, c.canonicalise()
}

func (v *Value) IsScientific() bool {
	return v.Offset != 0 && (v.Offset < -25 || v.Offset > -5)
}

func (v *Value) IsZero() bool {
	return v.Num == 0
}

func (v *Value) String() string {
	if v.IsZero() {
		return "0"
	}
	if v.Native && v.IsScientific() {
		value := strconv.FormatUint(v.Num, 10)
		offset := strconv.FormatInt(v.Offset, 10)
		if v.Negative {
			return "-" + value + "e" + offset
		}
		return value + "e" + offset
	}
	value := big.NewInt(int64(v.Num))
	if v.Negative {
		value.Neg(value)
	}
	var offset *big.Int
	if v.Native {
		offset = big.NewInt(-6)
	} else {
		offset = big.NewInt(v.Offset)
	}
	exp := offset.Exp(bigTen, offset.Neg(offset), nil)
	rat := big.NewRat(0, 1).SetFrac(value, exp)
	left := rat.FloatString(0)
	if rat.IsInt() {
		return left
	}
	length := len(left)
	if v.Negative {
		length -= 1
	}
	return strings.TrimRight(rat.FloatString(32-length), "0")
}

func (a *Amount) String() string {
	switch {
	case a.Native:
		return a.Value.String() + "/XRP"
	case a.Issuer.IsZero():
		return a.Value.String() + "/" + a.Currency.String()
	default:
		issuer, _ := a.Issuer.MarshalText()
		return a.Value.String() + "/" + a.Currency.String() + "/" + string(issuer)
	}
}

func (a *Amount) JSON() string {
	b, _ := a.MarshalText()
	return string(b)
}
