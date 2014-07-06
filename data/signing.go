package data

import (
	"github.com/donovanhide/ripple/crypto"
)

func Sign(s Signer, key crypto.Key) error {
	signingHash, err := SigningHash(s)
	if err != nil {
		return err
	}
	sig, err := key.Sign(signingHash.Bytes())
	if err != nil {
		return err
	}
	copy(s.GetPublicKey().Bytes(), key.PublicCompressed())
	*s.GetSignature() = VariableLength(sig)
	return nil
}

func CheckSignature(s Signer) (bool, error) {
	signingHash, err := SigningHash(s)
	if err != nil {
		return false, err
	}
	return crypto.Verify(s.GetPublicKey().Bytes(), s.GetSignature().Bytes(), signingHash.Bytes())
}
