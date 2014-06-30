package data

import (
	"errors"
	"fmt"
	"reflect"
)

// ReadWire parses types received via the peer network
func ReadWire(r Reader, typ NodeType, ledgerSequence uint32) (Hashable, error) {
	version, err := readHashPrefix(r)
	if err != nil {
		return nil, err
	}
	switch version {
	case HP_LEAF_NODE:
		return readLedgerEntry(r)
	case HP_TRANSACTION_NODE:
		return readTransactionWithMetadata(r, ledgerSequence)
	case HP_INNER_NODE:
		return readCompressedInnerNode(r, typ)
	default:
		return nil, fmt.Errorf("Unknown hash prefix: %s", version.String())
	}
}

// ReadPrefix parses types received from the nodestore
func ReadPrefix(r Reader) (Storer, error) {
	header, err := readHeader(r)
	if err != nil {
		return nil, err
	}
	version, err := readHashPrefix(r)
	if err != nil {
		return nil, err
	}
	switch {
	case version == HP_INNER_NODE:
		return readInnerNode(r, header.NodeType)
	case header.NodeType == NT_LEDGER:
		return ReadLedger(r)
	case header.NodeType == NT_TRANSACTION_NODE:
		return readTransactionWithMetadata(r, header.LedgerSequence)
	case header.NodeType == NT_ACCOUNT_NODE:
		return readLedgerEntry(r)
	default:
		return nil, fmt.Errorf("Unknown node type")
	}
}

func ReadLedger(r Reader) (*Ledger, error) {
	ledger := new(Ledger)
	return ledger, read(r, &ledger.LedgerHeader)
}

func ReadValidation(r Reader) (*Validation, error) {
	validation := new(Validation)
	v := reflect.ValueOf(validation)
	if err := readObject(r, &v); err != nil {
		return nil, err
	}
	return validation, nil
}

func ReadTransaction(r Reader) (Transaction, error) {
	txType, err := expectType(r, "TransactionType")
	if err != nil {
		return nil, err
	}
	tx := TxFactory[txType]()
	v := reflect.ValueOf(tx)
	if err := readObject(r, &v); err != nil {
		return nil, err
	}
	return tx, nil
}

// ReadTransactionAndMetadata combines the inputs from the two
// readers into a TransactionWithMetaData
func ReadTransactionAndMetadata(tx, meta Reader, hash Hash256, ledger uint32) (*TransactionWithMetaData, error) {
	t, err := ReadTransaction(tx)
	if err != nil {
		return nil, err
	}
	txm := &TransactionWithMetaData{
		Transaction:    t,
		LedgerSequence: ledger,
	}
	m := reflect.ValueOf(&txm.MetaData)
	if err := readObject(meta, &m); err != nil {
		return nil, err
	}
	txm.SetHash(hash[:])
	return txm, nil
}

// For internal use when reading Prefix format
func readTransactionWithMetadata(r Reader, ledger uint32) (*TransactionWithMetaData, error) {
	br, err := NewVariableByteReader(r)
	if err != nil {
		return nil, err
	}
	tx, err := ReadTransaction(br)
	if err != nil {
		return nil, err
	}
	txm := &TransactionWithMetaData{
		Transaction:    tx,
		LedgerSequence: ledger,
	}
	br, err = NewVariableByteReader(r)
	if err != nil {
		return nil, err
	}
	meta := reflect.ValueOf(&txm.MetaData)
	if err := readObject(br, &meta); err != nil {
		return nil, err
	}
	hash, err := readHash(r)
	if err != nil {
		return nil, err
	}
	txm.SetHash(hash[:])
	return txm, nil
}

func readHashPrefix(r Reader) (HashPrefix, error) {
	var version HashPrefix
	return version, read(r, &version)
}

func readHeader(r Reader) (*NodeHeader, error) {
	header := new(NodeHeader)
	return header, read(r, header)
}

func readHash(r Reader) (*Hash256, error) {
	var h Hash256
	n, err := r.Read(h[:])
	switch {
	case err != nil:
		return nil, err
	case n != len(h):
		return nil, fmt.Errorf("Bad hash")
	default:
		return &h, nil
	}
}

func readInnerNode(r Reader, typ NodeType) (*InnerNode, error) {
	var inner InnerNode
	inner.Type = typ
	for i := range inner.Children {
		if _, err := r.Read(inner.Children[i][:]); err != nil {
			return nil, err
		}
	}
	return &inner, nil
}

func readCompressedInnerNode(r Reader, typ NodeType) (*InnerNode, error) {
	var inner InnerNode
	inner.Type = typ
	var entry CompressedNodeEntry
	for read(r, &entry) == nil {
		inner.Children[entry.Pos] = entry.Hash
	}
	return &inner, nil
}

func readLedgerEntry(r Reader) (LedgerEntry, error) {
	leType, err := expectType(r, "LedgerEntryType")
	if err != nil {
		return nil, err
	}
	le := LedgerEntryFactory[leType]()
	v := reflect.ValueOf(le)
	// LedgerEntries have 32 bytes of index suffixed
	// but don't have a variable bytes indicator
	lr := LimitedByteReader(r, int64(r.Len()-32))
	if err := readObject(lr, &v); err != nil {
		return nil, err
	}
	return le, nil
}

func expectType(r Reader, expected string) (uint16, error) {
	enc, err := readEncoding(r)
	if err != nil {
		return 0, err
	}
	name := encodings[*enc]
	if name != expected {
		return 0, fmt.Errorf("Unexpected type: %s expected: %s", name, expected)
	}
	var typ uint16
	return typ, read(r, &typ)
}

var (
	errorEndOfObject = errors.New("EndOfObject")
	errorEndOfArray  = errors.New("EndOfArray")
)

func readObject(r Reader, v *reflect.Value) error {
	var err error
	for enc, err := readEncoding(r); err == nil; enc, err = readEncoding(r) {
		name := encodings[*enc]
		// fmt.Println(name, v, v.IsValid(), enc.typ, enc.field)
		switch enc.typ {
		case ST_ARRAY:
			if name == "EndOfArray" {
				return errorEndOfArray
			}
			array := getField(v, enc)
		loop:
			for {
				child := reflect.New(array.Type().Elem()).Elem()
				err := readObject(r, &child)
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
			case "EndOfObject":
				return errorEndOfObject
			case "PreviousFields", "NewFields", "FinalFields":
				var fields Fields
				f := reflect.ValueOf(&fields)
				v.Elem().FieldByName(name).Set(f)
				if readObject(r, &f); err != nil && err != errorEndOfObject {
					return err
				}
			case "ModifiedNode", "DeletedNode", "CreatedNode":
				var node AffectedNode
				n := reflect.ValueOf(&node)
				var effect NodeEffect
				e := reflect.ValueOf(&effect)
				e.Elem().FieldByName(name).Set(n)
				v.Set(e.Elem())
				return readObject(r, &n)
			case "Memo":
				var memo Memo
				m := reflect.ValueOf(&memo)
				inner := reflect.ValueOf(&memo.Memo)
				err := readObject(r, &inner)
				v.Set(m.Elem())
				return err
			default:
				panic(fmt.Sprintf("Unknown object: %+v", enc))
			}
		default:
			if v.Kind() == reflect.Struct {
				return fmt.Errorf("Unexpected object: %s for field: %s", v.Type(), name)
			}
			field := getField(v, enc)
			switch v := field.Addr().Interface().(type) {
			case Wire:
				if err := v.Unmarshal(r); err != nil {
					return err
				}
			default:
				if err := read(r, v); err != nil {
					return err
				}
			}
		}
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
