package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/donovanhide/ripple/crypto"
	"github.com/donovanhide/ripple/data"
	"os"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func parseSeed(s string) *crypto.RootDeterministicKey {
	seed, err := crypto.NewRippleHashCheck(s, crypto.RIPPLE_FAMILY_SEED)
	checkErr(err)
	key, err := crypto.GenerateRootDeterministicKey(seed.Payload())
	checkErr(err)
	return key
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

func sign(c *cli.Context, tx data.Transaction, sequence int32) {
	priv, err := key.GenerateAccountKey(sequence)
	checkErr(err)
	id, err := key.GenerateAccountId(sequence)
	checkErr(err)
	pub, err := priv.PublicAccountKey()
	checkErr(err)
	base := tx.GetBase()
	base.SigningPubKey = new(data.PublicKey)
	base.Sequence = 44193
	base.Flags = new(uint32)
	copy(base.Account[:], id.Payload())
	copy(base.SigningPubKey[:], pub.Payload())
	if c.GlobalString("fee") != "" {
		base.Fee.Native = true
		checkErr(base.Fee.Parse(c.GlobalString("fee")))
	}
	checkErr(data.Sign(priv, tx))
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
	payment.TransactionType = data.PAYMENT
	sign(c, payment, 0)
	fmt.Printf("%X\n", payment.Raw())
	out, err := json.Marshal(payment)
	checkErr(err)
	fmt.Println(string(out))
	tx2, err := data.NewDecoder(bytes.NewReader(payment.Raw())).Transaction()
	checkErr(err)
	out2, err := json.Marshal(tx2)
	checkErr(err)
	fmt.Println(string(out2))
}

func common(c *cli.Context) error {
	key = parseSeed(c.String("seed"))
	return nil
}

var key *crypto.RootDeterministicKey
var account *data.Account

func main() {
	app := cli.NewApp()
	app.Name = "tx"
	app.Usage = "create a Ripple transaction"
	app.Version = "0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{"seed,s", "", "the seed for the submitting account"},
		cli.StringFlag{"fee,f", "", "the fee you want to pay"},
	}
	app.Before = common
	app.Commands = []cli.Command{{
		Name:        "payment",
		ShortName:   "p",
		Usage:       "create a payment",
		Description: "destination and amount are required",
		Action:      payment,
		Flags: []cli.Flag{
			cli.IntFlag{"tag,t", 0, "destination tag"},
			cli.StringFlag{"invoice,i", "", "invoice id (will be passed through SHA512Half)"},
			cli.StringFlag{"paths", "", "paths"},
			cli.StringFlag{"sendmax,m", "", "maximum to send"},
			cli.BoolTFlag{"direct,d", "look for direct path"},
			cli.BoolFlag{"partial,p", "permit partial payment"},
			cli.BoolFlag{"limit,l", "limit quality"},
		},
	}}
	app.Run(os.Args)
}
