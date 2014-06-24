package crypto

import (
	"crypto/sha512"
	"math/big"
)

var maxNonce = big.NewInt(0).SetUint64(1 << 23)

type ProofOfWork struct {
	Challenge  *big.Int
	Target     *big.Int
	Iterations int
	first      []byte
	second     []byte
}

func (pow *ProofOfWork) Next(nonce []byte) *big.Int {
	first := make([]byte, 96)
	copy(first, pow.Challenge.Bytes())
	copy(first[64-len(nonce):], nonce)
	hasher := sha512.New()
	for i := pow.Iterations - 1; i >= 0; i-- {
		hasher.Write(first)
		copy(first[64:], hasher.Sum(nil)[:32])
		copy(pow.second[i*32:], first[64:])
		hasher.Reset()
	}
	hasher.Write(pow.second)
	return big.NewInt(0).SetBytes(hasher.Sum(nil)[:32])
}

func NewProofOfWork(challenge, target []byte, iterations uint32) *ProofOfWork {
	return &ProofOfWork{
		Challenge:  big.NewInt(0).SetBytes(challenge),
		Target:     big.NewInt(0).SetBytes(target),
		Iterations: int(iterations),
		second:     make([]byte, iterations*32),
	}
}

func (pow *ProofOfWork) Solve() []byte {
	for nonce := big.NewInt(0); ; nonce.Add(nonce, one) {
		result := pow.Next(nonce.Bytes())
		switch {
		case nonce.Cmp(maxNonce) >= 0:
			return []byte(nil)
		case result.Cmp(pow.Target) <= 0:
			return nonce.Bytes()
		}
	}
}

func (pow *ProofOfWork) Check(nonce []byte) bool {
	return pow.Next(nonce).Cmp(pow.Target) <= 0
}
