package data

import (
	"fmt"
)

func format(h Hashable, format string, values ...interface{}) string {
	prefix := h.GetType() + ": "
	switch v := h.(type) {
	case Transaction:
		prefix += "Fee: %d Flags: %08X "
		base := v.GetBase()
		var flags uint32
		if base.Flags != nil {
			flags = *base.Flags
		}
		values = append([]interface{}{base.Fee.Num, flags}, values...)
	default:
	}
	return fmt.Sprintf(prefix+format, values...)
}

func (l *Ledger) String() string {
	return format(l, "%d %s", l.LedgerSequence, l.Hash().TruncatedString(8))
}

func (v *Validation) String() string {
	return format(v, "%d %d %s %d %s", v.LedgerSequence, v.BaseFee, v.LedgerHash.TruncatedString(8), v.SigningTime, v.SigningPubKey.String())
}

func (p *Proposal) String() string {
	return format(p, "%d", p.Sequence)
}

func (m *MetaData) String() string {
	return format(m, "")
}

func (p *Payment) String() string {
	return format(p, "%s => %s Amount: %s ", p.Account, p.Destination, p.Amount)
}

func (o *OfferCreate) String() string {
	ratio, err := o.TakerPays.Divide(&o.TakerGets)
	if err != nil {
		return "Bad OfferCreate"
	}
	return format(o, "%s Sequence: %d Pays: %s Gets: %s Ratio: %s", o.Account, o.Sequence, o.TakerPays, o.TakerGets, ratio.Value)
}

func (o *OfferCancel) String() string {
	return format(o, "%s Sequence: %d", o.Account, o.Sequence)
}

func (a *AccountSet) String() string {
	return format(a, "%s %d", a.Account, a.Sequence)
}

func (t *TrustSet) String() string {
	return format(t, "%s", t.Account)
}

func (f *SetFee) String() string {
	return format(f, "%d", f.BaseFee)
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
