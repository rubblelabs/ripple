package data

import (
	"fmt"
	"sort"
)

// Transfer is a directional representation of a RippleState or AccountRoot balance change.
// Payments and OfferCreates lead to the creation of zero or more Transfers.
//
// 	TransitFee is earned by the Issuer
// 	QualityIn and QualityOut are earned by the Liquidity Provider and can be negative.
//
// Four scenarios:
// 	1. XRP -> XRP
// 	2. XRP -> IOU/Issuer 			Requires an orderbook
// 	3. IOU/Issuer -> XRP			Requires an orderbook
// 	4. IOU/IssuerA <-> IOU/IssuerB		Also known as Rippling, requires an account which trusts both currency/issuer pairs
type Transfer struct {
	Source             Account
	Destination        Account
	SourceBalance      Amount
	DestinationBalance Amount
	Change             Value
	TransitFee         *Value // Applies to all transfers except XRP -> XRP
	QualityIn          *Value // Applies to IOU -> IOU transfers
	QualityOut         *Value // Applies to IOU -> IOU transfers
}

type Balance struct {
	Account  Account
	Balance  Value
	Change   Value
	Currency Currency
}

func (b Balance) String() string {
	return fmt.Sprintf("Account: %-34s  Currency: %s Balance: %20s Change: %20s", b.Account, b.Currency, b.Balance, b.Change)
}

type BalanceSlice []Balance

func (s BalanceSlice) Len() int      { return len(s) }
func (s BalanceSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s BalanceSlice) Less(i, j int) bool {
	switch {
	case !s[i].Currency.Equals(s[j].Currency):
		return s[i].Currency.Less(s[j].Currency)
	case s[i].Change.Abs().Equals(*s[j].Change.Abs()):
		return s[i].Change.negative != s[j].Change.negative
	default:
		return s[i].Change.Abs().Less(*s[j].Change.Abs())
	}
}

func (s *BalanceSlice) Add(account *Account, balance, change *Value, currency *Currency) {
	*s = append(*s, Balance{*account, *balance, *change, *currency})
}

func (txm *TransactionWithMetaData) Balances() (BalanceSlice, error) {
	if txm.GetTransactionType() != OFFER_CREATE && txm.GetTransactionType() != PAYMENT {
		return nil, nil
	}
	var (
		balances BalanceSlice
		account  = txm.Transaction.GetBase().Account
	)
	for _, node := range txm.MetaData.AffectedNodes {
		switch {
		case node.CreatedNode != nil:
			switch node.CreatedNode.LedgerEntryType {
			case ACCOUNT_ROOT:
				created := node.CreatedNode.NewFields.(*AccountRoot)
				balances.Add(created.Account, &zeroNative, created.Balance, &zeroCurrency)
			case RIPPLE_STATE:
				// New trust line
				state := node.CreatedNode.NewFields.(*RippleState)
				balances.Add(&state.HighLimit.Issuer, &zeroNonNative, state.Balance.Value, &state.Balance.Currency)
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
				var (
					previous = node.ModifiedNode.PreviousFields.(*AccountRoot)
					current  = node.ModifiedNode.FinalFields.(*AccountRoot)
				)
				if previous.Balance == nil {
					// ownercount change
					continue
				}
				change, err := NewAmount(int64(current.Balance.num - previous.Balance.num))
				if err != nil {
					return nil, err
				}
				// Add fee and see if change is non-zero
				if current.Account.Equals(account) {
					change.Value, err = change.Value.Add(txm.GetBase().Fee)
					if err != nil {
						return nil, err
					}
				}
				if change.num != 0 {
					balances.Add(current.Account, current.Balance, change.Value, &zeroCurrency)
				}
			case RIPPLE_STATE:
				// Changed non-native balance
				var (
					previous = node.ModifiedNode.PreviousFields.(*RippleState)
					current  = node.ModifiedNode.FinalFields.(*RippleState)
				)
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
