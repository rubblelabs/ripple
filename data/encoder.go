package data

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
)

func Raw(h Hashable) (Hash256, []byte, error) {
	return raw(h, h.Prefix(), nil, false)
}

func NodeId(h Hashable) (Hash256, error) {
	nodeid, _, err := raw(h, h.Prefix(), nil, false)
	return nodeid, err
}

func SigningHash(s Signable) (Hash256, []byte, error) {
	return raw(s, s.SigningPrefix(), nil, true)
}

func MultiSigningHash(s MultiSignable, account Account) (Hash256, []byte, error) {
	return raw(s, s.MultiSigningPrefix(), account.Bytes(), true)
}

func Node(h Storer) (Hash256, []byte, error) {
	var header bytes.Buffer
	for _, v := range []interface{}{h.Ledger(), h.Ledger(), h.NodeType(), h.Prefix()} {
		if err := write(&header, v); err != nil {
			return zero256, nil, err
		}
	}
	key, value, err := raw(h, h.Prefix(), nil, true)
	if err != nil {
		return zero256, nil, err
	}
	return key, append(header.Bytes(), value...), nil
}

func raw(value interface{}, prefix HashPrefix, suffix []byte, ignoreSigningFields bool) (Hash256, []byte, error) {
	buf := new(bytes.Buffer)
	hasher := sha512.New()
	multi := io.MultiWriter(buf, hasher)
	if err := write(hasher, prefix); err != nil {
		return zero256, nil, err
	}
	if err := writeRaw(multi, value, ignoreSigningFields); err != nil {
		return zero256, nil, err
	}
	if err := write(hasher, suffix); err != nil {
		return zero256, nil, err
	}
	var hash Hash256
	copy(hash[:], hasher.Sum(nil))
	return hash, buf.Bytes(), nil
}

// Disgusting node format and ordering handled here
func writeRaw(w io.Writer, value interface{}, ignoreSigningFields bool) error {
	switch v := value.(type) {
	case *Ledger:
		values := []interface{}{
			v.LedgerSequence,
			v.TotalXRP,
			v.PreviousLedger,
			v.TransactionHash,
			v.StateHash,
			v.ParentCloseTime,
			v.CloseTime,
			v.CloseResolution,
			v.CloseFlags,
		}
		for _, value := range values {
			if err := write(w, value); err != nil {
				return err
			}
		}
		return nil
	case *InnerNode:
		return write(w, v.Children)
	case *Validation:
		return encode(w, value, ignoreSigningFields)
	case *Proposal:
		if ignoreSigningFields {
			return writeValues(w, v.SigningValues())
		} else {
			return write(w, v)
		}
	case *TransactionWithMetaData:
		txid, tx, err := Raw(v.Transaction)
		if err != nil {
			return err
		}
		if err := writeVariableLength(w, tx); err != nil {
			return err
		}
		var meta bytes.Buffer
		if err := encode(&meta, &v.MetaData, false); err != nil {
			return err
		}
		if err := writeVariableLength(w, meta.Bytes()); err != nil {
			return err
		}
		return write(w, txid)
	case Transaction:
		return encode(w, value, ignoreSigningFields)
	case LedgerEntry:
		if err := encode(w, v, ignoreSigningFields); err != nil {
			return err
		}
		index, err := LedgerIndex(v)
		if err != nil {
			return err
		}
		return write(w, *index)
	default:
		return fmt.Errorf("Unknown type")
	}
}

func encode(w io.Writer, value interface{}, ignoreSigningFields bool) error {
	v := reflect.Indirect(reflect.ValueOf(value))
	fields := getFields(&v, 0, ignoreSigningFields)
	// fmt.Println(fields.String())
	return fields.Each(func(e enc, v interface{}) error {
		if err := writeEncoding(w, e); err != nil {
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

type field struct {
	encoding enc
	value    interface{}
	children fieldSlice
}

type fieldSlice []field

func (s fieldSlice) Len() int           { return len(s) }
func (s fieldSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s fieldSlice) Less(i, j int) bool { return s[i].encoding.Priority() < s[j].encoding.Priority() }

func (s fieldSlice) Sort() { sort.Sort(s) }

func (s *fieldSlice) Append(e enc, v interface{}, children fieldSlice) {
	*s = append(*s, field{e, v, children})
}

func getFields(v *reflect.Value, depth int, ignoreSigningFields bool) fieldSlice {
	// fmt.Println(v, v.Kind(), v.Type().Name())
	length := v.NumField()
	fields := make(fieldSlice, 0, length)
	typ := v.Type()
	for i := 0; i < length; i++ {
		fieldName := typ.Field(i).Name
		if fieldName == "Hash" || fieldName == "Id" {
			continue
		}
		// Stops LedgerEntryType being encoded for Fields
		if fieldName == "LedgerEntryType" && depth > 1 && typ.Name() == "leBase" {
			continue
		}
		encoding := reverseEncodings[fieldName]
		if ignoreSigningFields && encoding.SigningField() {
			continue
		}
		f := v.Field(i)
		// fmt.Println(fieldName, encoding, f, f.Kind())
		if f.Kind() == reflect.Interface {
			f = f.Elem()
		}
		if f.Kind() == reflect.Ptr {
			f = f.Elem()
		}
		if !f.IsValid() || (f.Kind() == reflect.Slice && f.Len() == 0) {
			continue
		}
		switch encoding.typ {
		case ST_UINT8, ST_UINT16, ST_UINT32, ST_UINT64:
			fields.Append(encoding, f.Addr().Interface(), nil)
		case ST_HASH96, ST_HASH128, ST_HASH160, ST_HASH192, ST_HASH256, ST_HASH384, ST_HASH512, ST_AMOUNT, ST_VL, ST_ACCOUNT, ST_PATHSET, ST_VECTOR256:
			fields.Append(encoding, f.Addr().Interface(), nil)
		case ST_ARRAY:
			var children fieldSlice
			for i := 0; i < f.Len(); i++ {
				f2 := f.Index(i)
				children = append(children, getFields(&f2, depth+1, ignoreSigningFields)...)
			}
			children.Append(reverseEncodings["EndOfArray"], nil, nil)
			fields.Append(encoding, nil, children)
		case ST_OBJECT:
			children := getFields(&f, depth+1, ignoreSigningFields)
			children.Append(reverseEncodings["EndOfObject"], nil, nil)
			fields.Append(encoding, nil, children)
		default:
			fields = append(fields, getFields(&f, depth+1, ignoreSigningFields)...)
		}
	}
	fields.Sort()
	return fields
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

func (f fieldSlice) String() string {
	var s []string
	f.Each(func(e enc, v interface{}) error {
		s = append(s, fmt.Sprintf("%s:%d:%d:%v", encodings[e], e.typ, e.field, v))
		return nil
	})
	return strings.Join(s, "\n")
}
