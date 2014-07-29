package data

import (
	"fmt"
	"sort"
)

type LedgerEntryState uint8

const (
	Created LedgerEntryState = iota
	Modified
	Deleted
	Touched
)

type AffectedNode struct {
	FinalFields       LedgerEntry `json:",omitempty"`
	LedgerEntryType   LedgerEntryType
	LedgerIndex       *Hash256    `json:",omitempty"`
	PreviousFields    LedgerEntry `json:",omitempty"`
	NewFields         LedgerEntry `json:",omitempty"`
	PreviousTxnID     *Hash256    `json:",omitempty"`
	PreviousTxnLgrSeq *uint32     `json:",omitempty"`
}

type NodeEffect struct {
	ModifiedNode *AffectedNode `json:",omitempty"`
	CreatedNode  *AffectedNode `json:",omitempty"`
	DeletedNode  *AffectedNode `json:",omitempty"`
}

type NodeEffects []NodeEffect

type MetaData struct {
	AffectedNodes     NodeEffects
	TransactionIndex  uint32
	TransactionResult TransactionResult
	DeliveredAmount   *Amount `json:",omitempty"`
}

type TransactionSlice []*TransactionWithMetaData

func (s TransactionSlice) Len() int      { return len(s) }
func (s TransactionSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TransactionSlice) Less(i, j int) bool {
	if s[i].LedgerSequence == s[j].LedgerSequence {
		return s[i].MetaData.TransactionIndex < s[j].MetaData.TransactionIndex
	}
	return s[i].LedgerSequence < s[j].LedgerSequence
}

func (s TransactionSlice) Sort() { sort.Sort(s) }

type TransactionWithMetaData struct {
	Transaction
	MetaData       MetaData `json:"meta"`
	LedgerSequence uint32   `json:"ledger_index"`
	Id             Hash256  `json:"-"`
}

func (t TransactionWithMetaData) GetType() string    { return "TransactionWithMetadata" }
func (t TransactionWithMetaData) Prefix() HashPrefix { return HP_TRANSACTION_NODE }
func (t TransactionWithMetaData) NodeType() NodeType { return NT_TRANSACTION_NODE }
func (t TransactionWithMetaData) Ledger() uint32     { return t.LedgerSequence }
func (t TransactionWithMetaData) NodeId() *Hash256   { return &t.Id }

func NewTransactionWithMetadata(typ TransactionType) *TransactionWithMetaData {
	return &TransactionWithMetaData{Transaction: TxFactory[typ]()}
}

func (effect *NodeEffect) AffectedNode() (*AffectedNode, LedgerEntry, LedgerEntryState) {
	switch {
	case effect.CreatedNode != nil && effect.CreatedNode.NewFields != nil:
		return effect.CreatedNode, effect.CreatedNode.NewFields, Created
	case effect.DeletedNode != nil && effect.DeletedNode.FinalFields != nil:
		return effect.DeletedNode, effect.DeletedNode.FinalFields, Deleted
	case effect.ModifiedNode != nil && effect.ModifiedNode.FinalFields != nil:
		return effect.ModifiedNode, effect.ModifiedNode.FinalFields, Modified
	case effect.ModifiedNode != nil && effect.ModifiedNode.FinalFields == nil:
		return effect.ModifiedNode, nil, Touched
	default:
		panic(fmt.Sprintf("Unknown LedgerEntryState: %+v", effect))
	}
}
