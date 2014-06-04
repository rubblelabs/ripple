package websockets

import (
	"encoding/json"
	"github.com/donovanhide/ripple/data"
)

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

// Fields from subscribed server status stream messages
type ServerStreamMsg struct {
	Status     string `json:"server_status"`
	LoadBase   int    `json:"load_base"`
	LoadFactor int    `json:"load_factor"`
}

// Map message types to the appropriate data structure
var streamMessageFactory = map[string]func() interface{}{
	"ledgerClosed": func() interface{} { return &LedgerStreamMsg{} },
	"transaction":  func() interface{} { return &TransactionStreamMsg{} },
	"serverStatus": func() interface{} { return &ServerStreamMsg{} },
}

type SubscribeCommand struct {
	Command
	Streams []string `json:"streams"`
	Result  *struct {
		// Contains one or both of these, depending what streams were subscribed
		*LedgerStreamMsg
		*ServerStreamMsg
	} `json:"result,omitempty"`
}

func Subscribe(ledger, transactions, server bool) *SubscribeCommand {
	streams := []string{}
	if ledger {
		streams = append(streams, "ledger")
	}
	if transactions {
		streams = append(streams, "transactions")
	}
	if server {
		streams = append(streams, "server")
	}
	return &SubscribeCommand{
		Command: newCommand("subscribe"),
		Streams: streams,
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
