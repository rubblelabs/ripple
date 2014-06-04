package websockets

import (
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
	LedgerIndex  uint32 `json:"ledger_index"`
	Transactions bool   `json:"transactions"`
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
func GetLedger(ledger uint32) *LedgerCommand {
	return &LedgerCommand{
		Command:      newCommand("ledger"),
		LedgerIndex:  ledger,
		Transactions: true,
	}
}

type TransactionCommand struct {
	Command
	Transaction string `json:"transaction"`
	Binary      bool   `json:"binary"`
	Result      *struct {
		Hash           string `json:"hash"`
		LedgerSequence uint32 `json:"ledger_index"`
		Meta           string `json:"meta"`
		Transaction    string `json:"tx"`
		Validated      bool   `json:"validated"`
	} `json:"result,omitempty"`
}

// Creates new `tx` command to request a transaction by hash
func GetTransaction(hash string) *TransactionCommand {
	return &TransactionCommand{
		Command:     newCommand("tx"),
		Transaction: hash,
		Binary:      true,
	}
}
