package data

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/donovanhide/ripple/crypto"
	"strconv"
)

func (v *Value) MarshalText() ([]byte, error) {
	if v.Native {
		return []byte(strconv.FormatUint(v.Num, 10)), nil
	}
	return []byte(v.String()), nil
}

// Interpret as XRP in drips
func (v *Value) UnmarshalText(b []byte) (err error) {
	v.Native = true
	return v.Parse(string(b))
}

func (a *Amount) MarshalJSON() ([]byte, error) {
	if a.Native {
		return a.Value.MarshalText()
	}
	return json.Marshal(
		struct {
			Value    *Value   `json:"value"`
			Currency Currency `json:"currency"`
			Issuer   Account  `json:"issuer"`
		}{a.Value, a.Currency, a.Issuer})
}

func (a *Amount) UnmarshalJSON(b []byte) (err error) {
	a.Value = &Value{}

	// Try interpret as IOU
	var m map[string]string
	err = json.Unmarshal(b, &m)
	if err == nil {
		if err = a.Currency.UnmarshalText([]byte(m["currency"])); err != nil {
			return
		}

		a.Value.Native = false
		if err = a.Value.Parse(m["value"]); err != nil {
			return
		}

		if err = a.Issuer.UnmarshalText([]byte(m["issuer"])); err != nil {
			return
		}
		return
	}

	// Interpret as XRP in drips
	if err = a.Value.UnmarshalText(b[1 : len(b)-1]); err != nil {
		return
	}

	return
}

func (c Currency) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

func (c *Currency) UnmarshalText(text []byte) error {
	tmp, err := NewCurrency(string(text))
	if err != nil {
		return err
	}

	copy(c[:], tmp.Bytes())
	return nil
}

func (h Hash128) MarshalText() ([]byte, error) {
	return b2h(h[:]), nil
}

func (h Hash160) MarshalText() ([]byte, error) {
	return b2h(h[:]), nil
}

func (h Hash256) MarshalText() ([]byte, error) {
	return b2h(h[:]), nil
}

// Expects variable length hex
func (v *VariableLength) UnmarshalText(b []byte) error {
	_, err := hex.Decode(v.Bytes(), b)
	return err
}

func (v VariableLength) MarshalText() ([]byte, error) {
	return b2h(v), nil
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

// Expects base58-encoded account id
func (a *Account) UnmarshalText(text []byte) error {
	tmp, err := crypto.NewRippleHash(string(text))
	if err != nil {
		return err
	}
	if tmp.Version() != crypto.RIPPLE_ACCOUNT_ID {
		return fmt.Errorf("Incorrect version for Account: %d", tmp.Version())
	}

	copy(a[:], tmp.Payload())
	return nil
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

// Expects public key hex
func (p *PublicKey) UnmarshalText(text []byte) (err error) {

	var b []byte

	if b, err = hex.DecodeString(string(text)); err != nil {
		return err
	}

	copy(p[:], b)
	return nil
}
