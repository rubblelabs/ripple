package ledger

import (
	"fmt"
	"github.com/donovanhide/ripple/data"
	"github.com/donovanhide/ripple/storage"
	"github.com/golang/glog"
	"time"
)

type MissingLedgers struct {
	Request  *data.LedgerRange
	Response chan data.LedgerSlice
}

type Manager struct {
	Missing  chan *MissingLedgers
	Incoming chan interface{}
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
		Missing:  make(chan *MissingLedgers),
		Incoming: make(chan interface{}, 1000),
		db:       db,
		ledgers:  ledgers,
		stats:    make(map[string]uint64),
	}, nil
}

func (m *Manager) Start() {
	m.started = time.Now()
	tick := time.NewTicker(time.Minute)
	for {
		select {
		case <-tick.C:
			glog.Infoln("Manager:", m.String())
		case in := <-m.Incoming:
			switch v := in.(type) {
			case *data.Ledger:
				m.stats["ledgers"]++
				wait := m.ledgers.Set(v.LedgerSequence)
				glog.V(2).Infof("Manager: Received: %d %0.04f/secs ", v.LedgerSequence, wait.Seconds())
				if err := m.db.Insert(v); err != nil {
					glog.Errorln("Manager: Ledger Insert:", err.Error())
				}
			case data.Transaction:
				m.stats["transactions"]++
				if err := m.db.Insert(v); err != nil {
					glog.Errorln("Manager: Transaction Insert:", err.Error())
				}
			}
		case missing := <-m.Missing:
			m.ledgers.Extend(missing.Request.End)
			missing.Response <- m.ledgers.TakeMiddle(missing.Request)
		}
	}
}

func (m *Manager) String() string {
	diff := time.Now().Sub(m.started).Seconds()
	ledgers, transactions := m.stats["ledgers"], m.stats["transactions"]
	ledgerRate, txRate := float64(ledgers)/diff, float64(transactions)/diff
	return fmt.Sprintf("%d %0.4f/sec Tx: %d %0.4f/sec Got: %d Max: %d", ledgers, ledgerRate, transactions, txRate, m.ledgers.Count(), m.ledgers.Max())
}
