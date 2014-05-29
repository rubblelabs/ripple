package ledger

import (
	"sort"
)

type Queue []*LedgerState

type QueueFunc func(*LedgerState, *LedgerState) error

func (q Queue) Len() int           { return len(q) }
func (q Queue) Less(i, j int) bool { return q[i].Sequence() > q[j].Sequence() }
func (q Queue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func (q *Queue) Add(s *LedgerState) {
	*q = append(*q, s)
	sort.Sort(*q)
}

func (q *Queue) AddEmpty() {
	last := (*q)[len(*q)-1].Sequence()
	q.Add(NewEmptyLedgerState(last - 1))
}

func (q *Queue) Do(f QueueFunc) error {
	for {
		switch {
		case len(*q) < 2:
			return nil
		case (*q)[0].Sequence()-(*q)[1].Sequence() != 1:
			return nil
		default:
			if err := f((*q)[0], (*q)[1]); err != nil {
				return err
			}
			q.Pop()
		}
	}
}

func (q *Queue) Pop() bool {
	if len(*q) == 0 {
		return false
	}
	*q = (*q)[1:]
	return true
}
