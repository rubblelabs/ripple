package ledger

import (
	"github.com/donovanhide/ripple/data"
	"sort"
)

type CanonicalTxSet struct {
	s                []data.Transaction
	lastClosedLedger data.Hash256
}

func (s CanonicalTxSet) Len() int      { return len(s.s) }
func (s CanonicalTxSet) Swap(i, j int) { s.s[i], s.s[j] = s.s[j], s.s[i] }
func (s CanonicalTxSet) Less(i, j int) bool {
	l, r := s.s[i].GetBase(), s.s[j].GetBase()
	la, ra := l.Account.Hash256().Xor(s.lastClosedLedger), r.Account.Hash256().Xor(s.lastClosedLedger)
	cmp := la.Compare(ra)
	switch {
	case cmp < 0:
		return true
	case cmp > 0:
		return false
	default:
		return l.Sequence < r.Sequence
	}
}

func (s CanonicalTxSet) Sort(lastClosedLedger data.Hash256) {
	s.lastClosedLedger = lastClosedLedger
	sort.Sort(s)
}

func (s *CanonicalTxSet) Add(tx data.Transaction) {
	(*s).s = append((*s).s, tx)
}
