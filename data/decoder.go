package data

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type Decoder struct {
	r Reader
}

func NewDecoder(r Reader) *Decoder {
	return &Decoder{r}
}

func (dec *Decoder) Wire(typ NodeType) (Hashable, error) {
	version, err := dec.HashPrefix()
	if err != nil {
		return nil, err
	}
	switch version {
	case HP_LEAF_NODE:
		return dec.LedgerEntry()
	case HP_TRANSACTION_NODE:
		return dec.TransactionWithMetadata()
	case HP_INNER_NODE:
		return dec.CompressedInnerNode(typ)
	default:
		return nil, fmt.Errorf("Unknown hash prefix: %s", version.String())
	}
}

func (dec *Decoder) Prefix() (Hashable, error) {
	header, err := dec.Header()
	if err != nil {
		return nil, err
	}
	version, err := dec.HashPrefix()
	if err != nil {
		return nil, err
	}
	switch {
	case version == HP_INNER_NODE:
		return dec.InnerNode(header.NodeType)
	case header.NodeType == NT_LEDGER:
		return dec.Ledger()
	case header.NodeType == NT_TRANSACTION:
		return dec.Transaction()
	case header.NodeType == NT_TRANSACTION_NODE:
		return dec.TransactionWithMetadata()
	case header.NodeType == NT_ACCOUNT_NODE:
		return dec.LedgerEntry()
	default:
		return nil, fmt.Errorf("Unknown node type")
	}
}

func (dec *Decoder) Ledger() (*Ledger, error) {
	ledger := new(Ledger)
	return ledger, dec.read(&ledger.LedgerHeader)
}

func (dec *Decoder) Validation() (*Validation, error) {
	validation := new(Validation)
	v := reflect.ValueOf(validation)
	if err := dec.readObject(&v); err != nil {
		return nil, err
	}
	return validation, nil
}

func (dec *Decoder) HashPrefix() (HashPrefix, error) {
	var version HashPrefix
	return version, dec.read(&version)
}

func (dec *Decoder) Header() (*NodeHeader, error) {
	header := new(NodeHeader)
	return header, dec.read(header)
}

func (dec *Decoder) Hash() (*Hash256, error) {
	var h Hash256
	n, err := dec.r.Read(h[:])
	switch {
	case err != nil:
		return nil, err
	case n != len(h):
		return nil, fmt.Errorf("Bad hash")
	default:
		return &h, nil
	}
}

func (dec *Decoder) InnerNode(typ NodeType) (*InnerNode, error) {
	var inner InnerNode
	inner.Type = typ
	for i := range inner.Children {
		if _, err := dec.r.Read(inner.Children[i][:]); err != nil {
			return nil, err
		}
	}
	return &inner, nil
}

func (dec *Decoder) CompressedInnerNode(typ NodeType) (*InnerNode, error) {
	var inner InnerNode
	inner.Type = typ
	var entry CompressedNodeEntry
	for dec.read(&entry) == nil {
		inner.Children[entry.Pos] = entry.Hash
	}
	return &inner, nil
}

func (dec *Decoder) TransactionWithMetadata() (*TransactionWithMetaData, error) {
	br, err := NewVariableByteReader(dec.r)
	if err != nil {
		return nil, err
	}
	tx, err := NewDecoder(br).Transaction()
	if err != nil {
		return nil, err
	}
	txMeta := &TransactionWithMetaData{
		Transaction: tx,
	}
	br, err = NewVariableByteReader(dec.r)
	if err != nil {
		return nil, err
	}
	meta := reflect.ValueOf(&txMeta.MetaData)
	if err := NewDecoder(br).readObject(&meta); err != nil {
		return nil, err
	}
	return txMeta, nil
}

func (dec *Decoder) Transaction() (Transaction, error) {
	txType, err := dec.expectType("TransactionType")
	if err != nil {
		return nil, err
	}
	tx := TxFactory[txType]()
	v := reflect.ValueOf(tx)
	if err := dec.readObject(&v); err != nil {
		return nil, err
	}
	return tx, nil
}

func (dec *Decoder) LedgerEntry() (LedgerEntry, error) {
	leType, err := dec.expectType("LedgerEntryType")
	if err != nil {
		return nil, err
	}
	le := LedgerEntryFactory[leType]()
	v := reflect.ValueOf(le)
	// LedgerEntries have 32 bytes of hash suffixed
	// but don't have a variable bytes indicator
	lr := LimitedByteReader(dec.r, int64(dec.r.Len()-32))
	if err := NewDecoder(lr).readObject(&v); err != nil {
		return nil, err
	}
	return le, nil
}

func (dec *Decoder) next() (string, error) {
	var e enc
	if b, err := dec.r.ReadByte(); err != nil {
		return "", err
	} else {
		e.typ = b >> 4
		e.field = b & 0xF
	}
	var err error
	if e.typ == 0 {
		if e.typ, err = dec.r.ReadByte(); err != nil {
			return "", err
		}
	}
	if e.field == 0 {
		if e.field, err = dec.r.ReadByte(); err != nil {
			return "", err
		}
	}
	return encodings[e], nil
}

func (dec *Decoder) expectType(expected string) (uint16, error) {
	name, err := dec.next()
	if err != nil {
		return 0, err
	}
	if name != expected {
		return 0, fmt.Errorf("Unexpected type: %s expected: %s", name, expected)
	}
	var typ uint16
	return typ, dec.read(&typ)
}

func (dec *Decoder) read(dest interface{}) error {
	return binary.Read(dec.r, binary.BigEndian, dest)
}

func (dec *Decoder) readObject(v *reflect.Value) error {
	var err error
	for name, err := dec.next(); err == nil; name, err = dec.next() {
		// fmt.Println(name, v, v.IsValid())
		switch name {
		case "EndOfObject":
			return nil
		case "EndOfArray":
			continue
		case "PreviousFields", "NewFields", "FinalFields":
			ledgerEntryType := uint16(v.Elem().FieldByName("LedgerEntryType").Uint())
			le := fieldsFactory[ledgerEntryType]()
			lePtr := reflect.ValueOf(le)
			if err := dec.readObject(&lePtr); err != nil {
				return err
			}
			v.Elem().FieldByName(name).Set(lePtr)
		case "ModifiedNode", "DeletedNode", "CreatedNode":
			var node AffectedNode
			n := reflect.ValueOf(&node)
			if err := dec.readObject(&n); err != nil {
				return err
			}
			var effect NodeEffect
			e := reflect.ValueOf(&effect)
			e.Elem().FieldByName(name).Set(n)
			affected := v.Elem().FieldByName("AffectedNodes")
			affected.Set(reflect.Append(affected, e.Elem()))
		case "Memo":
			var memo Memo
			m := reflect.ValueOf(&memo.Memo)
			if err := dec.readObject(&m); err != nil {
				return err
			}
			memos := v.Elem().FieldByName("Memos")
			memos.Set(reflect.Append(memos, reflect.ValueOf(memo)))
		default:
			// fmt.Println(v, name)
			field := v.Elem().FieldByName(name)
			if field.Kind() == reflect.Ptr {
				field.Set(reflect.New(field.Type().Elem()))
				field = field.Elem()
			}
			if !field.IsValid() {
				return fmt.Errorf("Unknown Field: %s", name)
			}
			switch f := field.Addr().Interface().(type) {
			case Wire:
				if err := f.Unmarshal(dec.r); err != nil {
					return err
				}
			case *uint64, *uint32, *uint16, *uint8, *TransactionResult, *LedgerEntryType, *TransactionType:
				if err := dec.read(f); err != nil {
					return err
				}
			default:
				if err := dec.readObject(&field); err != nil {
					return err
				}
			}
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}
