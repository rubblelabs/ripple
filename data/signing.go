package data

import (
	"fmt"
	"github.com/donovanhide/ripple/crypto"
)

func CheckSignature(h Hashable) (bool, error) {
	switch v := h.(type) {
	case *Validation:
		if err := NewEncoder().Validation(v, true); err != nil {
			return false, err
		}
		return crypto.Verify(v.SigningPubKey.Bytes(), v.Signature.Bytes(), v.Hash().Bytes())
	case *SetFee, *Amendment:
		return true, nil
	case Transaction:
		if err := NewEncoder().Transaction(v, true); err != nil {
			return false, err
		}
		base := v.GetBase()
		return crypto.Verify(base.SigningPubKey.Bytes(), base.TxnSignature.Bytes(), v.Hash().Bytes())
	default:
		return false, fmt.Errorf("Not a signed type")
	}
}

// Fills the raw field with a signed version of the encoding
func Sign(key crypto.Key, tx Transaction) error {
	enc := NewEncoder()
	if err := enc.Transaction(tx, true); err != nil {
		return err
	}
	sig, err := key.Sign(tx.Raw())
	if err != nil {
		return nil
	}
	vlSign := VariableLength(sig)
	tx.GetBase().TxnSignature = &vlSign
	return enc.Transaction(tx, false)
}
