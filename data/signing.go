package data

import "github.com/rubblelabs/ripple/crypto"

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
	//
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
	s.SetSigners(signers)

	hash, _, err := Raw(s)
	if err != nil {
		return err
	}
	copy(s.GetHash().Bytes(), hash.Bytes())
	return nil
}
