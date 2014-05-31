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
	String() string
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
