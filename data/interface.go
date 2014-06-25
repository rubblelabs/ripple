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

type Signer interface {
	SigningHash() (Hash256, error)
	GetPublicKey() *PublicKey
	GetSignature() *VariableLength
}

type Router interface {
	SuppressionId() Hash256
}

type Storer interface {
	NodeId() Hash256
}

type Wire interface {
	Unmarshal(Reader) error
	Marshal(io.Writer) error
}

type LedgerEntry interface {
	Hashable
	GetLedgerEntryType() LedgerEntryType
}

type Transaction interface {
	Hashable
	Signer
	GetTransactionType() TransactionType
	GetBase() *TxBase
	PathSet() PathSet
}
