package data

type Query struct {
	MinLedger *uint64
	MaxLedger *uint64
	Account   *Account
	Limit     *uint64
	Order     []string
	Instance  interface{}
}
