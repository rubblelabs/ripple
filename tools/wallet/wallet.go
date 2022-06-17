package wallet

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/maybeTomorrow/ripple/crypto"
	"log"
)

type Wallet struct {
	PrivateKey *btcec.PrivateKey
	PublicKey  *btcec.PublicKey
}

func NewWallet(pk []byte) *Wallet {
	w := &Wallet{}
	pvk, pub := btcec.PrivKeyFromBytes(btcec.S256(), pk)
	w.PrivateKey = pvk
	w.PublicKey = pub
	return w
}

func NewWalletFromSeed(pk []byte) *Wallet {
	w := &Wallet{}
	pvk, pub := btcec.PrivKeyFromBytes(btcec.S256(), pk)
	w.PrivateKey = pvk
	w.PublicKey = pub
	return w
}

func (e *Wallet) ClassicAddress() string {

	hs, err := crypto.NewAccountId(crypto.Sha256RipeMD160(e.PublicKey.SerializeCompressed()))
	if err != nil {
		log.Println(err)
	}
	return hs.String()
}

func (e *Wallet) Sign(msg []byte) (*btcec.Signature, error) {
	return e.PrivateKey.Sign(msg)
}
