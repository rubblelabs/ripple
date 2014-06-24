package crypto

import (
	"fmt"
	"math/big"
)

type Hash interface {
	Network() HashNetwork
	Version() HashVersion
	Payload() []byte
	PayloadTrimmed() []byte
	Value() *big.Int
	String() string
	Clone() Hash
	MarshalText() ([]byte, error)
}

// First byte is the network
// Second byte is the version
// Remaining bytes are the payload
type hash []byte

func NewRippleHash(s string) (Hash, error) {
	// Special case which will deal short addresses
	switch s {
	case "0":
		return newHashFromString(ACCOUNT_ZERO, RIPPLE)
	case "1":
		return newHashFromString(ACCOUNT_ONE, RIPPLE)
	default:
		return newHashFromString(s, RIPPLE)
	}
}

// Checks hash matches expected version
func NewRippleHashCheck(s string, version HashVersion) (Hash, error) {
	hash, err := NewRippleHash(s)
	if err != nil {
		return nil, err
	}
	if hash.Version() != version {
		want := hashTypes[RIPPLE][version].Description
		got := hashTypes[RIPPLE][hash.Version()].Description
		return nil, fmt.Errorf("Bad version for: %s expected: %s got: %s ", s, want, got)
	}
	return hash, nil
}

func NewRippleAccount(b []byte) (Hash, error) {
	return newHash(b, RIPPLE, RIPPLE_ACCOUNT_ID)
}

func NewRipplePublicNode(b []byte) (Hash, error) {
	return newHash(b, RIPPLE, RIPPLE_NODE_PUBLIC)
}

func NewRipplePrivateNode(b []byte) (Hash, error) {
	return newHash(b, RIPPLE, RIPPLE_NODE_PRIVATE)
}

func NewRipplePublicAccount(b []byte) (Hash, error) {
	return newHash(b, RIPPLE, RIPPLE_ACCOUNT_PUBLIC)
}

func NewRipplePrivateAccount(b []byte) (Hash, error) {
	return newHash(b, RIPPLE, RIPPLE_ACCOUNT_PRIVATE)
}

func NewRippleFamilyGenerator(b []byte) (Hash, error) {
	return newHash(b, RIPPLE, RIPPLE_FAMILY_GENERATOR)
}

func NewRippleFamilySeed(b []byte) (Hash, error) {
	return newHash(b, RIPPLE, RIPPLE_FAMILY_SEED)
}

func GenerateFamilySeed(password string) (Hash, error) {
	return NewRippleFamilySeed(Sha512Quarter([]byte(password)))
}

func NewBitcoinAddress(b []byte) (Hash, error) {
	return newHash(b, BITCOIN, BITCOIN_ADDRESS)
}

func NewLitecoinAddress(b []byte) (Hash, error) {
	return newHash(b, LITECOIN, LITECOIN_ADDRESS)
}

func newHash(b []byte, network HashNetwork, version HashVersion) (Hash, error) {
	n := hashTypes[network][version].Payload
	if len(b) > n {
		return nil, fmt.Errorf("Hash is wrong size, expected: %d got: %d", n, len(b))
	}
	return append(hash{byte(network), byte(version)}, b...), nil
}

func newHashFromString(s string, network HashNetwork) (Hash, error) {
	decoded, err := Base58Decode(s, alphabets[network])
	if err != nil {
		return nil, err
	}
	return append(hash{byte(network)}, decoded[:len(decoded)-4]...), nil
}

func (h hash) String() string {
	b := append(hash{byte(h.Version())}, h.Payload()...)
	return Base58Encode(b, alphabets[h.Network()])
}

func (h hash) Network() HashNetwork {
	return HashNetwork(h[0])
}

func (h hash) Version() HashVersion {
	return HashVersion(h[1])
}

func (h hash) Payload() []byte {
	return h[2:]
}

// Return a slice of the payload with leading zeroes omitted
func (h hash) PayloadTrimmed() []byte {
	payload := h.Payload()
	for i := range payload {
		if payload[i] != 0 {
			return payload[i:]
		}
	}
	return payload[len(payload)-1:]
}

func (h hash) Value() *big.Int {
	return big.NewInt(0).SetBytes(h.Payload())
}

func (h hash) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

func (h hash) Clone() Hash {
	c := make(hash, len(h))
	copy(c, h)
	return c
}
