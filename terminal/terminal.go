// Utiltities for formatting Ripple data in a terminal
package terminal

import (
	"github.com/donovanhide/ripple/data"
	"github.com/donovanhide/ripple/websockets"
	"github.com/fatih/color"
	"reflect"
)

type Flag uint32

const (
	NoIndent Flag = 1 << iota
	Indent
	DoubleIndent
	TripleIndent

	ShowLedgerSequence
)

var indents = map[Flag]string{
	NoIndent:     "",
	Indent:       "    ",
	DoubleIndent: "        ",
	TripleIndent: "            ",
}

var Default = NoIndent

var (
	ledgerStyle  = color.New(color.FgRed, color.Underline)
	txStyle      = color.New(color.FgGreen)
	tradeStyle   = color.New(color.FgBlue)
	balanceStyle = color.New(color.FgMagenta)
	infoStyle    = color.New(color.FgRed)
)

type bundle struct {
	color  *color.Color
	format string
	values []interface{}
	flag   Flag
}

func newTxBundle(txm *data.TransactionWithMetaData, flag Flag) *bundle {
	var (
		base   = txm.GetBase()
		format = "%-11s %-8s %s%s %-34s "
		values = []interface{}{base.GetType(), base.Fee, txm.MetaData.TransactionResult.Tick(), base.MemoSymbol(), base.Account}
		style  = txStyle
	)
	if !txm.MetaData.TransactionResult.Success() {
		style = infoStyle
	}
	if flag&ShowLedgerSequence > 0 {
		format = "%-9d " + format
		values = append([]interface{}{txm.LedgerSequence}, values...)
	}
	switch tx := txm.Transaction.(type) {
	case *data.Payment:
		format += "=> %-34s %-60s"
		values = append(values, []interface{}{tx.Destination, tx.Amount}...)
	case *data.OfferCreate:
		format += "%-60s %-60s"
		values = append(values, []interface{}{tx.TakerPays, tx.TakerGets}...)
	case *data.OfferCancel:
		format += "%-9d"
		values = append(values, tx.Sequence)
	case *data.AccountSet:
		format += "%-9d"
		values = append(values, tx.Sequence)
	case *data.TrustSet:
		format += "%-60s %d %d"
		values = append(values, tx.LimitAmount, tx.QualityIn, tx.QualityOut)
	}
	return &bundle{
		color:  style,
		format: format,
		values: values,
		flag:   flag,
	}
}

func newBundle(value interface{}, flag Flag) *bundle {
	switch v := reflect.Indirect(reflect.ValueOf(value)).Interface().(type) {
	case data.Ledger:
		return &bundle{
			color:  ledgerStyle,
			format: "Ledger %d closed at %s",
			values: []interface{}{v.LedgerSequence, v.CloseTime},
			flag:   flag,
		}
	case websockets.LedgerStreamMsg:
		return &bundle{
			color:  ledgerStyle,
			format: "Ledger %d closed at %s with %d transactions",
			values: []interface{}{v.LedgerSequence, v.LedgerTime.String(), v.TxnCount},
			flag:   flag,
		}
	case websockets.ServerStreamMsg:
		return &bundle{
			color:  infoStyle,
			format: "Server Status: %s (%d/%d)",
			values: []interface{}{v.Status, v.LoadFactor, v.LoadBase},
			flag:   flag,
		}
	case data.TransactionWithMetaData:
		return newTxBundle(&v, flag)
	case data.Trade:
		return &bundle{
			color:  tradeStyle,
			format: "Trade: %-34s => %-34s %-18s %60s => %-60s",
			values: []interface{}{v.Seller, v.Buyer, v.Price(), v.Paid, v.Got},
			flag:   flag,
		}
	case data.Balance:
		return &bundle{
			color:  balanceStyle,
			format: "Balance: %-34s  Currency: %s Balance: %20s Change: %20s",
			values: []interface{}{v.Account, v.Currency, v.Balance, v.Change},
			flag:   flag,
		}
	default:
		return &bundle{
			color:  infoStyle,
			format: "%s",
			values: []interface{}{v},
			flag:   flag,
		}
	}
}

func Println(value interface{}, flag Flag) (int, error) {
	b := newBundle(value, flag)
	return b.color.Printf(indents[flag]+b.format+"\n", b.values...)
}

func Sprint(value interface{}, flag Flag) string {
	b := newBundle(value, flag)
	return b.color.SprintfFunc()(indents[flag]+b.format, b.values...)
}
