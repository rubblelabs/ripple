package ledger

import (
	"fmt"
	"github.com/donovanhide/ripple/data"
	"github.com/donovanhide/ripple/storage"
	"github.com/donovanhide/ripple/terminal"
	"github.com/golang/glog"
	"time"
)

type Manager struct {
	missing  chan chan *data.Work
	incoming chan []data.Hashable
	current  chan uint32
	db       storage.DB
	ledgers  *data.LedgerSet
	started  time.Time
	stats    map[string]uint64
}

func NewManager(db storage.DB) (*Manager, error) {
	start := time.Now()
	ledgers, err := db.Ledger()
	if err != nil {
		return nil, err
	}
	glog.Infof("Manager: Created Ledger in %0.4f secs", time.Now().Sub(start).Seconds())
	return &Manager{
		missing:  make(chan chan *data.Work),
		incoming: make(chan []data.Hashable, 1000),
		current:  make(chan uint32),
		db:       db,
		ledgers:  ledgers,
		stats:    make(map[string]uint64),
	}, nil
}

func (m *Manager) Start() {
	m.started = time.Now()
	tick := time.NewTicker(time.Minute)
	var held CanonicalTxSet
	for {
		select {
		case <-tick.C:
			glog.Infoln("Manager:", m.String())
		case current := <-m.current:
			if current > m.ledgers.Max() {
				m.ledgers.Extend(current)
				glog.Infoln(current, m.ledgers.Max())
			}
		case in := <-m.incoming:
			for _, item := range in {
				terminal.Println(item, terminal.ShowTransactionId)
				switch v := item.(type) {
				case *data.Validation:
					continue
				case *data.Proposal:
					continue
				case *data.Ledger:
					m.stats["ledgers"]++
					wait := m.ledgers.Set(v.LedgerSequence)
					glog.V(2).Infof("Manager: Received: %d %0.04f/secs ", v.LedgerSequence, wait.Seconds())
					if err := m.db.Insert(v); err != nil {
						glog.Errorln("Manager: Ledger Insert:", err.Error())
					}
				case *data.TransactionWithMetaData:
					m.stats["transactions"]++
					if err := m.db.Insert(v); err != nil {
						glog.Errorln("Manager: Transaction Insert:", err.Error())
					}
				case data.Transaction:
					held.Add(v)
				}
			}
		case missing := <-m.missing:
			continue
			work := <-missing
			m.ledgers.Extend(work.End)
			work.MissingLedgers = m.ledgers.TakeMiddle(work.LedgerRange)
			missing <- work
		}
	}
}

func (m *Manager) Current(current uint32) {
	m.current <- current
}

func (m *Manager) Submit(items []data.Hashable) {
	m.incoming <- items
}

func (m *Manager) Missing(*data.LedgerRange) *data.Work {
	c := make(chan *data.Work)
	m.missing <- c
	return <-c
}
func (m *Manager) Copy() *RadixMap { return nil }

func (m *Manager) String() string {
	diff := time.Now().Sub(m.started).Seconds()
	ledgers, transactions := m.stats["ledgers"], m.stats["transactions"]
	ledgerRate, txRate := float64(ledgers)/diff, float64(transactions)/diff
	return fmt.Sprintf("%d %0.4f/sec Tx: %d %0.4f/sec Got: %d Max: %d", ledgers, ledgerRate, transactions, txRate, m.ledgers.Count(), m.ledgers.Max())
}
