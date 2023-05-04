package data

import (
	"encoding/hex"
	"github.com/maybeTomorrow/ripple/crypto"
	"log"
)

func Sign(s Signable, key crypto.Key, sequence *uint32) error {
	s.InitialiseForSigning()
	copy(s.GetPublicKey().Bytes(), key.Public(sequence))
	hash, msg, err := SigningHash(s)

	asd := hex.EncodeToString(msg)
	log.Println("msg:", asd)
	log.Println("hashis :" + hash.String())

	if err != nil {
		return err
	}
	sig, err := crypto.Sign(key, hash.Bytes(), append(s.SigningPrefix().Bytes(), msg...))
	if err != nil {
		return err
	}
	*s.GetSignature() = VariableLength(sig)
	hash, _, err = Raw(s)
	if err != nil {
		return err
	}
	copy(s.GetHash().Bytes(), hash.Bytes())
	return nil
}

func SignMulti(s Signable, keys []crypto.Key, sequence *uint32) ([]Signer, error) {
	s.InitialiseForSigning()
	//copy(s.GetPublicKey().Bytes(), key.Public(sequence))
	hash, msg, err := SigningHash(s)

	asd := hex.EncodeToString(msg)
	log.Println(asd)

	list := make([]Signer, 0)

	for _, key := range keys {
		sig, err := crypto.Sign(key, hash.Bytes(), append(s.SigningPrefix().Bytes(), msg...))
		if err != nil {
			return nil, err
		}
		sg := Signer{}
		pb := &PublicKey{}
		ts := VariableLength(sig)
		copy(pb.Bytes(), key.Public(sequence))
		sg.SigningPubKey = pb
		sg.TxnSignature = &ts
		ai, _ := crypto.NewAccountId(key.Id(nil))
		f, _ := NewAccountFromAddress(ai.String())
		sg.Account = *f
		list = append(list, sg)
	}

	s.SetSingers(list)

	hash, _, err = Raw(s)
	if err != nil {
		return list, err
	}
	copy(s.GetHash().Bytes(), hash.Bytes())
	return list, nil
}

func CheckSignature(s Signable) (bool, error) {
	hash, msg, err := SigningHash(s)
	if err != nil {
		return false, err
	}
	return crypto.Verify(s.GetPublicKey().Bytes(), hash.Bytes(), msg, s.GetSignature().Bytes())
}
