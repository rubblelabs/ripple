package crypto

import "crypto"

const (
	al_ed25519   = "ed25519"
	al_secp256k1 = "ecdsa-secp256k1"
	al_sm2       = "sm2"
)

type Wallet struct {
	Algorithm string
	crypto.PrivateKey
}

func NewWalletFromString(pk, algorithm string) *Wallet {
	wa := &Wallet{}
	switch algorithm {
	case al_ed25519:
		wa.PrivateKey, _ = NewEd25519KeyFromString(pk)
		break
	case al_sm2:
		wa.PrivateKey, _ = NewSm2KeyFromString(pk)
		break
	default:
		wa.PrivateKey, _ = NewECDSAKeyFromString(pk)
		break
	}
	return wa
}
