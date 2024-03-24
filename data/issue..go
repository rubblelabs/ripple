package data

import "fmt"

type Issue struct {
	Currency Currency `json:"currency"`
	Issuer   Account  `json:"issuer,omitempty"`
}

func (i Issue) String() string {
	if i.Currency.IsNative() {
		return i.Currency.String()
	}
	return fmt.Sprintf("%s/%s", i.Currency, i.Issuer)
}
