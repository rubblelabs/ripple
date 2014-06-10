package data

import (
	"reflect"
)

type AffectedNode struct {
	FinalFields       interface{}     `json:",omitempty"`
	LedgerEntryType   LedgerEntryType `json:",omitempty"`
	LedgerIndex       *Hash256        `json:",omitempty"`
	PreviousFields    interface{}     `json:",omitempty"`
	NewFields         interface{}     `json:",omitempty"`
	PreviousTxnID     *Hash256        `json:",omitempty"`
	PreviousTxnLgrSeq *uint32         `json:",omitempty"`
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

type TransactionWithMetaData struct {
	Transaction
	MetaData       MetaData `json:"meta"`
	LedgerSequence uint32   `json:"ledger_index"`
}

type AccountRootDelta struct {
	Account           *Account
	FlagsDelta        *LedgerEntryFlag
	BalanceDelta      *Value
	TransferRateDelta int64
	RegularKey        *RegularKey
}

type RippleStateDelta struct {
	Account         *Account
	FlagsDelta      *LedgerEntryFlag
	BalanceDelta    *Amount
	QualityInDelta  int64
	QualityOutDelta int64
}

type OfferDelta struct {
	Account *Account
	Paid    *Amount
	Got     *Amount
}

type NodeDiff map[string]interface{}

func (m *MetaData) GetType() string { return "Metadata" }

func (m *NodeEffect) Action() string {
	switch {
	case m.ModifiedNode != nil:
		return "Modified"
	case m.DeletedNode != nil:
		return "Deleted"
	default:
		return "Created"
	}
}

func (m *NodeEffect) Diff() (NodeDiff, error) {
	diff := make(NodeDiff)
	switch {
	case m.CreatedNode != nil:
		fields := reflect.ValueOf(m.CreatedNode.NewFields).Elem()
		for i := 0; i < fields.NumField(); i++ {
			field, typ := fields.Field(i), fields.Type().Field(i)
			if field.IsNil() {
				continue
			}
			diff[typ.Name] = field.Elem().Interface()
		}
	case m.DeletedNode != nil:
		// previous := reflect.ValueOf(m.DeletedNode.PreviousFields).Elem()
		// final := reflect.ValueOf(m.DeletedNode.FinalFields).Elem()
		// for i := 0; i < final.NumField(); i++ {
		// 	typ := final.Type().Field(i)
		// 	fmt.Println(typ.Name, previous.Field(i).Interface(), final.Field(i).Interface())
		// 	// diff[typ.Name] = field.Interface()
		// }
	case m.ModifiedNode != nil && m.ModifiedNode.PreviousFields == nil:
		break
	case m.ModifiedNode != nil:
		previous := reflect.ValueOf(m.ModifiedNode.PreviousFields).Elem()
		final := reflect.ValueOf(m.ModifiedNode.FinalFields).Elem()
		for i := 0; i < final.NumField(); i++ {
			p, f := previous.Field(i), final.Field(i)
			if p.IsNil() || f.IsNil() {
				continue
			}
			typ := final.Type().Field(i)
			switch p.Interface().(type) {
			case *Value:

			case *Amount:
				change, err := p.Interface().(*Amount).Subtract(f.Interface().(*Amount))
				if err != nil {
					return nil, err
				}
				diff[typ.Name] = change
			case *uint32:
				diff[typ.Name] = int64(f.Elem().Uint()) - int64(p.Elem().Uint())
			}
		}
	}
	return diff, nil
}
