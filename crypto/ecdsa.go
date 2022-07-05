package crypto

import (
	"crypto/rand"
	"encoding/binary"
	"math/big"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

const (
	PubKeyBytesLenCompressed = 33
)

var (
	order = secp256k1.S256().N
	zero  = big.NewInt(0)
	one   = big.NewInt(1)
)

type ecdsaKey struct {
	*secp256k1.PrivateKey
}

func newKey(seed []byte) *secp256k1.PrivateKey {
	inc := big.NewInt(0).SetBytes(seed)
	inc.Lsh(inc, 32)
	for key := big.NewInt(0); ; inc.Add(inc, one) {
		key.SetBytes(Sha512Half(inc.Bytes()))
		if key.Cmp(zero) > 0 && key.Cmp(order) < 0 {
			return secp256k1.PrivKeyFromBytes(key.Bytes())
		}
	}
}

// If seed is nil, generate a random one
func NewECDSAKey(seed []byte) (*ecdsaKey, error) {
	if seed == nil {
		seed = make([]byte, 16)
		if _, err := rand.Read(seed); err != nil {
			return nil, err
		}
	}
	return &ecdsaKey{newKey(seed)}, nil
}

func (k *ecdsaKey) generateKey(sequence uint32) *secp256k1.PrivateKey {
	seed := make([]byte, PubKeyBytesLenCompressed+4)
	copy(seed, k.PubKey().SerializeCompressed())
	binary.BigEndian.PutUint32(seed[PubKeyBytesLenCompressed:], sequence)
	key := newKey(seed).ToECDSA()
	key.D.Add(key.D, k.ToECDSA().D).Mod(key.D, order)
	return secp256k1.PrivKeyFromBytes(key.D.Bytes())
}

func (k *ecdsaKey) Id(sequence *uint32) []byte {
	if sequence == nil {
		return Sha256RipeMD160(k.PubKey().SerializeCompressed())
	}
	return Sha256RipeMD160(k.Public(sequence))
}

func (k *ecdsaKey) Private(sequence *uint32) []byte {
	if sequence == nil {
		return k.ToECDSA().D.Bytes()
	}
	return k.generateKey(*sequence).ToECDSA().D.Bytes()
}

func (k *ecdsaKey) Public(sequence *uint32) []byte {
	if sequence == nil {
		return k.PubKey().SerializeCompressed()
	}
	return k.generateKey(*sequence).PubKey().SerializeCompressed()
}
