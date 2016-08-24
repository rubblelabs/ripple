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
	OwnerFunds      Value          `json:"owner_funds"`
	Quality         NonNativeValue `json:"quality"`
	TakerGetsFunded Amount         `json:"taker_gets_funded"`
	TakerPaysFunded Amount         `json:"taker_pays_funded"`
}

type AccountOffer struct {
	Flags     uint32         `json:"flags"`
	Quality   NonNativeValue `json:"quality"`
	Sequence  uint32         `json:"seq"`
	TakerGets Amount         `json:"taker_gets"`
	TakerPays Amount         `json:"taker_pays"`
}

type AccountOfferSlice []AccountOffer

func (s AccountOfferSlice) Len() int           { return len(s) }
func (s AccountOfferSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s AccountOfferSlice) Less(i, j int) bool { return s[i].Sequence < s[j].Sequence }

type AccountLine struct {
	Account      Account        `json:"account"`
	Balance      Value          `json:"balance"`
	Currency     Currency       `json:"currency"`
	Limit        NonNativeValue `json:"limit"`
	LimitPeer    NonNativeValue `json:"limit_peer"`
	NoRipple     bool           `json:"no_ripple"`
	NoRipplePeer bool           `json:"no_ripple_peer"`
	QualityIn    uint32         `json:"quality_in"`
	QualityOut   uint32         `json:"quality_out"`
}

type AccountLineSlice []AccountLine

func (s AccountLineSlice) Len() int      { return len(s) }
func (s AccountLineSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s AccountLineSlice) Less(i, j int) bool {
	if s[i].Currency.Equals(s[j].Currency) {
		return !s[i].Balance.Abs().Less(*s[j].Balance.Abs())
	}
	return s[i].Currency.Less(s[j].Currency)
}
