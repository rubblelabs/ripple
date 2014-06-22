package data

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"strings"
)

type pathEntry uint8

const (
	PATH_BOUNDARY pathEntry = 0xFF
	PATH_END      pathEntry = 0x00

	PATH_ACCOUNT  pathEntry = 0x01
	PATH_REDEEM   pathEntry = 0x02
	PATH_CURRENCY pathEntry = 0x10
	PATH_ISSUER   pathEntry = 0x20
)

type Path struct {
	Account  *Account
	Currency *Currency
	Issuer   *Account
}

type PathSet []Paths
type Paths []Path

func (p Path) pathEntry() pathEntry {
	var entry pathEntry
	if p.Account != nil {
		entry |= PATH_ACCOUNT
	}
	if p.Currency != nil {
		entry |= PATH_CURRENCY
	}
	if p.Issuer != nil {
		entry |= PATH_ISSUER
	}
	return entry
}

func (p Paths) Signature() (uint32, error) {
	checksum := crc32.NewIEEE()
	for _, path := range p {
		b := append(path.Account.Bytes(), append(path.Currency.Bytes(), path.Issuer.Bytes()...)...)
		if _, err := checksum.Write(b); err != nil {
			return 0, err
		}
	}
	return checksum.Sum32(), nil
}

func (p Paths) String() string {
	var s []string
	for _, path := range p {
		s = append(s, path.String())
	}
	return strings.Join(s, " => ")
}

func (p Path) String() string {
	var s []string
	if p.Account != nil {
		s = append(s, p.Account.String())
	}
	if p.Currency != nil {
		s = append(s, p.Currency.String())
	}
	if p.Issuer != nil {
		s = append(s, p.Issuer.String())
	}
	return strings.Join(s, "/")
}

func (p Path) MarshalJSON() ([]byte, error) {
	typ := p.pathEntry()
	return json.Marshal(struct {
		Account  *Account  `json:"account,omitempty"`
		Currency *Currency `json:"currency,omitempty"`
		Issuer   *Account  `json:"issuer,omitempty"`
		Type     pathEntry `json:"type"`
		TypeHex  string    `json:"type_hex"`
	}{
		p.Account,
		p.Currency,
		p.Issuer,
		typ,
		fmt.Sprintf("%016X", uint64(typ)),
	})
}
