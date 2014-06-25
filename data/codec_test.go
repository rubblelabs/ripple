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
		ledger, err := NewDecoder(test.Reader()).Ledger()
		c.Assert(err, IsNil)
		msg := dump(test, ledger)
		c.Assert(NewEncoder().Node(ledger), IsNil, msg)
		c.Assert(string(b2h(ledger.Raw())[26:]), Equals, test.Encoded, msg)
	}
}

func (s *CodecSuite) TestParseTransactions(c *C) {
	for _, test := range internal.Transactions {
		tx, err := NewDecoder(test.Reader()).Transaction()
		c.Assert(err, IsNil)
		msg := dump(test, tx)
		signable := tx.GetTransactionType() != SET_FEE && tx.GetTransactionType() != AMENDMENT
		ok, err := CheckSignature(tx)
		if signable {
			c.Assert(err, IsNil, msg)
		}
		c.Assert(ok, Equals, signable, msg)
		c.Assert(NewEncoder().Transaction(tx, false), IsNil, msg)
		c.Assert(string(b2h(tx.Raw())), Equals, test.Encoded, msg)
	}
}

func (s *CodecSuite) TestValidations(c *C) {
	for _, test := range internal.Validations {
		v, err := NewDecoder(test.Reader()).Validation()
		c.Assert(err, IsNil)
		msg := dump(test, v)
		ok, err := CheckSignature(v)
		c.Assert(ok, Equals, true, msg)
		c.Assert(err, IsNil, msg)
		c.Assert(NewEncoder().Validation(v, false), IsNil, msg)
		c.Assert(string(b2h(v.Raw())), Equals, test.Encoded, msg)
	}
}

func (s *CodecSuite) TestParseNodes(c *C) {
	for _, test := range internal.Nodes {
		n, err := NewDecoder(test.Reader()).Prefix()
		msg := dump(test, n)
		c.Assert(err, IsNil, msg)
		c.Assert(NewEncoder().Node(n), IsNil, msg)
		c.Assert(string(b2h(n.Raw()))[16:], Equals, test.Encoded[16:], msg)
	}
}
