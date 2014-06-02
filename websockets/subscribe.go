package websockets

import (
	"github.com/donovanhide/ripple/data"
)

// Map message types to the appropriate data structure
var streamMessageFactory = map[string]func() interface{}{
	"ledgerClosed": func() interface{} { return &LedgerStreamMsg{} },
	"transaction":  func() interface{} { return &TransactionStreamMsg{} },
	"serverStatus": func() interface{} { return &ServerStreamMsg{} },
}

type subscribeCommand struct {
	Command
	Streams []string `json:"streams"`
}

// Fields from subscribed ledger stream messages
type LedgerStreamMsg struct {
	FeeBase          uint64          `json:"fee_base"`
	FeeRef           uint64          `json:"fee_ref"`
	LedgerSequence   uint32          `json:"ledger_index"`
	LedgerHash       string          `json:"ledger_hash"`
	LedgerTime       data.RippleTime `json:"ledger_time"`
	ReserveBase      uint64          `json:"reserve_base"`
	ReserveIncrement uint64          `json:"reserve_inc"`
	ValidatedLedgers string          `json:"validated_ledgers"`
	TxnCount         uint32          `json:"txn_count"` // Only streamed, not in the subscribe result.
}

type SubscribeLedgerCommand struct {
	subscribeCommand
	Result *struct {
		LedgerStreamMsg
	} `json:"result,omitempty"`
}

// Creates new `subscribe` command to subscribe to ledgers
func SubscribeLedgers() *SubscribeLedgerCommand {
	return &SubscribeLedgerCommand{
		subscribeCommand: subscribeCommand{
			Command: newCommand("subscribe"),
			Streams: []string{"ledger"},
		},
	}
}

// Fields from subscribed transaction stream messages
type TransactionStreamMsg struct {
	//Transaction         data.Transaction
	EngineResult        string `json:"engine_result"`
	EngineResultCode    int    `json:"engine_result_code"`
	EngineResultMessage string `json:"engine_result_message"`
	LedgerHash          string `json:"ledger_hash"`
	LedgerSequence      uint32 `json:"ledger_index"`
	//Meta                interface{} `json:"meta_data"`
	Status    string
	Validated bool
}
type SubscribeTransactionCommand struct {
	subscribeCommand

	// Result object is empty in response
}

// Creates new `subscribe` command to subscribe to transactions
func SubscribeTransactions() *SubscribeTransactionCommand {
	return &SubscribeTransactionCommand{
		subscribeCommand: subscribeCommand{
			Command: newCommand("subscribe"),
			Streams: []string{"transactions"},
		},
	}
}

// Fields from subscribed server status stream messages
type ServerStreamMsg struct {
	Status     string `json:"server_status"`
	LoadBase   int    `json:"load_base"`
	LoadFactor int    `json:"load_factor"`
}

type SubscribeServerCommand struct {
	subscribeCommand
	Result *struct {
		ServerStreamMsg
	} `json:"result,omitempty"`
}

// Creates new `subscribe` command to subscribe to server status
func SubscribeServer() *SubscribeServerCommand {
	return &SubscribeServerCommand{
		subscribeCommand: subscribeCommand{
			Command: newCommand("subscribe"),
			Streams: []string{"server"},
		},
	}
}
