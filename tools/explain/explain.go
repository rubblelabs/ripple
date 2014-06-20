// Tool to explain transactions either individually, in a ledger or belonging to an account.
package main

import (
	"flag"
	"fmt"
	"github.com/donovanhide/ripple/data"
	"github.com/donovanhide/ripple/websockets"
	"os"
	"regexp"
	"strconv"
)

var argumentRegex = regexp.MustCompile(`(^[0-9a-fA-F]{64}$)|(^\d+$)|(^[r][a-km-zA-HJ-NP-Z0-9]{26,34}$)`)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func explain(txm *data.TransactionWithMetaData) {
	fmt.Println(txm.String())
	trades, err := txm.Trades()
	checkErr(err)
	for _, trade := range trades {
		fmt.Println("  ", trade.String())
	}
	balances, err := txm.Balances()
	checkErr(err)
	for _, balance := range balances {
		fmt.Println("  ", balance.String())
	}
}

func showUsage() {
	fmt.Println("Usage: explain [tx hash|ledger sequence|ripple address]")
	os.Exit(1)
}

var flags = flag.NewFlagSet("Flags", flag.ExitOnError)

var host = flags.String("host", "wss://s-east.ripple.com:443", "websockets host")

func main() {
	flags.Parse(os.Args[2:])
	if len(os.Args) == 1 {
		showUsage()
	}
	matches := argumentRegex.FindStringSubmatch(os.Args[1])
	r, err := websockets.NewRemote(*host)
	checkErr(err)
	fmt.Println("Connected to: ", *host)
	go r.Run()
	switch {
	case len(matches) == 0:
		showUsage()
	case len(matches[1]) > 0:
		hash, err := data.NewHash256(matches[1])
		checkErr(err)
		fmt.Println("Getting transaction: ", hash.String())
		result, err := r.Tx(*hash)
		checkErr(err)
		explain(&result.TransactionWithMetaData)
	case len(matches[2]) > 0:
		seq, err := strconv.ParseUint(matches[2], 10, 32)
		checkErr(err)
		ledger, err := r.Ledger(seq, true)
		checkErr(err)
		fmt.Println("Getting transactions for: ", seq)
		for _, tx := range ledger.Ledger.Transactions {
			explain(tx)
		}
	case len(matches[3]) > 0:
		account, err := data.NewAccountFromAddress(matches[3])
		checkErr(err)
		fmt.Println("Getting transactions for: ", account.String())
		for tx := range r.AccountTx(*account, 20) {
			explain(tx)
		}
	}
}
