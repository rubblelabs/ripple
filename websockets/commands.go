package websockets

import (
	"fmt"
	"github.com/donovanhide/ripple/data"
	"sync/atomic"
)

var counter uint64

type Command struct {
	Id      uint64 `json:"id"`
	Command string `json:"command"`
	Response
}

// Fields that are in every json response
type Response struct {
	Id           uint64
	Type         string
	Status       string
	Error        string
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func newCommand(command string) Command {
	return Command{
		Id:      atomic.AddUint64(&counter, 1),
		Command: command,
	}
}

type LedgerCommand struct {
	Command
	LedgerIndex  interface{} `json:"ledger_index"`
	Accounts     bool        `json:"accounts"`
	Transactions bool        `json:"transactions"`
	Expand       bool        `json:"expand"`
	Result       *struct {
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
	} `json:"result,omitempty"`
}

// Creates new `ledger` command to request a ledger by index
func Ledger(ledger interface{}, transactions bool) *LedgerCommand {
	return &LedgerCommand{
		Command:      newCommand("ledger"),
		LedgerIndex:  ledger,
		Transactions: transactions,
		Expand:       true,
	}
}

type TxCommand struct {
	Command
	Transaction string `json:"transaction"`
	Result      *struct {
		data.TransactionWithMetaData
		Validated bool `json:"validated"`
	} `json:"result,omitempty"`
}

// Creates new `tx` command to request a transaction by hash
func Tx(hash string) *TxCommand {
	return &TxCommand{
		Command:     newCommand("tx"),
		Transaction: hash,
	}
}

type SubmitCommand struct {
	Command
	TxBlob string `json:"tx_blob"`
	Result *struct {
		//FIXME(luke): TransactionResult doesn't support > 255 (tem, etc.)
		//EngineResult        data.TransactionResult `json:"engine_result"`
		EngineResult        string      `json:"engine_result"`
		EngineResultCode    int         `json:"engine_result_code"`
		EngineResultMessage string      `json:"engine_result_message"`
		TxBlob              string      `json:"tx_blob"`
		Tx                  interface{} `json:"tx_json"`
	}
}

func Submit(tx data.Transaction) *SubmitCommand {
	return &SubmitCommand{
		Command: newCommand("submit"),
		TxBlob:  fmt.Sprintf("%X", tx.Raw()),
	}
}
