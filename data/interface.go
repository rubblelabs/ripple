package data

import (
	"io"
)

type Hashable interface {
	GetType() string
	Prefix() HashPrefix
	GetHash() *Hash256
}

type Signable interface {
	Hashable
	InitialiseForSigning()
	SigningPrefix() HashPrefix
	GetPublicKey() *PublicKey
	GetSignature() *VariableLength
}

type MultiSignable interface {
	Hashable
	InitialiseForSigning()
	MultiSigningPrefix() HashPrefix
	GetPublicKey() *PublicKey
	GetSignature() *VariableLength
	SetSigners([]Signer)
	GetSigners() []Signer
}

type Router interface {
	Hashable
	SuppressionId() Hash256
}

type Storer interface {
	Hashable
	Ledger() uint32
	NodeType() NodeType
	NodeId() *Hash256
}

type LedgerEntry interface {
	Storer
	GetLedgerEntryType() LedgerEntryType
	GetLedgerIndex() *Hash256
	GetPreviousTxnId() *Hash256
	Affects(Account) bool
}

type Transaction interface {
	Signable
	GetTransactionType() TransactionType
	GetBase() *TxBase
	PathSet() PathSet
}

type Wire interface {
	Unmarshal(Reader) error
	Marshal(io.Writer) error
}
