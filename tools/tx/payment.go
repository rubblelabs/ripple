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

func NewOfferCreate(account data.Account, takerGets data.Amount, takerPays data.Amount) *data.OfferCreate {
	tr := (data.TxFactory[data.OFFER_CREATE]()).(*data.OfferCreate)
	tf := data.NoneFlags
	tr.Flags = &tf
	tr.TakerGets = takerGets
	tr.Account = account
	tr.TakerPays = takerPays
	f, _ := data.NewNativeValue(12)
	tr.Fee = *f
	return tr
}
