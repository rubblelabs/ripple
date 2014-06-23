package data

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"reflect"
	"sort"
	"strings"
)

type Encoder struct {
	buf   bytes.Buffer
	hash  hash.Hash
	multi io.Writer
}

func NewEncoder() *Encoder {
	enc := &Encoder{
		hash: sha512.New(),
	}
	enc.multi = io.MultiWriter(&enc.buf, enc.hash)
	return enc
}

func (enc *Encoder) Hex(w io.Writer, h Hashable) error {
	if err := enc.Node(h); err != nil {
		return err
	}
	_, err := fmt.Fprintf(w, "%s:%s\n", h.Hash().String(), b2h(enc.buf.Bytes()))
	return err
}

func (enc *Encoder) SigningHash(tx Transaction) ([]byte, error) {
	if err := enc.Transaction(tx, true); err != nil {
		return nil, err
	}
	enc.reset()
	if err := write(enc.hash, append(HP_TRANSACTION_SIGN.Bytes(), tx.Raw()...)); err != nil {
		return nil, err
	}
	return enc.hash.Sum(nil), nil

}

func (enc *Encoder) Transaction(tx Transaction, ignoreSigningFields bool) error {
	enc.reset()
	if err := enc.HashPrefix(enc.hash, tx); err != nil {
		return err
	}
	if err := enc.raw(enc.multi, tx, ignoreSigningFields); err != nil {
		return err
	}
	tx.SetHash(enc.hash.Sum(nil))
	tx.SetRaw(enc.buf.Bytes())
	return nil
}

func (enc *Encoder) Validation(v *Validation, ignoreSigningFields bool) error {
	enc.reset()
	if err := enc.HashPrefix(enc.hash, v); err != nil {
		return err
	}
	if err := enc.raw(enc.multi, v, ignoreSigningFields); err != nil {
		return err
	}
	v.SetHash(enc.hash.Sum(nil))
	v.SetRaw(enc.buf.Bytes())
	return nil
}

func (enc *Encoder) SigningProposal(p *Proposal) error {
	enc.reset()
	if err := enc.HashPrefix(enc.hash, p); err != nil {
		return err
	}
	values := []interface{}{p.Sequence, p.CloseTime.Uint32(), p.PreviousLedger.Bytes(), p.LedgerHash}
	for _, v := range values {
		if err := write(enc.hash, v); err != nil {
			return err
		}
	}
	p.SetHash(enc.hash.Sum(nil))
	p.SetRaw(enc.buf.Bytes())
	return nil
}

func (enc *Encoder) Node(h Hashable) error {
	enc.reset()
	if err := enc.Ledger(&enc.buf, h); err != nil {
		return err
	}
	if err := enc.Ledger(&enc.buf, h); err != nil {
		return err
	}
	if err := enc.NodeType(&enc.buf, h); err != nil {
		return err
	}
	if err := enc.HashPrefix(enc.multi, h); err != nil {
		return err
	}
	switch v := h.(type) {
	case *Ledger:
		if err := write(enc.multi, v.LedgerHeader); err != nil {
			return err
		}
	case *InnerNode:
		if err := write(enc.multi, v.Children); err != nil {
			return err
		}
	case *TransactionWithMetaData:
		var tx, meta bytes.Buffer
		txid := sha512.New()
		if err := write(txid, HP_TRANSACTION_ID); err != nil {
			return err
		}
		multi := io.MultiWriter(txid, &tx)
		if err := enc.raw(multi, v.Transaction, false); err != nil {
			return err
		}
		if err := enc.raw(&meta, &v.MetaData, false); err != nil {
			return err
		}
		if err := writeVariableLength(enc.multi, tx.Bytes()); err != nil {
			return err
		}
		if err := writeVariableLength(enc.multi, meta.Bytes()); err != nil {
			return nil
		}
		if err := write(enc.multi, txid.Sum(nil)[:32]); err != nil {
			return err
		}
	case LedgerEntry:
		if err := enc.raw(enc.multi, v, false); err != nil {
			return err
		}
		index, err := LedgerIndex(v)
		if err != nil {
			return err
		}
		if err = write(enc.multi, *index); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown type")
	}
	h.SetHash(enc.hash.Sum(nil))
	h.SetRaw(enc.buf.Bytes())
	return nil
}
func (enc *Encoder) reset() {
	enc.buf.Reset()
	enc.hash.Reset()
}

func write(w io.Writer, v interface{}) error {
	return binary.Write(w, binary.BigEndian, v)
}

func (enc *Encoder) Ledger(w io.Writer, h Hashable) error {
	switch v := h.(type) {
	case *Ledger:
		return write(w, v.LedgerSequence)
	case *TransactionWithMetaData:
		return write(w, v.LedgerSequence)
	case Transaction:
		return write(w, uint32(0))
	case *InnerNode, LedgerEntry:
		return write(w, uint32(0))
	default:
		return fmt.Errorf("Unknown type")
	}
}

func (enc *Encoder) NodeType(w io.Writer, h Hashable) error {
	switch v := h.(type) {
	case *Ledger:
		return write(w, NT_LEDGER)
	case *InnerNode:
		return write(w, v.Type)
	case Transaction:
		return write(w, NT_TRANSACTION_NODE)
	case LedgerEntry:
		return write(w, NT_ACCOUNT_NODE)
	default:
		return fmt.Errorf("Unknown type")
	}
}

func (enc *Encoder) HashPrefix(w io.Writer, h Hashable) error {
	switch h.(type) {
	case *Ledger:
		return write(w, HP_LEDGER_MASTER)
	case *InnerNode:
		return write(w, HP_INNER_NODE)
	case *TransactionWithMetaData:
		return write(w, HP_TRANSACTION_NODE)
	case Transaction:
		return write(w, HP_TRANSACTION_SIGN)
	case LedgerEntry:
		return write(w, HP_LEAF_NODE)
	case *Proposal:
		return write(w, HP_PROPOSAL)
	case *Validation:
		return write(w, HP_VALIDATION)
	default:
		return fmt.Errorf("Unknown type")
	}
}

type field struct {
	encoding enc
	value    interface{}
	children fieldSlice
}

type fieldSlice []field

func (s fieldSlice) Len() int           { return len(s) }
func (s fieldSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s fieldSlice) Less(i, j int) bool { return s[i].encoding.Priority() < s[j].encoding.Priority() }

func (s *fieldSlice) Append(e enc, v interface{}, children fieldSlice) {
	*s = append(*s, field{e, v, children})
}

func (s fieldSlice) Each(f func(e enc, v interface{}) error) error {
	for _, field := range s {
		if err := f(field.encoding, field.value); err != nil {
			return err
		}
		if err := field.children.Each(f); err != nil {
			return err
		}
	}
	return nil
}

func (s fieldSlice) Sort() { sort.Sort(s) }

func (f fieldSlice) String() string {
	var s []string
	f.Each(func(e enc, v interface{}) error {
		s = append(s, fmt.Sprintf("%s:%d:%d:%v", encodings[e], e.typ, e.field, v))
		return nil
	})
	return strings.Join(s, "\n")
}

func getFields(v *reflect.Value) fieldSlice {
	// fmt.Println(v, v.Kind(), v.Type().Name())
	var fields fieldSlice
	for i, length := 0, v.NumField(); i < length; i++ {
		f := v.Field(i)
		fieldName := v.Type().Field(i).Name
		encoding := reverseEncodings[fieldName]
		// fmt.Println(fieldName, encoding, f, f.Kind())
		if f.Kind() == reflect.Interface {
			f = f.Elem()
		}
		if f.Kind() == reflect.Ptr {
			f = f.Elem()
		}
		if !f.IsValid() || !f.CanInterface() || (f.Kind() == reflect.Slice && f.Len() == 0) {
			continue
		}
		switch encoding.typ {
		case ST_UINT8, ST_UINT16, ST_UINT32, ST_UINT64:
			fields.Append(encoding, f.Addr().Interface(), nil)
		case ST_HASH128, ST_HASH256, ST_AMOUNT, ST_VL, ST_ACCOUNT, ST_HASH160, ST_PATHSET, ST_VECTOR256:
			fields.Append(encoding, f.Addr().Interface(), nil)
		case ST_ARRAY:
			var children fieldSlice
			for i := 0; i < f.Len(); i++ {
				f2 := f.Index(i)
				children = append(children, getFields(&f2)...)
			}
			children.Append(reverseEncodings["EndOfArray"], nil, nil)
			fields.Append(encoding, nil, children)
		case ST_OBJECT:
			children := getFields(&f)
			children.Append(reverseEncodings["EndOfObject"], nil, nil)
			fields.Append(encoding, nil, children)
		default:
			fields = append(fields, getFields(&f)...)
		}
	}
	fields.Sort()
	return fields
}

func (encoder *Encoder) raw(w io.Writer, value interface{}, ignoreSigningFields bool) error {
	v := reflect.Indirect(reflect.ValueOf(value))
	fields := getFields(&v)
	// fmt.Println(fields.String())
	return fields.Each(func(e enc, v interface{}) error {
		if ignoreSigningFields && e.SigningField() {
			return nil
		}
		if err := encoder.writeEncoding(w, e); err != nil {
			return err
		}
		var err error
		switch v2 := v.(type) {
		case Wire:
			err = v2.Marshal(w)
		case nil:
			break
		default:
			err = write(w, v2)
		}
		return err
	})
}

func (enc *Encoder) writeEncoding(w io.Writer, e enc) error {
	switch {
	case e.typ < 16 && e.field < 16:
		return write(w, e.typ<<4|e.field)
	case e.typ < 16:
		return write(w, [2]uint8{e.typ << 4, e.field})
	case e.field < 16:
		return write(w, [2]uint8{e.field, e.typ})
	default:
		return write(w, [3]uint8{0, e.typ, e.field})
	}
}
