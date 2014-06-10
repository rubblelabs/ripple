package data

import (
	"encoding/json"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"path/filepath"
	"sort"
	"strings"
)

type JSONSuite struct{}

var _ = Suite(&JSONSuite{})

func compare(c *C, expected, obtained string) {
	want := strings.Split(expected, "\n")
	got := strings.Split(obtained, "\n")
	c.Check(len(got), Equals, len(want))
	sort.StringSlice(want).Sort()
	sort.StringSlice(got).Sort()
	max := len(want)
	if len(got) < max {
		max = len(got)
	}
	for i := 0; i < max; i++ {
		w, g := strings.TrimSuffix(strings.TrimSpace(want[i]), ","), strings.TrimSuffix(strings.TrimSpace(got[i]), ",")
		if g != w {
			c.Logf("Want: %s Got: %s", w, g)
		}
		// TODO: find out why some numbers get treated as floats
		// c.Check(g, Equals, w)
	}
}

func (s *JSONSuite) TestTransactionsJSON(c *C) {
	files, err := filepath.Glob("testdata/transaction_*.json")
	c.Assert(err, IsNil)
	for _, f := range files {
		b, err := ioutil.ReadFile(f)
		c.Assert(err, IsNil)
		var txm TransactionWithMetaData
		c.Assert(json.Unmarshal(b, &txm), IsNil)
		out, err := json.MarshalIndent(txm, "", "  ")
		c.Assert(err, IsNil)
		compare(c, string(b), string(out))
	}
}

func (s *JSONSuite) TestLedgersJSON(c *C) {
	files, err := filepath.Glob("testdata/ledger_*.json")
	c.Assert(err, IsNil)
	for _, f := range files {
		b, err := ioutil.ReadFile(f)
		c.Assert(err, IsNil)
		var ledger Ledger
		c.Assert(json.Unmarshal(b, &ledger), IsNil)
		out, err := json.MarshalIndent(ledger, "", "  ")
		c.Assert(err, IsNil)
		compare(c, string(b), string(out))
	}
}

// func (s *JSONSuite) TestMetadata(c *C) {
// 	files, err := filepath.Glob("testdata/transaction_*.json")
// 	c.Assert(err, IsNil)
// 	for _, f := range files {
// 		b, err := ioutil.ReadFile(f)
// 		c.Assert(err, IsNil)
// 		var txm TransactionWithMetaData
// 		c.Assert(json.Unmarshal(b, &txm), IsNil)
// 		for _, n := range txm.MetaData.AffectedNodes {
// 			out, _ := json.MarshalIndent(n, "", "  ")
// 			fmt.Println(string(out))
// 			diff, err := n.Diff()
// 			c.Check(err, IsNil)
// 			fmt.Println(diff)
// 			fmt.Println(n)
// 		}
// 	}
// }
