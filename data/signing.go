package data

import (
	"fmt"
	"sort"

	"github.com/rubblelabs/ripple/crypto"
)

func Sign(s Signable, key crypto.Key, sequence *uint32) error {
	s.InitialiseForSigning()
	copy(s.GetPublicKey().Bytes(), key.Public(sequence))
	hash, msg, err := SigningHash(s)
	if err != nil {
		return err
	}
	sig, err := crypto.Sign(key.Private(sequence), hash.Bytes(), append(s.SigningPrefix().Bytes(), msg...))
	if err != nil {
		return err
	}
	*s.GetSignature() = sig
	hash, _, err = Raw(s)
	if err != nil {
		return err
	}
	copy(s.GetHash().Bytes(), hash.Bytes())
	return nil
}

func CheckSignature(s Signable) (bool, error) {
	hash, msg, err := SigningHash(s)
	if err != nil {
		return false, err
	}
	return crypto.Verify(s.GetPublicKey().Bytes(), hash.Bytes(), msg, s.GetSignature().Bytes())
}

func MultiSign(s MultiSignable, key crypto.Key, sequence *uint32, account Account) error {
	s.InitialiseForSigning()
	hash, msg, err := MultiSigningHash(s, account)
	if err != nil {
		return err
	}
	msg = append(s.MultiSigningPrefix().Bytes(), msg...)
	msg = append(msg, account.Bytes()...)

	sig, err := crypto.Sign(key.Private(sequence), hash.Bytes(), msg)
	if err != nil {
		return err
	}
	*s.GetSignature() = sig
	// copy pub key only after the signing
	copy(s.GetPublicKey().Bytes(), key.Public(sequence))

	return nil
}

func SetSigners(s MultiSignable, signers ...Signer) error {
	sort.Slice(signers, func(i, j int) bool {
		return signers[i].Signer.Account.Less(signers[j].Signer.Account)
	})
	s.SetSigners(signers)

	hash, _, err := Raw(s)
	if err != nil {
		return err
	}
	copy(s.GetHash().Bytes(), hash.Bytes())
	return nil
}

func CheckMultiSignature(s MultiSignable) (bool, []Signer, error) {
	if len(s.GetSigners()) == 0 {
		return false, nil, fmt.Errorf("no signers in the multi-signable transaction")
	}
	signers := s.GetSigners()
	invalidSigners := make([]Signer, 0)
	for _, signer := range signers {
		account := signer.Signer.Account
		pubKey := signer.Signer.SigningPubKey
		signature := signer.Signer.TxnSignature

		hash, msg, err := MultiSigningHash(s, account)
		if err != nil {
			return false, nil, err
		}
		msg = append(s.MultiSigningPrefix().Bytes(), msg...)
		msg = append(msg, account.Bytes()...)

		valid, err := crypto.Verify(pubKey.Bytes(), hash.Bytes(), msg, signature.Bytes())
		if err != nil {
			return false, nil, err
		}
		if !valid {
			invalidSigners = append(invalidSigners, signer)
		}
	}

	return len(invalidSigners) == 0, invalidSigners, nil
}
