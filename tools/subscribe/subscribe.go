package main

import (
	"flag"
	"fmt"
	"github.com/donovanhide/ripple/terminal"
	"github.com/donovanhide/ripple/websockets"
	"os"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var (
	host = flag.String("host", "wss://s-east.ripple.com:443", "websockets host to connect to")
)

func main() {
	flag.Parse()
	r, err := websockets.NewRemote(*host)
	checkErr(err)
	go r.Run()
	// Subscribe to all streams
	confirmation, err := r.Subscribe(true, true, true)
	checkErr(err)
	fmt.Printf(
		"Subscribed at %d\n",
		confirmation.LedgerSequence,
	)

	// Consume messages as they arrive
	for {
		msg, ok := <-r.Incoming
		if !ok {
			return
		}

		switch msg := msg.(type) {
		case *websockets.LedgerStreamMsg:
			terminal.Println(msg, terminal.Default)
		case *websockets.TransactionStreamMsg:
			terminal.Println(msg.Transaction, terminal.Indent)
			trades, err := msg.Transaction.Trades()
			if err != nil {
				terminal.Println(err.Error(), terminal.Default)
			} else {
				for _, trade := range trades {
					terminal.Println(trade, terminal.DoubleIndent)
				}
			}
			balances, err := msg.Transaction.Balances()
			if err != nil {
				terminal.Println(err.Error(), terminal.Default)
			} else {
				for _, balance := range balances {
					terminal.Println(balance, terminal.DoubleIndent)
				}
			}
		case *websockets.ServerStreamMsg:
			terminal.Println(msg, terminal.Default)
		}
	}
}
