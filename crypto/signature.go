package crypto

import (
	"crypto/ed25519"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/tjfoc/gmsm/sm2"
	"math/big"
)

func Sign(privateKey Key, hash, msg []byte) ([]byte, error) {
	switch privateKey.Public(nil)[0] {
	case 0xED:
		return signEd25519(privateKey.Private(nil), msg)
	case 0x02, 0x03:
		return signECDSA(privateKey.Private(nil), hash)
	case 0x01, 0x00:
		return signSm2(privateKey.Private(nil), msg)
	default:
		return nil, fmt.Errorf("Unknown private key format")
	}
}

func Verify(publicKey, hash, msg, signature []byte) (bool, error) {
	switch publicKey[0] {
	case 0xED:
		return verifyEd25519(publicKey, signature, msg)
	case 0x02, 0x03:
		return verifyECDSA(publicKey, signature, hash)
	case 0x01, 0x00:
		return verifySm2(publicKey, signature, hash)
	default:
		return false, fmt.Errorf("Unknown public key format")
	}
}

func signEd25519(privateKey, msg []byte) ([]byte, error) {
	return ed25519.Sign(privateKey, msg)[:], nil
}

func verifyEd25519(pubKey, signature, msg []byte) (bool, error) {
	switch {
	case len(pubKey) != ed25519.PublicKeySize+1:
		return false, fmt.Errorf("Wrong public key length: %d", len(pubKey))
	case pubKey[0] != 0xED:
		return false, fmt.Errorf("Wrong public format:")
	case len(signature) != ed25519.SignatureSize:
		return false, fmt.Errorf("Wrong Signature length: %d", len(signature))
	default:
		return ed25519.Verify(pubKey[1:], msg, signature), nil
	}
}

// Returns DER encoded signature from input hash
func signECDSA(privateKey, hash []byte) ([]byte, error) {
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKey)
	sig, err := priv.Sign(hash)
	if err != nil {
		return nil, err
	}
	return sig.Serialize(), nil
}

// Verifies a hash using DER encoded signature
func verifyECDSA(pubKey, signature, hash []byte) (bool, error) {
	sig, err := btcec.ParseDERSignature(signature, btcec.S256())
	if err != nil {
		return false, err
	}
	pk, err := btcec.ParsePubKey(pubKey, btcec.S256())
	if err != nil {
		return false, nil
	}
	return sig.Verify(hash, pk), nil
}

func signSm2(privateKey, msg []byte) ([]byte, error) {
	priv := new(sm2.PrivateKey)
	priv.PublicKey.Curve = sm2.P256Sm2()
	priv.D = big.NewInt(0).SetBytes(privateKey)
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(priv.D.Bytes())
	return priv.Sign(nil, msg, nil)
}

// Verifies a hash using DER encoded signature
func verifySm2(pubKey, signature, msg []byte) (bool, error) {

	pub := sm2.Decompress(pubKey)

	return pub.Verify(msg, signature), nil

}
