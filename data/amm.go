package data

import "fmt"

var txAmmAcceptedMap = map[TransactionType]bool{AMM_DEPOSIT: true, AMM_WITHDRAW: true, AMM_CREATE: true, AMM_VOTE: true, AMM_BID: true, PAYMENT: true}

func (txm *TransactionWithMetaData) AMM() (*AMM, error) {
	if !txAmmAcceptedMap[txm.GetTransactionType()] {
		return nil, nil
	}
	for _, nodeAffect := range txm.MetaData.AffectedNodes {
		switch {
		case nodeAffect.CreatedNode != nil && nodeAffect.CreatedNode.LedgerEntryType == AMMROOT:
			ammParsed, ok := nodeAffect.CreatedNode.NewFields.(*AMM)
			if ok {
				if nodeAffect.CreatedNode.LedgerIndex != nil {
					ammParsed.LedgerIndex = nodeAffect.CreatedNode.LedgerIndex
				}
				return ammParsed, nil
			}
		case nodeAffect.ModifiedNode != nil && nodeAffect.ModifiedNode.LedgerEntryType == AMMROOT:
			ammParsed, ok := nodeAffect.ModifiedNode.FinalFields.(*AMM)
			if ok {
				if nodeAffect.ModifiedNode.LedgerIndex != nil {
					ammParsed.LedgerIndex = nodeAffect.ModifiedNode.LedgerIndex
				}
				return ammParsed, nil
			}
		}
	}
	return nil, fmt.Errorf("AMM not found")
}
