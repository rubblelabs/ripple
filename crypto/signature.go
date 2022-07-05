package crypto

import (
	"crypto/ed25519"
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
)

func Sign(privateKey, hash, msg []byte) ([]byte, error) {
	switch len(privateKey) {
	case ed25519.PrivateKeySize:
		return signEd25519(privateKey, msg), nil
	case secp256k1.PrivKeyBytesLen:
		return signECDSA(privateKey, hash), nil
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
	default:
		return false, fmt.Errorf("Unknown public key format")
	}
}

func signEd25519(privateKey, msg []byte) []byte {
	return ed25519.Sign(privateKey, msg)[:]
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
func signECDSA(privateKey, hash []byte) []byte {
	key := secp256k1.PrivKeyFromBytes(privateKey)
	return ecdsa.Sign(key, hash).Serialize()
}

// Verifies a hash using DER encoded signature
func verifyECDSA(pubKey, signature, hash []byte) (bool, error) {
	sig, err := ecdsa.ParseDERSignature(signature)
	if err != nil {
		return false, err
	}
	pk, err := secp256k1.ParsePubKey(pubKey)
	if err != nil {
		return false, nil
	}
	return sig.Verify(hash, pk), nil
}
