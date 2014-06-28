package data

import (
	"io"
)

type Hashable interface {
	GetType() string
	Prefix() HashPrefix
}

type Signer interface {
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
	Hash() Hash256
	Raw() []byte
	SetHash([]byte)
	SetRaw([]byte)
	Ledger() uint32
	NodeType() NodeType
}

type LedgerEntry interface {
	Storer
	GetLedgerEntryType() LedgerEntryType
}

type Transaction interface {
	Hashable
	Signer
	GetTransactionType() TransactionType
	GetBase() *TxBase
	PathSet() PathSet
}

type Wire interface {
	Unmarshal(Reader) error
	Marshal(io.Writer) error
}
