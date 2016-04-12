package data

import (
	"fmt"
	"strings"
)

type Asset struct {
	Currency string `json:"currency"`
	Issuer   string `json:"issuer,omitempty"`
}

func NewAsset(s string) (*Asset, error) {
	if s == "XRP" {
		return &Asset{
			Currency: s,
		}, nil
	}
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("bad asset: %s", s)
	}
	return &Asset{
		Currency: parts[0],
		Issuer:   parts[1],
	}, nil
}

func (a *Asset) IsNative() bool {
	return a.Currency == "XRP"
}

func (a Asset) String() string {
	if a.IsNative() {
		return a.Currency
	}
	return a.Currency + "/" + a.Issuer
}

type OrderBookOffer struct {
	Offer
	OwnerFunds      Value  `json:"owner_funds"`
	Quality         Value  `json:"quality"`
	TakerGetsFunded Amount `json:"taker_gets_funded"`
	TakerPaysFunded Amount `json:"taker_pays_funded"`
}

type AccountLine struct {
	Account    Account  `json:"account"`
	Balance    Value    `json:"balance"`
	Currency   Currency `json:"currency"`
	Limit      Value    `json:"limit"`
	LimitPeer  Value    `json:"limit_peer"`
	QualityIn  uint32   `json:"quality_in"`
	QualityOut uint32   `json:"quality_out"`
}
