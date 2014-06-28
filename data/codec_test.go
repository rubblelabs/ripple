package data

import (
	"encoding/json"
	internal "github.com/donovanhide/ripple/testing"
	. "launchpad.net/gocheck"
)

type CodecSuite struct{}

var _ = Suite(&CodecSuite{})

func dump(test internal.TestData, v interface{}) CommentInterface {
	out, _ := json.Marshal(v)
	return Commentf("Test: %s\nJSON:%s\n", test.Description, string(out))
}

func (s *CodecSuite) TestParseLedgerHeaders(c *C) {
	for _, test := range internal.LedgerHeaders {
		ledger, err := ReadLedger(test.Reader())
		c.Assert(err, IsNil)
		msg := dump(test, ledger)
		_, value, err := Node(ledger)
		c.Assert(err, IsNil, msg)
		c.Assert(string(b2h(value)[26:]), Equals, test.Encoded, msg)
		// c.Assert(key, Equals, ledger.Hash())
	}
}

func (s *CodecSuite) TestParseTransactions(c *C) {
	for _, test := range internal.Transactions {
		tx, err := ReadTransaction(test.Reader())
		c.Assert(err, IsNil)
		msg := dump(test, tx)
		signable := tx.GetTransactionType() != SET_FEE && tx.GetTransactionType() != AMENDMENT
		ok, err := CheckSignature(tx)
		if signable {
			c.Assert(err, IsNil, msg)
		}
		c.Assert(ok, Equals, signable, msg)
		_, raw, err := Raw(tx)
		c.Assert(err, IsNil, msg)
		c.Assert(string(b2h(raw)), Equals, test.Encoded, msg)
	}
}

func (s *CodecSuite) TestValidations(c *C) {
	for _, test := range internal.Validations {
		v, err := ReadValidation(test.Reader())
		c.Assert(err, IsNil)
		msg := dump(test, v)
		ok, err := CheckSignature(v)
		c.Assert(ok, Equals, true, msg)
		c.Assert(err, IsNil, msg)
		_, raw, err := Raw(v)
		c.Assert(err, IsNil, msg)
		c.Assert(string(b2h(raw)), Equals, test.Encoded, msg)
	}
}

func (s *CodecSuite) TestParseNodes(c *C) {
	for _, test := range internal.Nodes {
		n, err := ReadPrefix(test.Reader())
		msg := dump(test, n)
		c.Assert(err, IsNil, msg)
		_, value, err := Node(n)
		c.Assert(err, IsNil, msg)
		c.Assert(string(b2h(value))[16:], Equals, test.Encoded[16:], msg)
	}
}
