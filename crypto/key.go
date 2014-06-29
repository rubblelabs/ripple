package crypto

import (
	"crypto/rand"
	"github.com/conformal/btcec"
	"math/big"
)

type Key interface {
	Sign(b []byte) ([]byte, error)
	PublicCompressed() []byte
	PrivateBytes() []byte
}

type baseKey struct {
	priv btcec.PrivateKey
}

// Returns DER encoded signature from input hash
func (k *baseKey) Sign(hash []byte) ([]byte, error) {
	sig, err := k.priv.Sign(hash)
	if err != nil {
		return nil, err
	}
	return sig.Serialize(), nil
}

// Verifies a hash using DER encoded signature
func Verify(pubKey, signature, hash []byte) (bool, error) {
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

func (k *baseKey) PrivateBytes() []byte {
	return k.priv.D.Bytes()
}

func (k *baseKey) PublicCompressed() []byte {
	return (*btcec.PublicKey)(&k.priv.PublicKey).SerializeCompressed()
}

type RootDeterministicKey struct {
	baseKey
	Seed Hash
}

type AccountKey struct {
	baseKey
}

var order = btcec.S256().N

type genFunc func(*big.Int) bool

func newKey(priv, inc *big.Int, f genFunc) *baseKey {
	pk := big.NewInt(0).Set(priv)
	for ; f(pk); inc.Add(inc, one) {
		pk.SetBytes(Sha512Half(inc.Bytes()))
	}
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), pk.Bytes())
	return &baseKey{*privKey}
}

func ParsePublicKey(b []byte) (*btcec.PublicKey, error) {
	key, err := btcec.ParsePubKey(b, btcec.S256())
	return (*btcec.PublicKey)(key), err
}

func ParsePublicKeyFromHash(hash []byte) (*btcec.PublicKey, error) {
	h, err := NewRippleHash(string(hash))
	if err != nil {
		return nil, err
	}
	return ParsePublicKey(h.Payload())
}

// If seed is nil, generate a random one
func GenerateRootDeterministicKey(seed []byte) (*RootDeterministicKey, error) {
	if seed == nil {
		seed = make([]byte, 16)
		if _, err := rand.Read(seed); err != nil {
			return nil, err
		}
	}
	s, err := NewRippleFamilySeed(seed)
	if err != nil {
		return nil, err
	}
	inc := big.NewInt(0).SetBytes(seed)
	inc.Lsh(inc, 32)
	f := func(priv *big.Int) bool { return priv.Cmp(order) >= 0 }
	key := newKey(order, inc, f)
	key.priv.X, key.priv.Y = key.priv.ScalarBaseMult(key.PrivateBytes())
	return &RootDeterministicKey{
		baseKey: *key,
		Seed:    s,
	}, nil
}

func (r *RootDeterministicKey) GenerateAccountId(sequence int32) (Hash, error) {
	inc := big.NewInt(0).SetBytes(r.PublicCompressed())
	inc.Lsh(inc, 32).Add(inc, big.NewInt(int64(sequence))).Lsh(inc, 32).Add(inc, zero)
	f := func(priv *big.Int) bool { return priv.Cmp(order) >= 0 }
	key := newKey(order, inc, f)
	key.priv.X, key.priv.Y = key.priv.ScalarBaseMult(key.PrivateBytes())
	key.priv.X, key.priv.Y = key.priv.Add(key.priv.X, key.priv.Y, r.priv.X, r.priv.Y)
	b := Sha256RipeMD160(key.PublicCompressed())
	return NewRippleAccount(b)
}

func (r *RootDeterministicKey) GenerateAccountKey(sequence int32) (*AccountKey, error) {
	generator := big.NewInt(0).SetBytes(r.PublicCompressed())
	inc := big.NewInt(0).Set(generator)
	inc.Lsh(inc, 32).Add(inc, big.NewInt(int64(sequence))).Lsh(inc, 32).Add(inc, zero)
	f := func(priv *big.Int) bool { return priv.Cmp(generator) >= 0 || priv.Cmp(zero) <= 0 }
	key := newKey(zero, inc, f)
	key.priv.D.Add(key.priv.D, r.priv.D).Mod(key.priv.D, order)
	key.priv.X, key.priv.Y = key.priv.Curve.ScalarBaseMult(key.PrivateBytes())
	return &AccountKey{baseKey: *key}, nil
}

func (r *RootDeterministicKey) PublicGenerator() (Hash, error) {
	return NewRippleFamilyGenerator(r.PublicCompressed())
}

func (r *RootDeterministicKey) PublicNodeKey() (Hash, error) {
	return NewRipplePublicNode(r.PublicCompressed())
}

func (r *RootDeterministicKey) PrivateNodeKey() (Hash, error) {
	return NewRipplePrivateNode(r.PrivateBytes())
}

func (a *AccountKey) PublicAccountKey() (Hash, error) {
	return NewRipplePublicAccount(a.PublicCompressed())
}

func (a *AccountKey) PrivateAccountKey() (Hash, error) {
	return NewRipplePrivateAccount(a.PrivateBytes())
}
