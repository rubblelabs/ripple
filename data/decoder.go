package data

import (
	"encoding/binary"
	"errors"
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
		// TODO: What is the correct ledger sequence?
		return dec.TransactionWithMetadata(0)
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
		return dec.TransactionWithMetadata(header.LedgerIndex)
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

func (dec *Decoder) TransactionWithMetadata(ledger uint32) (*TransactionWithMetaData, error) {
	br, err := NewVariableByteReader(dec.r)
	if err != nil {
		return nil, err
	}
	tx, err := NewDecoder(br).Transaction()
	if err != nil {
		return nil, err
	}
	txMeta := &TransactionWithMetaData{
		Transaction:    tx,
		LedgerSequence: ledger,
	}
	br, err = NewVariableByteReader(dec.r)
	if err != nil {
		return nil, err
	}
	meta := reflect.ValueOf(&txMeta.MetaData)
	if err := NewDecoder(br).readObject(&meta); err != nil {
		return nil, err
	}
	hash := make([]byte, 32)
	n, err := dec.r.Read(hash)
	if err != nil {
		return nil, err
	}
	if n != 32 {
		return nil, fmt.Errorf("Bad hash")
	}
	txMeta.SetHash(hash)
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
	// LedgerEntries have 32 bytes of index suffixed
	// but don't have a variable bytes indicator
	lr := LimitedByteReader(dec.r, int64(dec.r.Len()-32))
	if err := NewDecoder(lr).readObject(&v); err != nil {
		return nil, err
	}
	return le, nil
}

func (dec *Decoder) next() (*enc, error) {
	var e enc
	if b, err := dec.r.ReadByte(); err != nil {
		return nil, err
	} else {
		e.typ = b >> 4
		e.field = b & 0xF
	}
	var err error
	if e.typ == 0 {
		if e.typ, err = dec.r.ReadByte(); err != nil {
			return nil, err
		}
	}
	if e.field == 0 {
		if e.field, err = dec.r.ReadByte(); err != nil {
			return nil, err
		}
	}
	return &e, nil
}

func (dec *Decoder) expectType(expected string) (uint16, error) {
	enc, err := dec.next()
	if err != nil {
		return 0, err
	}
	name := encodings[*enc]
	if name != expected {
		return 0, fmt.Errorf("Unexpected type: %s expected: %s", name, expected)
	}
	var typ uint16
	return typ, dec.read(&typ)
}

func (dec *Decoder) read(dest interface{}) error {
	return binary.Read(dec.r, binary.BigEndian, dest)
}

var (
	errorEndOfObject = errors.New("EndOfObject")
	errorEndOfArray  = errors.New("EndOfArray")
)

func (dec *Decoder) readObject(v *reflect.Value) error {
	var err error
	for enc, err := dec.next(); err == nil; enc, err = dec.next() {
		name := encodings[*enc]
		// fmt.Println(name, v, v.IsValid(), enc.typ, enc.field)
		if name == "EndOfArray" {
			return errorEndOfArray
		}
		if name == "EndOfObject" {
			return errorEndOfObject
		}
		switch enc.typ {
		case ST_ARRAY:
			array := getField(v, enc)
		loop:
			for {
				child := reflect.New(array.Type().Elem()).Elem()
				err := dec.readObject(&child)
				switch err {
				case errorEndOfArray:
					break loop
				case errorEndOfObject:
					array.Set(reflect.Append(*array, child))
				default:
					return err
				}
			}
		case ST_OBJECT:
			switch name {
			case "PreviousFields", "NewFields", "FinalFields":
				var fields Fields
				f := reflect.ValueOf(&fields)
				v.Elem().FieldByName(name).Set(f)
				if dec.readObject(&f); err != nil && err != errorEndOfObject {
					return err
				}
			case "ModifiedNode", "DeletedNode", "CreatedNode":
				var node AffectedNode
				n := reflect.ValueOf(&node)
				var effect NodeEffect
				e := reflect.ValueOf(&effect)
				e.Elem().FieldByName(name).Set(n)
				v.Set(e.Elem())
				return dec.readObject(&n)
			case "Memo":
				var memo Memo
				m := reflect.ValueOf(&memo)
				inner := reflect.ValueOf(&memo.Memo)
				err := dec.readObject(&inner)
				v.Set(m.Elem())
				return err
			default:
				panic(fmt.Sprintf("Unknown object: %+v", enc))
			}
		default:
			field := getField(v, enc)
			if w, ok := field.Addr().Interface().(Wire); ok {
				if err := w.Unmarshal(dec.r); err != nil {
					return err
				}
			} else {
				if err := dec.read(field.Addr().Interface()); err != nil {
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

func getField(v *reflect.Value, e *enc) *reflect.Value {
	name := encodings[*e]
	field := v.Elem().FieldByName(name)
	if field.Kind() == reflect.Ptr {
		field.Set(reflect.New(field.Type().Elem()))
		field = field.Elem()
	}
	return &field
}
