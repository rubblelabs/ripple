package main

import (
	"fmt"
	"github.com/donovanhide/ripple/websockets"
)

func main() {
	r, err := websockets.NewRemote("wss://s-east.ripple.com:443")
	if err != nil {
		panic(err)
	}
	go r.Run()

	// Subscribe to ledger stream
	r.Outgoing <- websockets.SubscribeLedgers()
	confirmation := <-r.Incoming
	fmt.Printf(
		"Subscribed at index %d to streams: %v\n",
		confirmation.(*websockets.SubscribeLedgerCommand).Result.LedgerSequence,
		confirmation.(*websockets.SubscribeLedgerCommand).Streams,
	)

	// Consume ledgers as they arrive
	for {
		ledger := <-r.Incoming
		fmt.Printf(
			"Ledger %d closed at %s with %d transactions\n",
			ledger.(*websockets.LedgerStreamMsg).LedgerSequence,
			ledger.(*websockets.LedgerStreamMsg).LedgerTime.String(),
			ledger.(*websockets.LedgerStreamMsg).TxnCount,
		)
	}
}
