package data

import (
	"io"
)

type Hashable interface {
	GetType() string
	Prefix() HashPrefix
	GetHash() *Hash256
}

type Signer interface {
	Hashable
	SigningPrefix() HashPrefix
	GetPublicKey() *PublicKey
	GetSignature() *VariableLength
}

type Router interface {
	Hashable
	SuppressionId() Hash256
}

type Storer interface {
	Hashable
	Ledger() uint32
	NodeType() NodeType
}

type LedgerEntry interface {
	Storer
	GetLedgerEntryType() LedgerEntryType
}

type Transaction interface {
	Signer
	GetTransactionType() TransactionType
	GetBase() *TxBase
	PathSet() PathSet
}

type Wire interface {
	Unmarshal(Reader) error
	Marshal(io.Writer) error
}
