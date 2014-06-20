package websockets

import (
	"encoding/json"
	"fmt"
	"github.com/donovanhide/ripple/data"
	"sync/atomic"
)

var counter uint64

type Syncer interface {
	Done()
}

type CommandError struct {
	Name    string `json:"error"`
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
}

type Command struct {
	*CommandError
	Id     uint64    `json:"id"`
	Name   string    `json:"command"`
	Type   string    `json:"type,omitempty"`
	Status string    `json:"status,omitempty"`
	Ready  chan bool `json:"-"`
}

func (c *Command) Done() {
	c.Ready <- true
}

func (c *Command) IncrementId() {
	c.Id = atomic.AddUint64(&counter, 1)
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("%s %d %s", e.Name, e.Code, e.Message)
}

func newCommand(command string) *Command {
	return &Command{
		Id:    atomic.AddUint64(&counter, 1),
		Name:  command,
		Ready: make(chan bool),
	}
}

type LedgerCommand struct {
	*Command
	LedgerIndex  interface{}   `json:"ledger_index"`
	Accounts     bool          `json:"accounts"`
	Transactions bool          `json:"transactions"`
	Expand       bool          `json:"expand"`
	Result       *LedgerResult `json:"result,omitempty"`
}

type LedgerResult struct {
	Ledger struct {
		LedgerSequence  uint32                `json:"ledger_index,string"`
		Accepted        bool                  `json:"accepted"`
		CloseTime       data.RippleTime       `json:"close_time"`
		Closed          bool                  `json:"closed"`
		Hash            data.Hash256          `json:"ledger_hash"`
		PreviousLedger  data.Hash256          `json:"parent_hash"`
		TotalXRP        uint64                `json:"total_coins,string"`
		AccountHash     data.Hash256          `json:"account_hash"`
		TransactionHash data.Hash256          `json:"transaction_hash"`
		Transactions    data.TransactionSlice `json:"transactions"`
	}
}

type AccountTxCommand struct {
	*Command
	Account   data.Account           `json:"account"`
	MinLedger int64                  `json:"ledger_index_min"`
	MaxLedger int64                  `json:"ledger_index_max"`
	Binary    bool                   `json:"binary,omitempty"`
	Forward   bool                   `json:"forward,omitempty"`
	Limit     int                    `json:"limit,omitempty"`
	Marker    map[string]interface{} `json:"marker,omitempty"`
	Result    *AccountTxResult       `json:"result,omitempty"`
}

type AccountTxResult struct {
	Marker       map[string]interface{} `json:"marker,omitempty"`
	Transactions data.TransactionSlice  `json:"transactions,omitempty"`
}

func newAccountTxCommand(account data.Account, pageSize int, marker map[string]interface{}) *AccountTxCommand {
	return &AccountTxCommand{
		Command:   newCommand("account_tx"),
		Account:   account,
		MinLedger: -1,
		MaxLedger: -1,
		Limit:     pageSize,
		Marker:    marker,
	}
}

type TxCommand struct {
	*Command
	Transaction data.Hash256 `json:"transaction"`
	Result      *TxResult    `json:"result,omitempty"`
}

type TxResult struct {
	data.TransactionWithMetaData
	Validated bool `json:"validated"`
}

// A shim to populate the Validated field before passing
// control on to TransactionWithMetaData.UnmarshalJSON
func (txr *TxResult) UnmarshalJSON(b []byte) error {
	var extract map[string]interface{}
	if err := json.Unmarshal(b, &extract); err != nil {
		return err
	}
	txr.Validated = extract["validated"].(bool)
	return json.Unmarshal(b, &txr.TransactionWithMetaData)
}

type SubmitCommand struct {
	*Command
	TxBlob string        `json:"tx_blob"`
	Result *SubmitResult `json:"result,omitempty"`
}

type SubmitResult struct {
	EngineResult        data.TransactionResult `json:"engine_result"`
	EngineResultCode    int                    `json:"engine_result_code"`
	EngineResultMessage string                 `json:"engine_result_message"`
	TxBlob              string                 `json:"tx_blob"`
	Tx                  interface{}            `json:"tx_json"`
}
