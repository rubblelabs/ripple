package data

import (
	"io"
)

type Hashable interface {
	GetType() string
	Hash() Hash256
	Raw() []byte
	SetHash([]byte)
	SetRaw([]byte)
}

type Wire interface {
	Unmarshal(Reader) error
	Marshal(io.Writer) error
}

type LedgerEntry interface {
	Hashable
	GetLedgerEntryType() LedgerEntryType
	// SetLedgerSequence(uint32)
	// SetTransactionIndex(uint32)
}

type Transaction interface {
	Hashable
	GetTransactionType() TransactionType
	GetAccount() string
	// GetAffectedNodes() []NodeEffect
	GetBase() *TxBase
}

type LedgerSync interface {
	GetMissingLedgers(*LedgerRange) []uint32
	GetMissingTransactions(*LedgerRange) []Hash256
	GetMissingTransactionStates(*LedgerRange) []Hash256
	SubmitLedger(*Ledger)
	SubmitTransaction(Transaction)
}
