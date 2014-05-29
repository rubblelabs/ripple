package crypto

import (
	"math/big"
)

var zero = big.NewInt(0)
var one = big.NewInt(1)

type HashVersion byte
type HashNetwork byte

const (
	ACCOUNT_ZERO = "rrrrrrrrrrrrrrrrrrrrrhoLvTp"
	ACCOUNT_ONE  = "rrrrrrrrrrrrrrrrrrrrBZbvji"
	NaN          = "rrrrrrrrrrrrrrrrrrrn5RM1rHd"
	ROOT         = "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh"
)

const (
	RIPPLE   HashNetwork = 0
	BITCOIN  HashNetwork = 1
	LITECOIN HashNetwork = 2
)

const (
	RIPPLE_NODE_PUBLIC      HashVersion = 28
	RIPPLE_NODE_PRIVATE     HashVersion = 32
	RIPPLE_ACCOUNT_ID       HashVersion = 0
	RIPPLE_ACCOUNT_PUBLIC   HashVersion = 35
	RIPPLE_ACCOUNT_PRIVATE  HashVersion = 34
	RIPPLE_FAMILY_GENERATOR HashVersion = 41
	RIPPLE_FAMILY_SEED      HashVersion = 33
	BITCOIN_ADDRESS         HashVersion = 0
	LITECOIN_ADDRESS        HashVersion = 48
)

var alphabets = [3]string{
	RIPPLE:   "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz",
	BITCOIN:  "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz",
	LITECOIN: "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz",
}

var hashTypes = map[HashNetwork]map[HashVersion]struct {
	Description       string
	Prefix            byte
	Payload           int
	MaximumCharacters int
}{
	RIPPLE: {
		RIPPLE_NODE_PUBLIC:      {"Validation public key for node.", 'n', 33, 53},
		RIPPLE_NODE_PRIVATE:     {"Validation private key for node.", 'p', 32, 52},
		RIPPLE_ACCOUNT_ID:       {"Short name for sending funds to an account.", 'r', 20, 35},
		RIPPLE_ACCOUNT_PUBLIC:   {"Account public key.", 'a', 33, 53},
		RIPPLE_ACCOUNT_PRIVATE:  {"Account private key.", 'p', 32, 52},
		RIPPLE_FAMILY_GENERATOR: {"Family public generator", 'f', 33, 53},
		RIPPLE_FAMILY_SEED:      {"Family seed.", 's', 16, 29},
	},
	BITCOIN: {
		BITCOIN_ADDRESS: {"Public address", '1', 20, 35},
	},
	LITECOIN: {
		LITECOIN_ADDRESS: {"Public address", 'L', 20, 35},
	},
}
