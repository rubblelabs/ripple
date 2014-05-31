package data

import (
	"fmt"
)

func format(h Hashable, format string, values ...interface{}) string {
	return fmt.Sprintf(h.GetType()+":"+format, values...)
}

func (l *Ledger) String() string {
	return format(l, "%d %s", l.LedgerSequence, l.Hash().TruncatedString(8))
}

func (v *Validation) String() string {
	return format(v, "%d %d %s %d %s", v.LedgerSequence, v.BaseFee, v.LedgerHash.String(), v.SigningTime, v.SigningPubKey.String())
}

func (p *Proposal) String() string {
	return format(p, "%d", p.Sequence)
}

func (m *MetaData) String() string {
	return format(m, "")
}

func (p *Payment) String() string {
	return format(p, "%s => %s Amount: %s", p.Account.String(), p.Destination.String(), p.Amount.String())
}

func (o *OfferCreate) String() string {
	return format(o, "%s Sequence: %d Pays: %s Gets: %s", o.Account.String(), o.Sequence, o.TakerPays.String(), o.TakerGets.String())
}

func (o *OfferCancel) String() string {
	return format(o, "%s Sequence: %d", o.Account, o.Sequence)
}

func (a *AccountSet) String() string {
	return format(a, "%s", a.Account)
}

func (t *TrustSet) String() string {
	return format(t, "%s", t.Account)
}

func (f *SetFee) String() string {
	return format(f, "%s", f.BaseFee)
}

func (a *Amendment) String() string {
	return format(a, "%s", a.Amendment)
}

func (s *SetRegularKey) String() string {
	return format(s, "%s %s", s.Account, s.RegularKey)
}

func (a *AccountRoot) String() string {
	return format(a, "%s %s", a.Account, a.Balance)
}

func (r *RippleState) String() string {
	return format(r, "%s %s %s", r.HighLimit, r.LowLimit, r.Balance)
}

func (o *Offer) String() string {
	return format(o, "%s Pays: %s Gets: %s", o.Account.String(), o.TakerGets.String(), o.TakerPays.String())
}

func (d *Directory) String() string {
	return format(d, "")
}

func (h *LedgerHashes) String() string {
	return format(h, "%d", len(h.Hashes))
}

func (s *FeeSetting) String() string {
	return format(s, "%d", s.BaseFee)
}

func (a *Amendments) String() string {
	return format(a, "%s", a.Amendments)
}
