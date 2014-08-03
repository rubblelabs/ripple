package websockets

import (
	"github.com/golang/glog"
	"github.com/rubblelabs/ripple/data"
	"time"
)

// If no new ledger is emitted in this time, disconnects from server
// to try another one
const TIMEOUT = time.Second * 30

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

type Manager struct {
	// Ledgers are emitted on this channel in the correct order
	ledgers chan *ManagedLedger

	// Ledgers are captured in whatever order they arrive on this channel
	unorderedLedgers chan *ManagedLedger

	// Current websockets.Remote
	remote *Remote

	// Whether the remote is currently connected
	remoteConnected bool

	// Next ledger to anticipate
	nextLedgerSequence uint32
}

// NewManager creates a new Manager and starts emitting ManagedLedgers
// starting with startLedgerIdx.
func NewManager(startLedgerIdx uint32) *Manager {
	return &Manager{
		ledgers:            make(chan *ManagedLedger),
		unorderedLedgers:   make(chan *ManagedLedger),
		nextLedgerSequence: startLedgerIdx,
	}
}

// Ledgers returns a channel that emits ManagedLedgers.
func (m *Manager) Ledgers() <-chan *ManagedLedger {
	return m.ledgers
}

// Accepts unordered ledgers, buffers, and emits them in order.
func (m *Manager) runLedgerBuffer() {
	ledgerBuffer := make(map[uint32]*ManagedLedger, 10)

	// If TIMEOUT elapses without receiving the next ledger,
	// kill the connection and resume on another server.
	// This can overcome network problems, server bugs, etc.
	watchdogTimer := time.NewTimer(TIMEOUT)

	for {
		select {
		case <-watchdogTimer.C:
			if m.remoteConnected {
				m.remote.Close()
				m.remoteConnected = false
			}
			watchdogTimer.Reset(TIMEOUT)

		case l := <-m.unorderedLedgers:
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
				if len(ledgerBuffer) > 0 && m.remoteConnected {
					go m.getLedger(m.nextLedgerSequence)
				}

				watchdogTimer.Reset(TIMEOUT)

			} else if l.LedgerSequence > m.nextLedgerSequence {
				// Ledger is not the one we need. Stash it.

				// If a gap was just created, start filling it
				if len(ledgerBuffer) == 0 && m.remoteConnected {
					go m.getLedger(m.nextLedgerSequence)
				}

				// Stash this ledger in the buffer
				ledgerBuffer[l.LedgerSequence] = l

			} else {
				glog.Errorf("Received old ledger: %d", l.LedgerSequence)
			}
		}
	}
}

// Runs forever
func (m *Manager) Run() {
	uriIndex := 0
	go m.runLedgerBuffer()

	for {
		m.handleConnection(URIs[uriIndex])

		// Increment to next URI
		uriIndex = (uriIndex + 1) % len(URIs)

		glog.Infof("Reconnecting in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}

// Establishes a websockets connection and receives streaming messages
func (m *Manager) handleConnection(uri string) {
	var err error
	m.remote, err = NewRemote(uri)
	if err != nil {
		glog.Errorln(err.Error())
		return
	}
	go m.remote.Run()
	m.remoteConnected = true
	defer func() {
		m.remoteConnected = false
	}()

	res, err := m.remote.Subscribe(true, true, false)
	if err != nil {
		glog.Errorf(err.Error())
		m.remote.Close()
		return
	}
	glog.Infof(
		"Subscribed at %d\n",
		res.LedgerSequence,
	)

	if m.nextLedgerSequence == 0 {
		// Starting at 0 indicates start at current ledger
		m.nextLedgerSequence = res.LedgerSequence + 1
	} else if m.nextLedgerSequence < res.LedgerSequence {
		// If we need ledgers from before the subscription point, start retreiving them
		go m.getLedger(m.nextLedgerSequence)
	}

	// Receive streaming messages and collate them into ManagedLedgers
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
			glog.V(2).Infof("Txn %s", msg.Transaction.GetHash())
			currLedger.Transactions = append(currLedger.Transactions, &msg.Transaction)

			if len(currLedger.Transactions) == int(currLedgerTxnCount) {
				m.unorderedLedgers <- currLedger
				currLedger = nil
				currLedgerTxnCount = 0
			}

		default:
			panic("Unknown incoming message")

		}
	}
}

// Requestes an individual ledger
func (m *Manager) getLedger(sequence uint32) {
	res, err := m.remote.Ledger(sequence, true)
	if err != nil {
		glog.Errorf(err.Error())
		return
	}
	m.unorderedLedgers <- &ManagedLedger{
		LedgerSequence: res.Ledger.LedgerSequence,
		CloseTime:      res.Ledger.CloseTime,
		Hash:           res.Ledger.Hash,
		Transactions:   res.Ledger.Transactions,
	}
}
