package websockets

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/rubblelabs/ripple/data"
	"time"
)

const (
	// Number of out-of-order buffers to keep
	ledgerBufferSize = 10

	// Timeout for getLedger command
	getLedgerTimeout = 30 * time.Second
)

// Roster of websockets servers to rotate through
var URIs []string = []string{
	"wss://s-west.ripple.com:443",
	"wss://s-east.ripple.com:443",
	"wss://s1.ripple.com:443",
}

// A minimal ledger structure containing fields that can be populated
// by streaming OR the `ledger` command
type ManagedLedger struct {
	LedgerSequence uint32          `json:"ledger_index,string"`
	CloseTime      data.RippleTime `json:"close_time"`
	Hash           data.Hash256    `json:"ledger_hash"`
	Transactions   data.TransactionSlice
}

// Manager produces ledgers with transactions in order, starting from any index.
// If it encouters an error, it reconnects to a different server and picks up
// where it left off. It is ideal for applications that need to consume all
// ledgers despite server errors, network errors, and application restarts.
type Manager struct {
	// Ledgers are emitted on this channel in the correct order
	ledgers chan *ManagedLedger

	// Ledgers are captured in whatever order they arrive on this channel
	unorderedLedgers chan *ManagedLedger

	// Current websockets.Remote
	remote *Remote

	// Ledgers to be requested individually using the "ledger" command
	ledgerRequests chan uint32

	// Next ledger to anticipate
	nextLedgerSequence uint32
}

// NewManager creates a Manager and starts emitting ManagedLedgers
// starting with startLedgerIdx.
func NewManager(startLedgerIdx uint32) *Manager {
	m := &Manager{
		ledgers:            make(chan *ManagedLedger),
		unorderedLedgers:   make(chan *ManagedLedger),
		ledgerRequests:     make(chan uint32, 1),
		nextLedgerSequence: startLedgerIdx,
	}
	go m.run()
	return m
}

// Ledgers returns a channel that emits ManagedLedgers.
func (m *Manager) Ledgers() <-chan *ManagedLedger {
	return m.ledgers
}

// Spawns the following goroutines which conspire to provide gapless, ordered
// ledgers complete with transactions.
//
// +----------------------+     +----------------+
// | consolidateStreams() |     | ledgerGetter() |
// +-----------+----------+     +--------+-------+
//             |                         |
//             +-------------+-----------+
//                           |
//                           | unorderedLedgers
//                           |
//                  +--------+--------+
//                  | bufferLedgers() |
//                  +--------+--------+
//                           |
//                           v Ledgers()
//
func (m *Manager) run() {
	uriIndex := 0
	go m.bufferLedgers()

	for {
		m.handleConnection(URIs[uriIndex])

		// Increment to next URI
		uriIndex = (uriIndex + 1) % len(URIs)

		glog.Infof("Reconnecting in 1 seconds...")
		time.Sleep(1 * time.Second)
	}
}

// Accepts unordered ledgers, buffers them, and emits them in order.
// If ledgers are missing, it requests them from ledgerGetter.
func (m *Manager) bufferLedgers() {
	ledgerBuffer := make(map[uint32]*ManagedLedger)

	for l := range m.unorderedLedgers {
		if l.LedgerSequence == m.nextLedgerSequence {
			// Ledger is the one we need next. Emit it.
			m.ledgers <- l
			m.nextLedgerSequence++

			// If we already have any of the next ledgers, emit them now too.
			for ; ledgerBuffer[m.nextLedgerSequence] != nil; m.nextLedgerSequence++ {
				m.ledgers <- ledgerBuffer[m.nextLedgerSequence]
				delete(ledgerBuffer, m.nextLedgerSequence)
			}

			// If we still have a gap, request the next ledger
			if len(ledgerBuffer) > 0 {
				m.ledgerRequests <- m.nextLedgerSequence
			}

		} else if l.LedgerSequence > m.nextLedgerSequence {
			// Ledger is not the one we need. Stash it.

			// If a gap was just created, start filling it
			if len(ledgerBuffer) == 0 {
				m.ledgerRequests <- m.nextLedgerSequence
			}

			// Stash this ledger in the buffer if there is room
			if len(ledgerBuffer) < ledgerBufferSize {
				ledgerBuffer[l.LedgerSequence] = l
			}

		} else {
			glog.Errorf("Received old ledger: %d", l.LedgerSequence)
		}
	}
}

// Establishes a websockets connection, spawns a ledgerGetter, and then
// passes control to consolidateStreams().
func (m *Manager) handleConnection(uri string) {
	var err error
	m.remote, err = NewRemote(uri)
	if err != nil {
		glog.Errorln(err.Error())
		return
	}

	res, err := m.remote.Subscribe(true, true, true)
	if err != nil {
		glog.Errorf(err.Error())
		m.remote.Close()
		return
	}
	glog.Infof("Subscribed at %d\n", res.LedgerSequence)

	if m.nextLedgerSequence == 0 {
		// Starting at 0 indicates start at current ledger
		m.nextLedgerSequence = res.LedgerSequence + 1
	} else {
		m.ledgerRequests <- m.nextLedgerSequence
	}

	quitLedgerGetter := make(chan struct{}, 1)
	go ledgerGetter(m.remote, m.ledgerRequests, m.unorderedLedgers, quitLedgerGetter)
	defer func() {
		quitLedgerGetter <- struct{}{}
	}()

	m.consolidateStreams()
}

// Receives streaming messages and collate them into ManagedLedgers
func (m *Manager) consolidateStreams() {
	var currLedgerTxnCount uint32
	var currLedger *ManagedLedger
	for msg := range m.remote.Incoming {

		switch msg := msg.(type) {
		case *LedgerStreamMsg:
			glog.V(2).Infof("Ledger %d", msg.LedgerSequence)
			currLedger = &ManagedLedger{
				LedgerSequence: msg.LedgerSequence,
				CloseTime:      msg.LedgerTime,
				Hash:           msg.LedgerHash,
			}
			currLedgerTxnCount = msg.TxnCount

		case *TransactionStreamMsg:
			if currLedger == nil {
				// If we get a transaction before we see our first ledger,
				// ignore it.
				continue
			}

			glog.V(2).Infof("Txn %s", msg.Transaction.GetHash())
			currLedger.Transactions = append(currLedger.Transactions, &msg.Transaction)

			if len(currLedger.Transactions) == int(currLedgerTxnCount) {
				m.unorderedLedgers <- currLedger
				currLedger = nil
				currLedgerTxnCount = 0
			}

		case *ServerStreamMsg:
			glog.V(1).Infof("Server message: %s", msg.Status)

		default:
			panic("Unknown incoming message")

		}
	}
}

// Monitors ledgerRequests for ledgers that need to be fetched individually.
// Fetches them synchronously and sends them to unorderedLedgers.
func ledgerGetter(r *Remote, ledgerRequests <-chan uint32, unorderedLedgers chan<- *ManagedLedger, quit <-chan struct{}) {
	for {
		select {
		case getLedgerIdx := <-ledgerRequests:
			ml, err := getLedger(r, getLedgerIdx)
			if err != nil {
				glog.Errorf(err.Error())
				r.Close()
				return
			}
			unorderedLedgers <- ml

		case <-quit:
			return
		}
	}
}

// Requestes an individual ledger with timeout
func getLedger(r *Remote, sequence uint32) (*ManagedLedger, error) {
	var ledger *ManagedLedger

	timeout := time.NewTimer(getLedgerTimeout)
	result := make(chan error, 1)

	go func() {
		res, err := r.Ledger(sequence, true)
		if err != nil {
			result <- err
			return
		}
		ledger = &ManagedLedger{
			LedgerSequence: res.Ledger.LedgerSequence,
			CloseTime:      res.Ledger.CloseTime,
			Hash:           res.Ledger.Hash,
			Transactions:   res.Ledger.Transactions,
		}
		result <- nil
	}()

	select {
	case err := <-result:
		return ledger, err

	case <-timeout.C:
		return nil, fmt.Errorf("getLedger Timeout")
	}
}
