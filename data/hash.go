package data

import (
	"encoding/hex"
	"fmt"
	"github.com/donovanhide/ripple/crypto"
)

type Hash128 [16]byte
type Hash160 [20]byte
type Hash256 [32]byte
type Vector256 []Hash256
type VariableLength []byte
type PublicKey [33]byte
type Account [20]byte
type RegularKey [20]byte

var zero256 Hash256
var zeroAccount Account

func (h Hash128) MarshalText() ([]byte, error) {
	return b2h(h[:]), nil
}

func (h Hash128) Bytes() []byte {
	return h[:]
}

func (h Hash160) MarshalText() ([]byte, error) {
	return b2h(h[:]), nil
}

func (h Hash160) Bytes() []byte {
	return h[:]
}

func NewHash256(s string) (Hash256, error) {
	var h Hash256
	n, err := hex.Decode(h[:], []byte(s))
	if err != nil {
		return h, err
	}
	if n != 32 {
		return h, fmt.Errorf("NewHash256: Wrong length %s", s)
	}
	return h, nil
}

func (h Hash256) MarshalText() ([]byte, error) {
	return b2h(h[:]), nil
}

func (h Hash256) IsZero() bool {
	return h == zero256
}

func (h Hash256) Bytes() []byte {
	return h[:]
}

func (h Hash256) String() string {
	return string(b2h(h[:]))
}

func (h Hash256) TruncatedString(length int) string {
	return string(b2h(h[:length]))
}

func (v VariableLength) MarshalText() ([]byte, error) {
	return b2h(v), nil
}

func (v *VariableLength) Bytes() []byte {
	if v != nil {
		return []byte(*v)
	}
	return []byte(nil)
}

func (p PublicKey) MarshalText() ([]byte, error) {
	if len(p) == 0 {
		return nil, nil
	}
	if pubKey, err := crypto.NewRipplePublicAccount(p[:]); err != nil {
		return nil, err
	} else {
		return []byte(pubKey.ToJSON()), nil
	}
}

func (p PublicKey) String() string {
	b, _ := p.MarshalText()
	return string(b)
}

func (p *PublicKey) Bytes() []byte {
	if p != nil {
		return p[:]
	}
	return []byte(nil)
}

func (a Account) MarshalText() ([]byte, error) {
	if len(a) == 0 {
		return nil, nil
	}
	if address, err := crypto.NewRippleAccount(a[:]); err != nil {
		return nil, err
	} else {
		return []byte(address.ToJSON()), nil
	}
}

func (a Account) String() string {
	b, _ := a.MarshalText()
	return string(b)
}

func (a Account) IsZero() bool {
	return a == zeroAccount
}

func (a *Account) Bytes() []byte {
	if a != nil {
		return a[:]
	}
	return []byte(nil)
}

func (r RegularKey) MarshalText() ([]byte, error) {
	if len(r) == 0 {
		return nil, nil
	}
	if address, err := crypto.NewRippleAccount(r[:]); err != nil {
		return nil, err
	} else {
		return []byte(address.ToJSON()), nil
	}
}

func (r *RegularKey) Bytes() []byte {
	if r != nil {
		return r[:]
	}
	return []byte(nil)
}
