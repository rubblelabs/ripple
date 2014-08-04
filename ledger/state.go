package ledger

import (
	"fmt"
	"github.com/rubblelabs/ripple/data"
	"github.com/rubblelabs/ripple/storage"
	"strconv"
	"strings"
)

type AccountState struct {
	AccountRoot data.AccountRoot
	RippleState []data.RippleState
}

type CurrencyPair struct {
	Left  data.Currency
	Right data.Currency
}

type Offers struct {
	Asks []data.Offer
	Bids []data.Offer
}

type LedgerState struct {
	*data.Ledger
	AccountState *RadixMap
	Transactions *RadixMap
	Books        map[CurrencyPair]Offers
	full         bool
}

func NewEmptyLedgerState(sequence uint32) *LedgerState {
	return &LedgerState{
		Ledger:       data.NewEmptyLedger(sequence),
		AccountState: NewEmptyRadixMap(),
		Transactions: NewEmptyRadixMap(),
		Books:        make(map[CurrencyPair]Offers),
	}
}

func NewLedgerStateFromDB(hash data.Hash256, db storage.DB) (*LedgerState, error) {
	node, err := db.Get(hash)
	if err != nil {
		return nil, err
	}
	ledger, ok := node.(*data.Ledger)
	if !ok {
		return nil, fmt.Errorf("NewLedgerStateFromDB: not a ledger:%s", hash.String())
	}
	return &LedgerState{
		Ledger:       ledger,
		AccountState: NewRadixMap(ledger.StateHash, db),
		Transactions: NewRadixMap(ledger.TransactionHash, db),
		Books:        make(map[CurrencyPair]Offers),
	}, nil
}

func (state *LedgerState) Sequence() uint32 {
	return state.Ledger.LedgerSequence
}

func (state *LedgerState) Fill() error {
	if err := state.AccountState.Fill(); err != nil {
		return err
	}
	return state.Transactions.Fill()
}

func (state *LedgerState) Summary() (string, error) {
	summary := make(map[string]uint64)
	var s []string
	if err := state.AccountState.Summary(summary); err != nil {
		return "", err
	}
	if err := state.Transactions.Summary(summary); err != nil {
		return "", err
	}
	for _, typ := range data.HashableTypes {
		s = append(s, strconv.FormatUint(summary[typ], 10))
	}
	return strings.Join(s, ","), nil
}
