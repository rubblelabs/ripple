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

type Response struct {
	Id           uint64
	Type         string
	Status       string
	Error        string
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type LedgerCommand struct {
	Command
	LedgerIndex  uint32 `json:"ledger_index"`
	Transactions bool   `json:"transactions"`
	Result       *struct {
		Ledger struct {
			LedgerSequence  string   `json:"ledger_index"`
			Accepted        bool     `json:"accepted"`
			CloseTime       uint64   `json:"close_time"`
			Closed          bool     `json:"closed"`
			Hash            string   `json:"ledger_hash"`
			PreviousLedger  string   `json:"parent_hash"`
			TotalXRP        string   `json:"total_coins"`
			AccountHash     string   `json:"account_hash"`
			TransactionHash string   `json:"transaction_hash"`
			Transactions    []string `json:"transactions"`
		}
	} `json:"result,omitempty"`
}

type TransactionCommand struct {
	Command
	Transaction string `json:"transaction"`
	Binary      bool   `json:"binary"`
	Result      *struct {
		Response
		Hash           string `json:"hash"`
		LedgerSequence uint32 `json:"ledger_index"`
		Meta           string `json:"meta"`
		Transaction    string `json:"tx"`
		Validated      bool   `json:"validated"`
	} `json:"result,omitempty"`
}

type StreamCommand struct {
	Command
	Streams []string `json:"streams"`
}

type LedgerStreamCommand struct {
	StreamCommand
	Result *struct {
		Response
		LedgerStream
	} `json:"result,omitempty"`
}

type LedgerStream struct {
	FeeBase          uint64          `json:"fee_base"`
	FeeRef           uint64          `json:"fee_ref"`
	LedgerSequence   uint32          `json:"ledger_index"`
	LedgerHash       string          `json:"ledger_hash"`
	LedgerTime       data.RippleTime `json:"ledger_time"`
	ReserveBase      uint64          `json:"reserve_base"`
	ReserveIncrement uint64          `json:"reserve_inc"`
	ValidatedLedgers string          `json:"validated_ledgers"`
}

func newCommand(command string) Command {
	return Command{
		Id:      atomic.AddUint64(&counter, 1),
		Command: command,
	}
}

func GetLedgerStream() *LedgerStreamCommand {
	return &LedgerStreamCommand{
		StreamCommand: StreamCommand{
			Command: newCommand("subscribe"),
			Streams: []string{"ledger"},
		},
	}
}

func GetLedger(ledger uint32) *LedgerCommand {
	return &LedgerCommand{
		Command:      newCommand("ledger"),
		LedgerIndex:  ledger,
		Transactions: true,
	}
}

func GetTransaction(hash string) *TransactionCommand {
	return &TransactionCommand{
		Command:     newCommand("tx"),
		Transaction: hash,
		Binary:      true,
	}
}
