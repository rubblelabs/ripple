package data

import (
	"fmt"
	"sort"
)

type Trade struct {
	Buyer    Account
	Seller   Account
	Price    Value
	Amount   Value
	Currency Currency
	Issuer   Account
}

type Balance struct {
	Account  Account
	Balance  Value
	Change   Value
	Currency Currency
}

func (t Trade) String() string {
	return fmt.Sprintf("%s/%-34s %-34s=>%-34s %18s@%18s", t.Currency, t.Issuer, t.Seller, t.Buyer, t.Amount, t.Price)
}

func (b Balance) String() string {
	return fmt.Sprintf("Account: %-34s  Currency: %s Balance: %20s Change: %20s", b.Account, b.Currency, b.Balance, b.Change)
}

type TradeSlice []Trade
type BalanceSlice []Balance

func (s TradeSlice) Len() int           { return len(s) }
func (s TradeSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s TradeSlice) Less(i, j int) bool { return s[i].Price.Less(s[j].Price) }

func (s BalanceSlice) Len() int      { return len(s) }
func (s BalanceSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s BalanceSlice) Less(i, j int) bool {
	switch {
	case !s[i].Currency.Equals(s[j].Currency):
		return s[i].Currency.Less(s[j].Currency)
	case s[i].Change.Abs().Equals(*s[j].Change.Abs()):
		return s[i].Change.Negative != s[j].Change.Negative
	default:
		return s[i].Change.Abs().Less(*s[j].Change.Abs())
	}
}

func (s *TradeSlice) Add(buyer, seller *Account, price, amount *Amount) {
	*s = append(*s, Trade{*buyer, *seller, *price.Value, *amount.Value, price.Currency, price.Issuer})
}

func (s TradeSlice) Sum() (*Amount, error) {
	if len(s) == 0 {
		return nil, nil
	}
	return nil, nil
	// sum, err := NewAmount(fmt.Sprintf("0/%s/%s", s[0].Currency, s[0].Issuer))
	// if err != nil {
	// 	return nil, err
	// }
	// for _, trade := range s {
	// 	if sum, err = sum.Add(&trade.Amount); err != nil {
	// 		return nil, err
	// 	}
	// }
	// return sum, nil
}

func (s *BalanceSlice) Add(account *Account, balance, change *Value, currency *Currency) {
	*s = append(*s, Balance{*account, *balance, *change, *currency})
}

func (txm *TransactionWithMetaData) Trades() (TradeSlice, error) {
	var (
		trades  TradeSlice
		account = txm.Transaction.GetBase().Account
	)
	for _, node := range txm.MetaData.AffectedNodes {
		switch {
		case node.CreatedNode != nil && node.CreatedNode.LedgerEntryType == OFFER:
			// No actual side effect
		case node.DeletedNode != nil && node.DeletedNode.LedgerEntryType == OFFER:
			// An OfferCreate specifying previous OfferSequence has no side effect
			if node.DeletedNode.PreviousFields == nil {
				continue
			}
			// Fully consumed offer
			previous, final := node.DeletedNode.PreviousFields.(*OfferFields), node.DeletedNode.FinalFields.(*OfferFields)
			price, err := previous.TakerPays.Divide(previous.TakerGets)
			if err != nil {
				return nil, err
			}
			trades.Add(&account, final.Account, price, previous.TakerGets)
		case node.ModifiedNode != nil && node.ModifiedNode.LedgerEntryType == OFFER:
			// No change?
			if node.ModifiedNode.PreviousFields == nil {
				continue
			}
			// Partially consumed offer
			previous, current := node.ModifiedNode.PreviousFields.(*OfferFields), node.ModifiedNode.FinalFields.(*OfferFields)
			paid, err := previous.TakerPays.Subtract(current.TakerPays)
			if err != nil {
				return nil, err
			}
			got, err := previous.TakerGets.Subtract(current.TakerGets)
			if err != nil {
				return nil, err
			}
			price, err := paid.Divide(got)
			if err != nil {
				return nil, err
			}
			trades.Add(&account, current.Account, price, got)
		}
	}
	sort.Sort(trades)
	return trades, nil
}

func (txm *TransactionWithMetaData) Balances() (BalanceSlice, error) {
	var (
		balances BalanceSlice
		account  = txm.Transaction.GetBase().Account
	)
	for _, node := range txm.MetaData.AffectedNodes {
		switch {
		case node.CreatedNode != nil:
			switch node.CreatedNode.LedgerEntryType {
			case ACCOUNT_ROOT:
				created := node.CreatedNode.NewFields.(*AccountRootFields)
				balances.Add(created.Account, &zeroNative, created.Balance, &zeroCurrency)
			case RIPPLE_STATE:
				// New trust line
				state := node.CreatedNode.NewFields.(*RippleStateFields)
				balances.Add(&account, &zeroNonNative, state.Balance.Value, &state.Balance.Currency)
			}
		case node.DeletedNode != nil:
			switch node.DeletedNode.LedgerEntryType {
			case RIPPLE_STATE:
				//?
			case ACCOUNT_ROOT:
				return nil, fmt.Errorf("Deleted AccountRoot!")
			}
		case node.ModifiedNode != nil:
			if node.ModifiedNode.PreviousFields == nil {
				// No change
				continue
			}
			switch node.ModifiedNode.LedgerEntryType {
			case ACCOUNT_ROOT:
				// Changed XRP Balance
				previous, current := node.ModifiedNode.PreviousFields.(*AccountRootFields), node.ModifiedNode.FinalFields.(*AccountRootFields)
				if previous.Balance == nil {
					// ownercount change
					continue
				}
				change, err := NewAmount(int64(current.Balance.Num - previous.Balance.Num))
				if err != nil {
					return nil, err
				}
				balances.Add(current.Account, current.Balance, change.Value, &zeroCurrency)
			case RIPPLE_STATE:
				// Changed non-native balance
				previous, current := node.ModifiedNode.PreviousFields.(*RippleStateFields), node.ModifiedNode.FinalFields.(*RippleStateFields)
				if previous.Balance == nil {
					//flag change
					continue
				}
				change, err := current.Balance.Subtract(previous.Balance)
				if err != nil {
					return nil, err
				}
				balances.Add(&current.HighLimit.Issuer, current.Balance.Value, change.Value, &current.Balance.Currency)
				balances.Add(&current.LowLimit.Issuer, current.Balance.Value.Negate(), change.Value.Negate(), &current.Balance.Currency)
			}
		}
	}
	sort.Sort(balances)
	return balances, nil
}
