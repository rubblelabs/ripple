package data

import (
	"fmt"
	"reflect"
)

// Horrible look up tables
// Could all this be one big map?

type LedgerEntryType uint16
type TransactionType uint16

const (
	ACCOUNT_ROOT  LedgerEntryType = 0x61 // 'a'
	DIRECTORY     LedgerEntryType = 0x64 // 'd'
	AMENDMENTS    LedgerEntryType = 0x66 // 'f'
	LEDGER_HASHES LedgerEntryType = 0x68 // 'h'
	OFFER         LedgerEntryType = 0x6f // 'o'
	RIPPLE_STATE  LedgerEntryType = 0x72 // 'r'
	FEE_SETTING   LedgerEntryType = 0x73 // 's'

	PAYMENT         TransactionType = 0
	ACCOUNT_SET     TransactionType = 3
	SET_REGULAR_KEY TransactionType = 5
	OFFER_CREATE    TransactionType = 7
	OFFER_CANCEL    TransactionType = 8
	TRUST_SET       TransactionType = 20
	AMENDMENT       TransactionType = 100
	SET_FEE         TransactionType = 101
)

var LedgerFactory = [...]func() Hashable{
	func() Hashable { return &Ledger{} },
}

var fieldsFactory = [...]func() interface{}{
	ACCOUNT_ROOT:  func() interface{} { return &AccountRootFields{} },
	DIRECTORY:     func() interface{} { return &DirectoryFields{} },
	AMENDMENTS:    func() interface{} { return &AmendmentsFields{} },
	LEDGER_HASHES: func() interface{} { return &LedgerHashesFields{} },
	OFFER:         func() interface{} { return &OfferFields{} },
	RIPPLE_STATE:  func() interface{} { return &RippleStateFields{} },
	FEE_SETTING:   func() interface{} { return &FeeSettingFields{} },
}

var LedgerEntryFactory = [...]func() LedgerEntry{
	ACCOUNT_ROOT:  func() LedgerEntry { return &AccountRoot{LeBase: LeBase{LedgerEntryType: ACCOUNT_ROOT}} },
	DIRECTORY:     func() LedgerEntry { return &Directory{LeBase: LeBase{LedgerEntryType: DIRECTORY}} },
	AMENDMENTS:    func() LedgerEntry { return &Amendments{LeBase: LeBase{LedgerEntryType: AMENDMENTS}} },
	LEDGER_HASHES: func() LedgerEntry { return &LedgerHashes{LeBase: LeBase{LedgerEntryType: LEDGER_HASHES}} },
	OFFER:         func() LedgerEntry { return &Offer{LeBase: LeBase{LedgerEntryType: OFFER}} },
	RIPPLE_STATE:  func() LedgerEntry { return &RippleState{LeBase: LeBase{LedgerEntryType: RIPPLE_STATE}} },
	FEE_SETTING:   func() LedgerEntry { return &FeeSetting{LeBase: LeBase{LedgerEntryType: FEE_SETTING}} },
}

var TxFactory = [...]func() Transaction{
	PAYMENT:         func() Transaction { return &Payment{TxBase: TxBase{TransactionType: PAYMENT}} },
	ACCOUNT_SET:     func() Transaction { return &AccountSet{TxBase: TxBase{TransactionType: ACCOUNT_SET}} },
	SET_REGULAR_KEY: func() Transaction { return &SetRegularKey{TxBase: TxBase{TransactionType: SET_REGULAR_KEY}} },
	OFFER_CREATE:    func() Transaction { return &OfferCreate{TxBase: TxBase{TransactionType: OFFER_CREATE}} },
	OFFER_CANCEL:    func() Transaction { return &OfferCancel{TxBase: TxBase{TransactionType: OFFER_CANCEL}} },
	TRUST_SET:       func() Transaction { return &TrustSet{TxBase: TxBase{TransactionType: TRUST_SET}} },
	AMENDMENT:       func() Transaction { return &Amendment{TxBase: TxBase{TransactionType: AMENDMENT}} },
	SET_FEE:         func() Transaction { return &SetFee{TxBase: TxBase{TransactionType: SET_FEE}} },
}

func GetTxFactoryByType(txType string) func() Transaction {
	return TxFactory[txTypes[txType]]
}

var ledgerEntryNames = [...]string{
	ACCOUNT_ROOT:  "AccountRoot",
	DIRECTORY:     "Directory",
	AMENDMENTS:    "Amendments",
	LEDGER_HASHES: "LedgerHashes",
	OFFER:         "Offer",
	RIPPLE_STATE:  "RippleState",
	FEE_SETTING:   "Fee",
}

var ledgerEntryTypes = map[string]LedgerEntryType{
	"AccountRoot":  ACCOUNT_ROOT,
	"Directory":    DIRECTORY,
	"Amendments":   AMENDMENTS,
	"LedgerHashes": LEDGER_HASHES,
	"Offer":        OFFER,
	"RippleState":  RIPPLE_STATE,
	"Fee":          FEE_SETTING,
}

var txNames = [...]string{
	PAYMENT:         "Payment",
	ACCOUNT_SET:     "AccountSet",
	SET_REGULAR_KEY: "SetRegularKey",
	OFFER_CREATE:    "OfferCreate",
	OFFER_CANCEL:    "OfferCancel",
	TRUST_SET:       "TrustSet",
	AMENDMENT:       "Amendment",
	SET_FEE:         "SetFee",
}

var txTypes = map[string]TransactionType{
	"Payment":       PAYMENT,
	"AccountSet":    ACCOUNT_SET,
	"SetRegularKey": SET_REGULAR_KEY,
	"OfferCreate":   OFFER_CREATE,
	"OfferCancel":   OFFER_CANCEL,
	"TrustSet":      TRUST_SET,
	"Amendment":     AMENDMENT,
	"SetFee":        SET_FEE,
}

var HashableTypes []string

func init() {
	HashableTypes = append(HashableTypes, []string{"LedgerMaster", "InnerNode"}...)
	for _, typ := range txNames {
		if len(typ) > 0 {
			HashableTypes = append(HashableTypes, typ)
		}
	}
	for _, typ := range ledgerEntryNames {
		if len(typ) > 0 {
			HashableTypes = append(HashableTypes, typ)
		}
	}
}

func NewHashable(typ reflect.Type) (Hashable, error) {
	if leType, ok := ledgerEntryTypes[typ.Name()]; ok {
		return LedgerEntryFactory[leType](), nil
	}
	if txType, ok := txTypes[typ.Name()]; ok {
		return TxFactory[txType](), nil
	}
	if typ.Name() == "Ledger" {
		return LedgerFactory[0](), nil
	}
	return nil, fmt.Errorf("NewHashable: Unknown type: %s ", typ.Name())
}
