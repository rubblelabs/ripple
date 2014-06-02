package data

import (
	"fmt"
)

type Currency [20]byte

var zeroCurrency Currency

func NewCurrency(s string) (Currency, error) {
	if s == "XRP" {
		return zeroCurrency, nil
	}
	var currency Currency
	if len(s) != 3 {
		return currency, fmt.Errorf("Bad Currency: %s", s)
	}
	copy(currency[12:], []byte(s))
	return currency, nil
}

func (c Currency) IsNative() bool {
	return c == zeroCurrency
}

func (c Currency) Equals(other Currency) bool {
	return c == other
}

func (c Currency) Clone() Currency {
	var n Currency
	copy(n[:], c[:])
	return n
}

func (c *Currency) Bytes() []byte {
	if c != nil {
		return c[:]
	}
	return []byte(nil)
}

func (c Currency) String() string {
	if c.IsNative() {
		return "XRP"
	}
	return string(c[12:15])
}

func (c Currency) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}
