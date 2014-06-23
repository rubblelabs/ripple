package data

import (
	"fmt"
	"github.com/donovanhide/ripple/crypto"
)

func CheckSymbol(h Hashable) string {
	if valid, err := CheckSignature(h); !valid || err != nil {
		return "✗"
	}
	return "✓"
}

func CheckSignature(h Hashable) (bool, error) {
	switch v := h.(type) {
	case *Validation:
		if err := NewEncoder().Validation(v, true); err != nil {
			return false, err
		}
		return crypto.Verify(v.SigningPubKey.Bytes(), v.Signature.Bytes(), v.Hash().Bytes())
	case *Proposal:
		if err := NewEncoder().SigningProposal(v); err != nil {
			return false, err
		}
		return crypto.Verify(v.PublicKey.Bytes(), v.Signature.Bytes(), v.Hash().Bytes())
	case *SetFee, *Amendment:
		return true, nil
	case Transaction:
		signingHash, err := NewEncoder().SigningHash(v)
		if err != nil {
			return false, err
		}
		base := v.GetBase()
		return crypto.Verify(base.SigningPubKey.Bytes(), base.TxnSignature.Bytes(), signingHash)
	default:
		return false, fmt.Errorf("Not a signed type")
	}
}

// Fills the raw field with a signed version of the encoding
func Sign(key crypto.Key, tx Transaction) error {
	enc := NewEncoder()
	signingHash, err := enc.SigningHash(tx)
	if err != nil {
		return err
	}
	sig, err := key.Sign(signingHash)
	if err != nil {
		return nil
	}
	vlSign := VariableLength(sig)
	tx.GetBase().TxnSignature = &vlSign
	return enc.Transaction(tx, false)
}
