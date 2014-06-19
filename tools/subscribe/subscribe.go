package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/donovanhide/ripple/data"
	"github.com/donovanhide/ripple/websockets"
	"github.com/fatih/color"
	"os"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	r, err := websockets.NewRemote("wss://s-east.ripple.com:443")
	checkErr(err)
	go r.Run()
	// Subscribe to all streams
	confirmation, err := r.Subscribe(true, true, true)
	checkErr(err)
	fmt.Printf(
		"Subscribed at %d\n",
		confirmation.LedgerSequence,
	)

	ledgerStyle := color.New(color.FgRed, color.Underline)
	transactionStyle := color.New(color.FgGreen)
	nodeStyle := color.New(color.FgBlue)
	serverStyle := color.New(color.FgMagenta)

	// Consume messages as they arrive
	for {
		msg, ok := <-r.Incoming
		if !ok {
			return
		}

		switch msg := msg.(type) {
		case *websockets.LedgerStreamMsg:
			ledgerStyle.Printf(
				"Ledger %d closed at %s with %d transactions\n",
				msg.LedgerSequence,
				msg.LedgerTime.String(),
				msg.TxnCount,
			)
		case *websockets.TransactionStreamMsg:
			transactionStyle.Printf(
				"    %-11s by %-34s Fee: %-8s Result: %s\n",
				msg.Transaction.GetTransactionType().String(),
				msg.Transaction.GetAccount(),
				msg.Transaction.GetBase().Fee,
				msg.EngineResult.String(),
			)
			trades, err := msg.Transaction.Trades()
			if err != nil {
				fmt.Println(err.Error())
			}
			for _, trade := range trades {
				nodeStyle.Printf("\t%s\n", trade.String())
			}
			balances, err := msg.Transaction.Balances()
			if err != nil {
				fmt.Println(err.Error())
			}
			for _, balance := range balances {
				nodeStyle.Printf("\t%s\n", balance.String())
			}

			// for _, n := range msg.Transaction.MetaData.AffectedNodes {
			// 	s := ExplainNodeEffect(&n)
			// 	if s != "" {
			// 		nodeStyle.Printf("        %s\n", s)
			// 	}
			// }
		case *websockets.ServerStreamMsg:
			serverStyle.Printf(
				"Server Status: %s (%d/%d)\n",
				msg.Status,
				msg.LoadFactor,
				msg.LoadBase,
			)
		}
	}
}

func ExplainNodeEffect(ne *data.NodeEffect) string {
	var op string
	var n *data.AffectedNode
	var fields interface{}

	switch {
	case ne.CreatedNode != nil:
		op = "Created"
		n = ne.CreatedNode
		fields = n.NewFields
	case ne.ModifiedNode != nil:
		op = "Modified"
		n = ne.ModifiedNode
		fields = n.FinalFields
	case ne.DeletedNode != nil:
		op = "Deleted"
		n = ne.DeletedNode
		fields = n.FinalFields
	}
	out, _ := json.Marshal(ne)
	fmt.Println(string(out))

	switch n.LedgerEntryType {
	// case data.DIRECTORY:
	// 	// Skip
	// 	return ""

	case data.OFFER:
		return fmt.Sprintf("%s Offer %s %s for %s",
			op,
			fields.(*data.OfferFields).Account,
			fields.(*data.OfferFields).TakerGets,
			fields.(*data.OfferFields).TakerPays,
		)

	default:
		return fmt.Sprintf("%s %s node: %s", op, n.LedgerEntryType, n.LedgerIndex)
	}
}
