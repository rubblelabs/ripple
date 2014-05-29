package crypto

import (
	"code.google.com/p/go.crypto/ripemd160"
	"crypto/sha256"
	"crypto/sha512"
)

func sha512Section(b []byte, n int) ([]byte, error) {
	hasher := sha512.New()
	if _, err := hasher.Write(b); err != nil {
		return nil, err
	}
	return hasher.Sum(nil)[:n], nil
}

// Returns first 32 bytes of a SHA512 of the input bytes
func Sha512Half(b []byte) ([]byte, error) {
	return sha512Section(b, 32)
}

// Returns first 16 bytes of a SHA512 of the input bytes
func Sha512Quarter(b []byte) ([]byte, error) {
	return sha512Section(b, 16)
}

func DoubleSha256(b []byte) ([]byte, error) {
	hasher := sha256.New()
	if _, err := hasher.Write(b); err != nil {
		return nil, err
	}
	sha := hasher.Sum(nil)
	hasher.Reset()
	if _, err := hasher.Write(sha); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func Sha256RipeMD160(b []byte) ([]byte, error) {
	ripe := ripemd160.New()
	sha := sha256.New()
	if _, err := sha.Write(b); err != nil {
		return nil, err
	}
	if _, err := ripe.Write(sha.Sum(nil)); err != nil {
		return nil, err
	}
	return ripe.Sum(nil), nil
}
