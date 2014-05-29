package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	internal "github.com/donovanhide/ripple/testing"
	. "launchpad.net/gocheck"
)

type CodecSuite struct{}

var _ = Suite(&CodecSuite{})

// func (s *CodecSuite) TestParseLedgerHeaders(c *C) {
// 	for _, test := range internal.LedgerHeaders {
// 		ledger, err := NewDecoder(test.Reader()).Ledger()
// 		c.Check(err, IsNil)
// 		out, _ := json.Marshal(ledger)
// 		fmt.Println(test.Description, string(out))
// 	}
// }

// func (s *CodecSuite) TestParseTransactions(c *C) {
// 	for _, test := range internal.Transactions {
// 		tx, err := NewDecoder(test.Reader()).Transaction()
// 		c.Check(err, IsNil, Commentf(test.Description))
// 		out, _ := json.Marshal(tx)
// 		fmt.Println(test.Description, string(out))
// 	}
// }

func (s *CodecSuite) TestParseNodes(c *C) {
	for _, test := range internal.Nodes {
		var buf bytes.Buffer
		n, err := NewDecoder(test.Reader()).Prefix()
		c.Check(err, IsNil, Commentf(test.Description))
		out, _ := json.MarshalIndent(n, "", "    ")
		fmt.Println(test.Description, string(out))
		c.Check(NewEncoder().Hex(&buf, n), IsNil)
		fmt.Printf("Hash: %s\nG: %s\nW: %s\n", buf.String()[:64], buf.String()[65:], test.Encoded)
	}
}
