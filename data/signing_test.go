package data

import (
	"reflect"
	"testing"

	"github.com/rubblelabs/ripple/crypto"
)

func TestMultiSignWithVerification(t *testing.T) {
	// generate signers seeds
	seed1 := genSeedFromPassword(t, "password1")
	seed2 := genSeedFromPassword(t, "password2")
	seq := uint32(0)
	key1 := seed1.Key(ECDSA)
	account1 := seed1.AccountId(ECDSA, &seq)
	key2 := seed2.Key(ECDSA)
	account2 := seed2.AccountId(ECDSA, &seq)

	// prepare first signature
	tx := buildPaymentTxForTheMultiSigning(t)
	if err := MultiSign(tx, key1, &seq, account1); err != nil {
		t.Fatal(err)
	}
	signer1 := Signer{
		Signer: SignerItem{
			Account:       account1,
			TxnSignature:  tx.TxnSignature,
			SigningPubKey: tx.SigningPubKey,
		},
	}
	// prepare second signature
	tx = buildPaymentTxForTheMultiSigning(t)
	if err := MultiSign(tx, key2, &seq, account2); err != nil {
		t.Fatal(err)
	}
	signer2 := Signer{
		Signer: SignerItem{
			Account:       account2,
			TxnSignature:  tx.TxnSignature,
			SigningPubKey: tx.SigningPubKey,
		},
	}
	// rebuild the tx and set signers
	tx = buildPaymentTxForTheMultiSigning(t)
	if err := SetSigners(tx, signer1, signer2); err != nil {
		t.Fatal(err)
	}
	// check that signature is valid
	valid, invalidSigners, err := CheckMultiSignature(tx)
	if err != nil {
		t.Fatalf("Failed to check the signature, err:%s", err)
	}
	if !valid {
		t.Fatal("Unexpected invalid signatures")
	}
	if len(invalidSigners) != 0 {
		t.Fatal("Unexpected invalid signers length")
	}

	// update one signature
	tx.Signers[0].Signer.TxnSignature = tx.Signers[1].Signer.TxnSignature
	valid, invalidSigners, err = CheckMultiSignature(tx)
	if err != nil {
		t.Fatalf("Failed to check the signature, err:%s", err)
	}
	if valid {
		t.Fatal("Unexpected valid signatures")
	}
	if len(invalidSigners) != 1 {
		t.Fatal("Unexpected invalid signers length")
	}
	if !reflect.DeepEqual(invalidSigners[0], tx.Signers[0]) {
		t.Fatalf("Unexpected signer, expected:%+v, got:%+v", invalidSigners[0], tx.Signers[0])
	}

	// update tx data to check that both signers are invalid now
	tx.Sequence = 123
	valid, invalidSigners, err = CheckMultiSignature(tx)
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("Unexpected valid signatures")
	}
	if len(invalidSigners) != 2 {
		t.Fatal("Unexpected invalid signers length")
	}
}

func buildPaymentTxForTheMultiSigning(t *testing.T) *Payment {
	amount, err := NewAmount("1")
	if err != nil {
		t.Fatal(err)
	}
	tx := Payment{
		Amount: *amount,
		TxBase: TxBase{
			Account:         zeroAccount,
			Sequence:        1,
			TransactionType: PAYMENT,
		},
	}
	// important for the multi-signing
	tx.TxBase.SigningPubKey = &PublicKey{}
	return &tx
}

func genSeedFromPassword(t *testing.T, password string) *Seed {
	seedFromPass, err := crypto.GenerateFamilySeed(password)
	if err != nil {
		t.Fatal(err)
	}
	seed, err := NewSeedFromAddress(seedFromPass.String())
	if err != nil {
		t.Fatal(err)
	}

	return seed
}
