package main

import (
	"flag"
	"log"
	"os"

	"github.com/rubblelabs/ripple/data"
	"github.com/rubblelabs/ripple/websockets"
)

func checkErr(err error, quit bool) {
	if err != nil {
		log.Println(err.Error())
		if quit {
			os.Exit(1)
		}
	}
}

var (
	host    = flag.String("host", "wss://s-east.ripple.com:443", "websockets host to connect to")
	account = flag.String("account", "", "optional account to monitor")
)

func main() {
	flag.Parse()
	var (
		filter *data.Account
		err    error
	)
	if len(*account) > 0 {
		filter, err = data.NewAccountFromAddress(*account)
		checkErr(err, true)
	}

	r, err := websockets.NewRemote(*host)
	checkErr(err, true)

	confirmation, err := r.Subscribe(true, true, false, false)
	checkErr(err, true)
	log.Printf("Subscribed at: %d ", confirmation.LedgerSequence)

	for {
		msg, ok := <-r.Incoming
		if !ok {
			return
		}
		switch msg := msg.(type) {
		case *websockets.TransactionStreamMsg:
			trades, err := data.NewTradeSlice(&msg.Transaction)
			checkErr(err, false)
			if filter != nil {
				trades = trades.Filter(*filter)
			}
			for _, trade := range trades {
				log.Println(trade)
			}
		}
	}
}
