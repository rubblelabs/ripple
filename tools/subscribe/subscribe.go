package main

import (
	"flag"
	"fmt"
	"github.com/donovanhide/ripple/websockets"
)

func main() {
	flag.Parse()
	r, err := websockets.NewRemote("wss://s-east.ripple.com:443")
	if err != nil {
		panic(err)
	}
	go r.Run()

	// Subscribe to all streams
	r.Outgoing <- websockets.Subscribe(true, true, true)
	confirmation := <-r.Incoming
	fmt.Printf(
		"Subscribed at %d to streams: %v\n",
		confirmation.(*websockets.SubscribeCommand).Result.LedgerSequence,
		confirmation.(*websockets.SubscribeCommand).Streams,
	)

	// Consume messages as they arrive
	for {
		msg := <-r.Incoming
		switch msg := msg.(type) {
		case *websockets.LedgerStreamMsg:
			fmt.Printf(
				"Ledger %d closed at %s with %d transactions\n",
				msg.LedgerSequence,
				msg.LedgerTime.String(),
				msg.TxnCount,
			)
		case *websockets.TransactionStreamMsg:
			fmt.Printf(
				"    %s by %s\n",
				msg.Transaction.GetTransactionType().String(),
				msg.Transaction.GetAccount(),
			)
		case *websockets.ServerStreamMsg:
			fmt.Printf(
				"Server Status: %s (%d/%d)\n",
				msg.Status,
				msg.LoadFactor,
				msg.LoadBase,
			)
		}
	}
}
