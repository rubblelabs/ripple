package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/donovanhide/ripple/crypto"
	"github.com/donovanhide/ripple/data"
	"github.com/donovanhide/ripple/websockets"
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
	base.Sequence = uint32(c.GlobalInt("sequence"))
	base.SigningPubKey = new(data.PublicKey)
	if c.GlobalInt("lastledger") > 0 {
		base.LastLedgerSequence = new(uint32)
		*base.LastLedgerSequence = uint32(c.GlobalInt("lastledger"))
	}
	if base.Flags == nil {
		base.Flags = new(data.TransactionFlag)
	}
	copy(base.Account[:], id.Payload())
	copy(base.SigningPubKey[:], pub.Payload())
	if c.GlobalString("fee") != "" {
		fee, err := data.NewNativeValue(int64(c.GlobalInt("fee")))
		checkErr(err)
		base.Fee = *fee
	}
	tx.GetBase().TxnSignature = &data.VariableLength{}
	checkErr(data.Sign(tx, priv))
}

func submitTx(tx data.Transaction) {
	r, err := websockets.NewRemote("wss://s-east.ripple.com:443")
	checkErr(err)
	go r.Run()
	result, err := r.Submit(tx)
	checkErr(err)
	fmt.Printf("%s: %s\n", result.EngineResult, result.EngineResultMessage)
	os.Exit(0)
}

func payment(c *cli.Context) {
	// Validate and parse required fields
	if c.String("dest") == "" || c.String("amount") == "" || key == nil {
		fmt.Println("Destination, amount, and seed are required")
		os.Exit(1)
	}
	destination, amount := parseAccount(c.String("dest")), parseAmount(c.String("amount"))

	// Create payment and sign it
	payment := &data.Payment{
		Destination: *destination,
		Amount:      *amount,
	}
	payment.TransactionType = data.PAYMENT

	payment.Flags = new(data.TransactionFlag)
	if c.Bool("nodirect") {
		*payment.Flags = *payment.Flags | data.TxNoDirectRipple
	}
	if c.Bool("partial") {
		*payment.Flags = *payment.Flags | data.TxPartialPayment
	}
	if c.Bool("limit") {
		*payment.Flags = *payment.Flags | data.TxLimitQuality
	}

	sign(c, payment, 0)
	hash, raw, err := data.Raw(payment)
	checkErr(err)
	fmt.Printf("Hash: %X\nRaw:%X\n", hash, raw)

	// Print it in JSON
	out, err := json.Marshal(payment)
	checkErr(err)
	fmt.Println(string(out))

	if c.GlobalBool("submit") {
		submitTx(payment)
	}
}

func common(c *cli.Context) error {
	if c.GlobalString("seed") == "" {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}
	if c.GlobalInt("sequence") == 0 {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}
	key = parseSeed(c.String("seed"))
	return nil
}

var key *crypto.RootDeterministicKey

func main() {
	app := cli.NewApp()
	app.Name = "tx"
	app.Usage = "create a Ripple transaction. Sequence and seed must be specified for every command."
	app.Version = "0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{"seed,s", "", "the seed for the submitting account"},
		cli.IntFlag{"fee,f", 10, "the fee you want to pay"},
		cli.IntFlag{"sequence,q", 0, "the sequence for the transaction"},
		cli.IntFlag{"lastledger,l", 0, "highest ledger number that the transaction can appear in"},
		cli.BoolFlag{"submit,t", "submits the transaction via websocket"},
	}
	app.Before = common
	app.Commands = []cli.Command{{
		Name:        "payment",
		ShortName:   "p",
		Usage:       "create a payment",
		Description: "seed, sequence, destination and amount are required",
		Action:      payment,
		Flags: []cli.Flag{
			cli.StringFlag{"dest,d", "", "destination account"},
			cli.StringFlag{"amount,a", "", "amount to send"},
			cli.IntFlag{"tag,t", 0, "destination tag"},
			cli.StringFlag{"invoice,i", "", "invoice id (will be passed through SHA512Half)"},
			cli.StringFlag{"paths", "", "paths"},
			cli.StringFlag{"sendmax,m", "", "maximum to send"},
			cli.BoolFlag{"nodirect,r", "do not look for direct path"},
			cli.BoolFlag{"partial,p", "permit partial payment"},
			cli.BoolFlag{"limit,l", "limit quality"},
		},
	}}
	app.Run(os.Args)
}
