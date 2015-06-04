package websockets

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/rubblelabs/ripple/data"
)

var counter uint64

type Syncer interface {
	Done()
	Fail(message string)
}

type CommandError struct {
	Name    string `json:"error"`
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
}

type Command struct {
	*CommandError
	Id     uint64        `json:"id"`
	Name   string        `json:"command"`
	Type   string        `json:"type,omitempty"`
	Status string        `json:"status,omitempty"`
	Ready  chan struct{} `json:"-"`
}

func (c *Command) Done() {
	c.Ready <- struct{}{}
}

func (c *Command) Fail(message string) {
	c.CommandError = &CommandError{
		Name:    "Client Error",
		Code:    -1,
		Message: message,
	}
	c.Ready <- struct{}{}
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
		Ready: make(chan struct{}),
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

func newBinaryLedgerDataCommand(ledger interface{}, marker *data.Hash256) *BinaryLedgerDataCommand {
	return &BinaryLedgerDataCommand{
		Command: newCommand("ledger_data"),
		Ledger:  ledger,
		Binary:  true,
		Marker:  marker,
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

type LedgerCommand struct {
	*Command
	LedgerIndex  interface{}   `json:"ledger_index"`
	Accounts     bool          `json:"accounts"`
	Transactions bool          `json:"transactions"`
	Expand       bool          `json:"expand"`
	Result       *LedgerResult `json:"result,omitempty"`
}

type LedgerResult struct {
	Ledger data.Ledger
}

type LedgerHeaderCommand struct {
	*Command
	Ledger interface{} `json:"ledger"`
	Result *LedgerHeaderResult
}

type LedgerHeaderResult struct {
	Ledger         data.Ledger
	LedgerSequence uint32              `json:"ledger_index"`
	Hash           *data.Hash256       `json:"ledger_hash,omitempty"`
	LedgerData     data.VariableLength `json:"ledger_data"`
}

type LedgerDataCommand struct {
	*Command
	Ledger interface{}       `json:"ledger"`
	Marker *data.Hash256     `json:"marker,omitempty"`
	Result *LedgerDataResult `json:"result,omitempty"`
}

type BinaryLedgerDataCommand struct {
	*Command
	Ledger interface{}             `json:"ledger"`
	Binary bool                    `json:"binary"`
	Marker *data.Hash256           `json:"marker,omitempty"`
	Result *BinaryLedgerDataResult `json:"result,omitempty"`
}

type LedgerDataResult struct {
	LedgerSequence uint32                `json:"ledger_index,string"`
	Hash           data.Hash256          `json:"ledger_hash"`
	Marker         *data.Hash256         `json:"marker"`
	State          data.LedgerEntrySlice `json:"state"`
}

type BinaryLedgerData struct {
	Data  string `json:"data"`
	Index string `json:"index"`
}

type BinaryLedgerDataResult struct {
	LedgerSequence uint32             `json:"ledger_index,string"`
	Hash           data.Hash256       `json:"ledger_hash"`
	Marker         *data.Hash256      `json:"marker"`
	State          []BinaryLedgerData `json:"state"`
}

type RipplePathFindCommand struct {
	*Command
	SrcAccount    data.Account          `json:"source_account"`
	SrcCurrencies *[]data.Currency      `json:"source_currencies,omitempty"`
	DestAccount   data.Account          `json:"destination_account"`
	DestAmount    data.Amount           `json:"destination_amount"`
	Result        *RipplePathFindResult `json:"result,omitempty"`
}

type RipplePathFindResult struct {
	Alternatives []struct {
		SrcAmount      data.Amount  `json:"source_amount"`
		PathsComputed  data.PathSet `json:"paths_computed,omitempty"`
		PathsCanonical data.PathSet `json:"paths_canonical,omitempty"`
	}
	DestAccount    data.Account    `json:"destination_account"`
	DestCurrencies []data.Currency `json:"destination_currencies"`
}

type AccountInfoCommand struct {
	*Command
	Account data.Account       `json:"account"`
	Result  *AccountInfoResult `json:"result,omitempty"`
}

type AccountInfoResult struct {
	LedgerSequence uint32           `json:"ledger_current_index"`
	AccountData    data.AccountRoot `json:"account_data"`
}
