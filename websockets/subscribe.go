package websockets

import (
	"encoding/json"
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

func (msg *TransactionStreamMsg) UnmarshalJSON(b []byte) (err error) {
	var tmp map[string]json.RawMessage
	err = json.Unmarshal(b, &tmp)
	if err != nil {
		return
	}

	// Basic fields
	if err = json.Unmarshal(tmp["engine_result"], &msg.EngineResult); err != nil {
		return
	}
	if err = json.Unmarshal(tmp["engine_result_code"], &msg.EngineResultCode); err != nil {
		return
	}
	if err = json.Unmarshal(tmp["engine_result_message"], &msg.EngineResultMessage); err != nil {
		return
	}
	if err = json.Unmarshal(tmp["ledger_hash"], &msg.LedgerHash); err != nil {
		return
	}
	if err = json.Unmarshal(tmp["ledger_index"], &msg.LedgerSequence); err != nil {
		return
	}
	if err = json.Unmarshal(tmp["status"], &msg.Status); err != nil {
		return
	}
	if err = json.Unmarshal(tmp["validated"], &msg.Validated); err != nil {
		return
	}
	if err = json.Unmarshal(tmp["transaction"], &msg.Transaction); err != nil {
		return
	}

	// Transaction stream places the metadata *outside* of the transaction object.
	// We'll put it into the TransactionWithMetaData struct
	if err = json.Unmarshal(tmp["meta"], &msg.Transaction.MetaData); err != nil {
		return
	}

	// TransactionWithMetaData has a field for LedgerSequence too...
	if err = json.Unmarshal(tmp["ledger_index"], &msg.Transaction.LedgerSequence); err != nil {
		return
	}

	return
}

// Fields from subscribed transaction stream messages
type TransactionStreamMsg struct {
	Transaction         data.TransactionWithMetaData
	EngineResult        data.TransactionResult `json:"engine_result"`
	EngineResultCode    int                    `json:"engine_result_code"`
	EngineResultMessage string                 `json:"engine_result_message"`
	LedgerHash          data.Hash256           `json:"ledger_hash"`
	LedgerSequence      uint32                 `json:"ledger_index"`
	Status              string
	Validated           bool
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
