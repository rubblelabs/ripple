package data

import (
	"bytes"
	"fmt"
)

// Helper functions for database/sql

func (h *Hash256) Scan(src interface{}) error {
	return scan(h[:], src, "Hash256")
}

func (a *Account) Scan(src interface{}) error {
	return scan(a[:], src, "Account")
}

func (a *PublicKey) Scan(src interface{}) error {
	return scan(a[:], src, "PublicKey")
}

func (v *Value) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Cannot scan %+v into Value", src)
	}
	return v.Unmarshal(bytes.NewReader(b))
}

func (a *Amount) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Cannot scan %+v into Amount", src)
	}
	return a.Unmarshal(bytes.NewReader(b))
}

func (t *RippleTime) Scan(src interface{}) error {
	v, ok := src.(int64)
	if !ok {
		return fmt.Errorf("Cannot scan %+v into RippleTime", src)
	}
	t.T = uint32(v)
	return nil
}

// support function for satisfying sql.Scanner interface
func scan(dest []byte, src interface{}, typ string) error {
	if src == nil {
		return nil
	}
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Cannot scan %+v into a %s", src, typ)
	}
	copy(dest, b)
	return nil
}
