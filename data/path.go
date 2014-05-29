package data

import (
	"encoding/json"
	"fmt"
)

type PathEntry uint8

type Path struct {
	Account  *Account
	Currency *Currency
	Issuer   *Account
}

type Paths [][]Path

func (p Path) PathEntry() PathEntry {
	var entry PathEntry
	if p.Account != nil {
		entry |= 0x01
	}
	if p.Currency != nil {
		entry |= 0x10
	}
	if p.Issuer != nil {
		entry |= 0x20
	}
	return entry
}

func (p Path) MarshalJSON() ([]byte, error) {
	typ := p.PathEntry()
	return json.Marshal(struct {
		Account  *Account  `json:"account,omitempty"`
		Currency *Currency `json:"currency,omitempty"`
		Issuer   *Account  `json:"issuer,omitempty"`
		Type     PathEntry `json:"type"`
		TypeHex  string    `json:"type_hex"`
	}{
		p.Account,
		p.Currency,
		p.Issuer,
		typ,
		fmt.Sprintf("%016X", uint64(typ)),
	})
}
