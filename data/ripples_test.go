package data

import (
	"encoding/json"
	// "fmt"
	"io/ioutil"

	. "gopkg.in/check.v1"
)

type RippleSuite struct{}

var _ = Suite(&RippleSuite{})

var expectedTradesAndBalances = map[string]struct {
	Balances    int
	Trades      int
	TotalTrades *Amount
}{
	"transaction_offercreate.json": {26, 8, amountCheck("8/BTC")},
}

func (s *JSONSuite) TestTradesAndBalances(c *C) {
	for file, expected := range expectedTradesAndBalances {
		b, err := ioutil.ReadFile("testdata/" + file)
		c.Assert(err, IsNil)
		var txm TransactionWithMetaData
		c.Assert(json.Unmarshal(b, &txm), IsNil)
		trades, err := txm.Trades()
		c.Check(err, IsNil)
		c.Check(len(trades), Equals, expected.Trades)
		balances, err := txm.Balances()
		c.Check(err, IsNil)
		c.Check(len(balances), Equals, expected.Balances)
		// sum, err := trades.Sum()
		// c.Check(err, IsNil)
		// c.Check(sum.Equals(*expected.TotalTrades), Equals, true)
	}
}
