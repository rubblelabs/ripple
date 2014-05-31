package ledger

import (
	"fmt"
	"github.com/donovanhide/ripple/data"
	"github.com/donovanhide/ripple/storage"
	"io"
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

func (state *LedgerState) WriteSummary(w io.Writer) error {
	summary := make(map[string]uint64)
	if err := state.AccountState.Summary(summary); err != nil {
		return err
	}
	if err := state.Transactions.Summary(summary); err != nil {
		return err
	}
	out := []string{strconv.FormatUint(uint64(state.Sequence()), 10)}
	for _, typ := range data.HashableTypes {
		out = append(out, strconv.FormatUint(summary[typ], 10))
	}
	if _, err := w.Write([]byte(strings.Join(out, ","))); err != nil {
		return err
	}
	_, err := w.Write([]byte{'\n'})
	return err
}
