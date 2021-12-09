package tx

import (
	"github.com/maybeTomorrow/ripple/data"
)

func NewPayment(account data.Account, to data.Account, amount data.Amount) *data.Payment {
	tr := (data.TxFactory[data.PAYMENT]()).(*data.Payment)
	tr.Amount = amount
	tr.Destination = to
	tr.Account = account
	f, _ := data.NewNativeValue(12)
	tr.Fee = *f
	return tr
}
