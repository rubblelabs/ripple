package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/donovanhide/ripple/data"
	"os"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func parseAccount(s string) *data.Account {
	account, err := data.NewAccountFromAddress(s)
	checkErr(err)
	return account
}

func parseAmount(s string) *data.Amount {
	amount, err := data.NewAmount(s)
	checkErr(err)
	return amount
}

func encode(tx data.Transaction) string {
	checkErr(data.NewEncoder().Transaction(tx, true))
	return string(tx.Raw())
}

func payment(c *cli.Context) {
	if c.Args().Get(1) == "" {
		fmt.Println("Destination and amount are required")
		os.Exit(1)
	}
	destination, amount := parseAccount(c.Args().Get(0)), parseAmount(c.Args().Get(1))
	payment := &data.Payment{
		Destination: *destination,
		Amount:      *amount,
	}
	fmt.Println("Payment:", destination.String(), amount.String())
	fmt.Printf("%X\n", encode(payment))
}

func main() {
	app := cli.NewApp()
	app.Name = "tx"
	app.Usage = "create a Ripple transaction"
	app.Version = "0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{"seed,s", "", "the seed for the submitting account"},
	}
	app.Commands = []cli.Command{{
		Name:        "payment",
		ShortName:   "p",
		Usage:       "create a payment",
		Description: "destination and amount are required",
		Action:      payment,
		Flags: []cli.Flag{
			cli.IntFlag{"tag,t", 0, "destination tag"},
			cli.StringFlag{"invoice,i", "", "invoice id (will be passed through SHA512Half)"},
			cli.StringFlag{"paths,p", "", "paths"},
			cli.StringFlag{"sendmax,m", "", "maximum to send"},
		},
	}}
	app.Run(os.Args)
}
