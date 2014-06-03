package data

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/donovanhide/ripple/crypto"
	"regexp"
	"strconv"
	"time"
)

// wrapper type to enable second level of marshalling
type txmJSON TransactionWithMetaData

var txmTransactionTypeRegex = regexp.MustCompile(`"TransactionType":.*"(.*)"`)
var txmHashRegex = regexp.MustCompile(`"hash":.*"(.*)"`)

func (txm *TransactionWithMetaData) UnmarshalJSON(b []byte) error {
	matches := txmTransactionTypeRegex.FindAllStringSubmatch(string(b), 1)
	if matches == nil {
		return fmt.Errorf("Not a valid transaction with metadata: Missing TransactionType")
	}
	txType := matches[0][1]
	matches = txmHashRegex.FindAllStringSubmatch(string(b), 1)
	if matches == nil {
		return fmt.Errorf("Not a valid transaction with metadata: Missing hash")
	}
	hash := matches[0][1]
	txm.Transaction = GetTxFactoryByType(txType)()
	h, err := hex.DecodeString(hash)
	if err != nil {
		return fmt.Errorf("Bad hash: %s", hash)
	}
	txm.SetHash(h)
	if err := json.Unmarshal(b, txm.Transaction); err != nil {
		return err
	}
	return json.Unmarshal(b, (*txmJSON)(txm))
}

const txmFormat = `%s,"hash":"%s","inLedger":%d,"ledger_index":%d,"meta":%s}`

func (txm TransactionWithMetaData) MarshalJSON() ([]byte, error) {
	// This is an evil hack to be revisited
	tx, err := json.Marshal(txm.Transaction)
	if err != nil {
		return nil, err
	}
	meta, err := json.Marshal(txm.MetaData)
	if err != nil {
		return nil, err
	}
	out := fmt.Sprintf(txmFormat, string(tx[:len(tx)-1]), txm.Hash().String(), txm.LedgerSequence, txm.LedgerSequence, string(meta))
	return []byte(out), nil
}

func (r TransactionResult) MarshalText() ([]byte, error) {
	return []byte(resultNames[r]), nil
}

func (r *TransactionResult) UnmarshalText(b []byte) error {
	if result, ok := reverseResults[string(b)]; ok {
		*r = result
		return nil
	}
	return fmt.Errorf("Unknown TransactionResult: %s", string(b))
}

func (l LedgerEntryType) MarshalText() ([]byte, error) {
	return []byte(ledgerEntryNames[l]), nil
}

func (l *LedgerEntryType) UnmarshalText(b []byte) error {
	if leType, ok := ledgerEntryTypes[string(b)]; ok {
		*l = leType
		return nil
	}
	return fmt.Errorf("Unknown LedgerEntryType: %s", string(b))
}

func (t TransactionType) MarshalText() ([]byte, error) {
	return []byte(txNames[t]), nil
}

func (t *TransactionType) UnmarshalText(b []byte) error {
	if txType, ok := txTypes[string(b)]; ok {
		*t = txType
		return nil
	}
	return fmt.Errorf("Unknown TransactionType: %s", string(b))
}

func (t *RippleTime) UnmarshalJSON(b []byte) error {
	if unix, err := strconv.ParseInt(string(b), 10, 64); err != nil {
		return fmt.Errorf("Bad RippleTime:%s", string(b))
	} else {
		*t = RippleTime(time.Unix(unix+rippleEpoch, 0))
	}
	return nil
}

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

type amountJSON struct {
	Value    *Value   `json:"value"`
	Currency Currency `json:"currency"`
	Issuer   Account  `json:"issuer"`
}

func (a *Amount) MarshalJSON() ([]byte, error) {
	if a.Native {
		return []byte(`"` + strconv.FormatUint(a.Num, 10) + `"`), nil
	}
	return json.Marshal(amountJSON{a.Value, a.Currency, a.Issuer})
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
	var err error
	*c, err = NewCurrency(string(text))
	return err
}

func (h Hash128) MarshalText() ([]byte, error) {
	return b2h(h[:]), nil
}

func (h Hash128) UnmarshalText(b []byte) error {
	_, err := hex.Decode(h[:], b)
	return err
}

func (h Hash160) MarshalText() ([]byte, error) {
	return b2h(h[:]), nil
}

func (h Hash160) UnmarshalText(b []byte) error {
	_, err := hex.Decode(h[:], b)
	return err
}

func (h Hash256) MarshalText() ([]byte, error) {
	return b2h(h[:]), nil
}

func (h *Hash256) UnmarshalText(b []byte) error {
	_, err := hex.Decode(h[:], b)
	return err
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

// Expects variable length hex
func (v *VariableLength) UnmarshalText(b []byte) error {
	var err error
	*v, err = hex.DecodeString(string(b))
	return err
}

func (v VariableLength) MarshalText() ([]byte, error) {
	return b2h(v), nil
}

func (p PublicKey) MarshalText() ([]byte, error) {
	return b2h(p[:]), nil
}

// Expects public key hex
func (p *PublicKey) UnmarshalText(b []byte) error {
	_, err := hex.Decode(p[:], b)
	return err
}
