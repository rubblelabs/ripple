package data

import (
	"encoding/json"
	// "fmt"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"path/filepath"
)

type RippleSuite struct{}

var _ = Suite(&RippleSuite{})

func (s *JSONSuite) TestMetadata(c *C) {
	files, err := filepath.Glob("testdata/transaction_*.json")
	c.Assert(err, IsNil)
	for _, f := range files {
		b, err := ioutil.ReadFile(f)
		c.Assert(err, IsNil)
		var txm TransactionWithMetaData
		c.Assert(json.Unmarshal(b, &txm), IsNil)
		trades, err := txm.Trades()
		c.Check(err, IsNil)
		c.Check(len(trades), Equals, 8)
		balances, err := txm.Balances()
		c.Check(err, IsNil)
		c.Check(len(balances), Equals, 26)
		// sum, err := trades.Sum()
		// c.Check(err, IsNil)
		// c.Check(sum.String(), Equals, "8/BTC")
		// fmt.Println(trades.String())
		// fmt.Println(balances.String())
	}
}
