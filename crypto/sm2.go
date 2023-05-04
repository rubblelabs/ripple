package crypto

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"github.com/tjfoc/gmsm/sm2"
	"log"
	"math/big"
)

type sm2Key struct {
	*sm2.PrivateKey
}

func newSm2Key(seed []byte) *sm2.PrivateKey {

	privKey, err := sm2.GenerateKey(bytes.NewReader(seed))
	if err != nil {
		log.Println(err)
	}
	return privKey
}

func NewSm2KeyFromString(p string) (*sm2Key, error) {
	pk, err := hex.DecodeString(p)
	if err != nil {
		return nil, err
	}

	priv := new(sm2.PrivateKey)
	priv.PublicKey.Curve = sm2.P256Sm2()
	priv.D = big.NewInt(0).SetBytes(pk)
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(priv.D.Bytes())

	r := &sm2Key{priv}
	return r, nil
}

// If seed is nil, generate a random one
func NewSm2Key(seed []byte) (*sm2Key, error) {
	if seed == nil {
		c := sm2.P256Sm2()
		seed = make([]byte, c.Params().BitSize/8+8)
		if _, err := rand.Read(seed); err != nil {
			return nil, err
		}
	}
	return &sm2Key{newSm2Key(seed)}, nil
}

func (k *sm2Key) Id(seq *uint32) []byte {
	return Sha256RipeMD160(sm2.Compress(&k.PublicKey))
}

func (k *sm2Key) Private(seq *uint32) []byte {

	return k.D.Bytes()

}

func (k *sm2Key) Public(seq *uint32) []byte {
	return sm2.Compress(&k.PublicKey)
}
