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
	CounterParty Account
	Balance      Value
	Change       Value
	Currency     Currency
	AMMID        *Hash256
}

func (b Balance) String() string {
	return fmt.Sprintf("CounterParty: %-34s  Currency: %s Balance: %20s Change: %20s AMMID: %v", b.CounterParty, b.Currency, b.Balance, b.Change, b.AMMID)
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

func (s *BalanceSlice) Add(counterparty *Account, balance, change *Value, currency *Currency, AMMID *Hash256) {
	if change == nil || currency == nil {
		return
	}
	*s = append(*s, Balance{*counterparty, *balance, *change, *currency, AMMID})
}

type BalanceMap map[Account]*BalanceSlice

func (m *BalanceMap) Add(account *Account, counterparty *Account, balance, change *Value, currency *Currency, AMMID *Hash256) {
	_, ok := (*m)[*account]
	if !ok {
		(*m)[*account] = &BalanceSlice{}
	}
	(*m)[*account].Add(counterparty, balance, change, currency, AMMID)
}

var txBalanceAcceptedMap = map[TransactionType]bool{OFFER_CREATE: true, PAYMENT: true, AMM_DEPOSIT: true, AMM_WITHDRAW: true, AMM_CREATE: true}

func (txm *TransactionWithMetaData) Balances() (BalanceMap, error) {
	if !txBalanceAcceptedMap[txm.GetTransactionType()] {
		return nil, nil
	}
	balanceMap := BalanceMap{}
	account := txm.Transaction.GetBase().Account
	for _, node := range txm.MetaData.AffectedNodes {
		switch {
		case node.CreatedNode != nil:
			switch node.CreatedNode.LedgerEntryType {
			case ACCOUNT_ROOT:
				created := node.CreatedNode.NewFields.(*AccountRoot)
				balanceMap.Add(created.Account, &zeroAccount, created.Balance, created.Balance, &zeroCurrency, created.AMMID)
			case RIPPLE_STATE:
				// New trust line
				state := node.CreatedNode.NewFields.(*RippleState)
				balanceMap.Add(&state.LowLimit.Issuer, &state.HighLimit.Issuer, state.Balance.Value, state.Balance.Value, &state.Balance.Currency, nil)
				balanceMap.Add(&state.HighLimit.Issuer, &state.LowLimit.Issuer, state.Balance.Value.Negate(), state.Balance.Value.Negate(), &state.Balance.Currency, nil)
			case AMMROOT:
				state := node.CreatedNode.NewFields.(*AMM)
				balanceMap.Add(&state.LPTokenBalance.Issuer,
					&state.LPTokenBalance.Issuer,
					state.LPTokenBalance.Value,
					state.LPTokenBalance.Value,
					&state.LPTokenBalance.Currency,
					nil)

			}
		case node.DeletedNode != nil:
			switch node.DeletedNode.LedgerEntryType {
			case RIPPLE_STATE:
				// A deletion (complete termination of a token) can lead to having to use
				// the deleted node to determine the balance of the counterparty.
				var previous, current *RippleState
				if node.DeletedNode.PreviousFields == nil || node.DeletedNode.FinalFields == nil {
					continue
				}
				previous = node.DeletedNode.PreviousFields.(*RippleState)
				current = node.DeletedNode.FinalFields.(*RippleState)

				if previous.Balance == nil {
					//flag change
					continue
				}
				change, err := current.Balance.Subtract(previous.Balance)
				if err != nil {
					return nil, err
				}
				balanceMap.Add(&current.LowLimit.Issuer, &current.HighLimit.Issuer, current.Balance.Value, change.Value, &current.Balance.Currency, nil)
				balanceMap.Add(&current.HighLimit.Issuer, &current.LowLimit.Issuer, current.Balance.Value.Negate(), change.Value.Negate(), &current.Balance.Currency, nil)
			case ACCOUNT_ROOT:
				return nil, fmt.Errorf("Deleted AccountRoot!")
			case AMMROOT:
				// 20230911, NvN: Looks like it is not required to handle this case
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
					balanceMap.Add(current.Account, &zeroAccount, current.Balance, change.Value, &zeroCurrency, current.AMMID)
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
				balanceMap.Add(&current.LowLimit.Issuer, &current.HighLimit.Issuer, current.Balance.Value, change.Value, &current.Balance.Currency, nil)
				balanceMap.Add(&current.HighLimit.Issuer, &current.LowLimit.Issuer, current.Balance.Value.Negate(), change.Value.Negate(), &current.Balance.Currency, nil)
			case AMMROOT:
				// Used to retrieve the currency & issuer for lptokens while parsing balances
				var (
					previous = node.ModifiedNode.PreviousFields.(*AMM)
					current  = node.ModifiedNode.FinalFields.(*AMM)
				)
				if previous.LPTokenBalance == nil {
					//Vote change, bid slot change
					continue
				}
				change, err := current.LPTokenBalance.Subtract(previous.LPTokenBalance)
				if err != nil {
					return nil, err
				}
				balanceMap.Add(&current.LPTokenBalance.Issuer, &current.LPTokenBalance.Issuer, current.LPTokenBalance.Value, change.Value, &current.LPTokenBalance.Currency, nil)
			}
		}
	}
	for _, balances := range balanceMap {
		sort.Sort(balances)
	}
	return balanceMap, nil
}
