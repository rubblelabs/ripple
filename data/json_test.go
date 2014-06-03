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

func (s *JSONSuite) TestTransactionsJSON(c *C) {
	files, err := filepath.Glob("testdata/*.json")
	c.Assert(err, IsNil)
	for _, f := range files {
		b, err := ioutil.ReadFile(f)
		c.Assert(err, IsNil)
		var txm TransactionWithMetaData
		c.Check(json.Unmarshal(b, &txm), IsNil)
		out, err := json.MarshalIndent(txm, "", "  ")
		c.Check(err, IsNil)
		want := strings.Split(string(b), "\n")
		got := strings.Split(string(out), "\n")
		c.Check(len(got), Equals, len(want))
		sort.StringSlice(want).Sort()
		sort.StringSlice(got).Sort()
		for i := range want {
			w, g := strings.TrimSuffix(strings.TrimSpace(want[i]), ","), strings.TrimSuffix(strings.TrimSpace(got[i]), ",")
			if g != w {
				c.Logf("Want: %s Got: %s", w, g)
			}
			// c.Check(g, Equals, w)
		}
	}
}
