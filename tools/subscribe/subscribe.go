package main

import (
	"flag"
	"fmt"
	"github.com/donovanhide/ripple/terminal"
	"github.com/donovanhide/ripple/websockets"
	"os"
)

func checkErr(err error, quit bool) {
	if err != nil {
		terminal.Println(err.Error(), terminal.Default)
		if quit {
			os.Exit(1)
		}
	}
}

var (
	host = flag.String("host", "wss://s-east.ripple.com:443", "websockets host to connect to")
)

func main() {
	flag.Parse()
	r, err := websockets.NewRemote(*host)
	checkErr(err, true)
	go r.Run()

	// Subscribe to all streams
	confirmation, err := r.Subscribe(true, true, true)
	checkErr(err, true)
	terminal.Println(fmt.Sprint("Subscribed at: ", confirmation.LedgerSequence), terminal.Default)

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
			for _, path := range msg.Transaction.PathSet() {
				terminal.Println(path, terminal.DoubleIndent)
			}
			trades, err := msg.Transaction.Trades()
			checkErr(err, false)
			for _, trade := range trades {
				terminal.Println(trade, terminal.DoubleIndent)
			}
			balances, err := msg.Transaction.Balances()
			checkErr(err, false)
			for _, balance := range balances {
				terminal.Println(balance, terminal.DoubleIndent)
			}
		case *websockets.ServerStreamMsg:
			terminal.Println(msg, terminal.Default)
		}
	}
}
