package main

import (
	"flag"
	"fmt"
	"github.com/rubblelabs/ripple/websockets"
)

//TODO(luke): Merge this tool with subscribe.go
func main() {
	flag.Parse()
	m := websockets.NewManager(7302386)
	go m.Run()

	for ledger := range m.Ledgers {
		fmt.Printf(
			"Ledger %d closed with %d txns\n",
			ledger.LedgerSequence,
			len(ledger.Transactions),
		)
		for _, tx := range ledger.Transactions {
			fmt.Printf(
				"    %s\n",
				tx.GetHash(),
			)
		}
	}
}
