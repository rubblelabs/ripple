package data

// Horrible look up tables
// Could all this be one big map?

type LedgerEntryType uint16
type TransactionType uint16

const (
	// LedgerEntryType values come from rippled's "LedgerFormats.h"
	SIGNER_LIST      LedgerEntryType = 0x53 // 'S'
	TICKET           LedgerEntryType = 0x54 // 'T'
	ACCOUNT_ROOT     LedgerEntryType = 0x61 // 'a'
	DIRECTORY        LedgerEntryType = 0x64 // 'd'
	AMENDMENTS       LedgerEntryType = 0x66 // 'f'
	LEDGER_HASHES    LedgerEntryType = 0x68 // 'h'
	OFFER            LedgerEntryType = 0x6f // 'o'
	RIPPLE_STATE     LedgerEntryType = 0x72 // 'r'
	FEE_SETTINGS     LedgerEntryType = 0x73 // 's'
	ESCROW           LedgerEntryType = 0x75 // 'u'
	PAY_CHANNEL      LedgerEntryType = 0x78 // 'x'
	CHECK            LedgerEntryType = 0x43 // 'C'
	DEPOSIT_PRE_AUTH LedgerEntryType = 0x70 // 'p'
	NEGATIVE_UNL     LedgerEntryType = 0x4e // 'N'
	NFTOKEN_PAGE     LedgerEntryType = 0x50 // 'P'
	NFTOKEN_OFFER    LedgerEntryType = 0x37 // '7'

	// TransactionType values come from rippled's "TxFormats.h"
	PAYMENT              TransactionType = 0
	ESCROW_CREATE        TransactionType = 1
	ESCROW_FINISH        TransactionType = 2
	ACCOUNT_SET          TransactionType = 3
	ESCROW_CANCEL        TransactionType = 4
	SET_REGULAR_KEY      TransactionType = 5
	OFFER_CREATE         TransactionType = 7
	OFFER_CANCEL         TransactionType = 8
	TICKET_CREATE        TransactionType = 10
	SIGNER_LIST_SET      TransactionType = 12
	PAYCHAN_CREATE       TransactionType = 13
	PAYCHAN_FUND         TransactionType = 14
	PAYCHAN_CLAIM        TransactionType = 15
	CHECK_CREATE         TransactionType = 16
	CHECK_CASH           TransactionType = 17
	CHECK_CANCEL         TransactionType = 18
	SET_DEPOSIT_PREAUTH  TransactionType = 19
	TRUST_SET            TransactionType = 20
	ACCOUNT_DELETE       TransactionType = 21
	NFTOKEN_MINT         TransactionType = 25
	NFTOKEN_BURN         TransactionType = 26
	NFTOKEN_CREATE_OFFER TransactionType = 27
	NFTOKEN_CANCEL_OFFER TransactionType = 28
	NFTOKEN_ACCEPT_OFFER TransactionType = 29

	AMENDMENT  TransactionType = 100
	SET_FEE    TransactionType = 101
	UNL_MODIFY TransactionType = 102
)

var LedgerFactory = [...]func() Hashable{
	func() Hashable { return &Ledger{} },
}

var LedgerEntryFactory = [...]func() LedgerEntry{
	ACCOUNT_ROOT:     func() LedgerEntry { return &AccountRoot{leBase: leBase{LedgerEntryType: ACCOUNT_ROOT}} },
	DIRECTORY:        func() LedgerEntry { return &Directory{leBase: leBase{LedgerEntryType: DIRECTORY}} },
	AMENDMENTS:       func() LedgerEntry { return &Amendments{leBase: leBase{LedgerEntryType: AMENDMENTS}} },
	LEDGER_HASHES:    func() LedgerEntry { return &LedgerHashes{leBase: leBase{LedgerEntryType: LEDGER_HASHES}} },
	OFFER:            func() LedgerEntry { return &Offer{leBase: leBase{LedgerEntryType: OFFER}} },
	RIPPLE_STATE:     func() LedgerEntry { return &RippleState{leBase: leBase{LedgerEntryType: RIPPLE_STATE}} },
	FEE_SETTINGS:     func() LedgerEntry { return &FeeSettings{leBase: leBase{LedgerEntryType: FEE_SETTINGS}} },
	ESCROW:           func() LedgerEntry { return &Escrow{leBase: leBase{LedgerEntryType: ESCROW}} },
	SIGNER_LIST:      func() LedgerEntry { return &SignerList{leBase: leBase{LedgerEntryType: SIGNER_LIST}} },
	TICKET:           func() LedgerEntry { return &Ticket{leBase: leBase{LedgerEntryType: TICKET}} },
	PAY_CHANNEL:      func() LedgerEntry { return &PayChannel{leBase: leBase{LedgerEntryType: PAY_CHANNEL}} },
	CHECK:            func() LedgerEntry { return &Check{leBase: leBase{LedgerEntryType: CHECK}} },
	DEPOSIT_PRE_AUTH: func() LedgerEntry { return &DepositPreAuth{leBase: leBase{LedgerEntryType: DEPOSIT_PRE_AUTH}} },
	NEGATIVE_UNL:     func() LedgerEntry { return &NegativeUNL{leBase: leBase{LedgerEntryType: NEGATIVE_UNL}} },
	NFTOKEN_PAGE:     func() LedgerEntry { return &NFTokenPage{leBase: leBase{LedgerEntryType: NFTOKEN_PAGE}} },
	NFTOKEN_OFFER:    func() LedgerEntry { return &NFTokenOffer{leBase: leBase{LedgerEntryType: NFTOKEN_OFFER}} },
}

var TxFactory = [...]func() Transaction{
	PAYMENT:              func() Transaction { return &Payment{TxBase: TxBase{TransactionType: PAYMENT}} },
	ACCOUNT_SET:          func() Transaction { return &AccountSet{TxBase: TxBase{TransactionType: ACCOUNT_SET}} },
	ACCOUNT_DELETE:       func() Transaction { return &AccountDelete{TxBase: TxBase{TransactionType: ACCOUNT_DELETE}} },
	SET_REGULAR_KEY:      func() Transaction { return &SetRegularKey{TxBase: TxBase{TransactionType: SET_REGULAR_KEY}} },
	OFFER_CREATE:         func() Transaction { return &OfferCreate{TxBase: TxBase{TransactionType: OFFER_CREATE}} },
	OFFER_CANCEL:         func() Transaction { return &OfferCancel{TxBase: TxBase{TransactionType: OFFER_CANCEL}} },
	TRUST_SET:            func() Transaction { return &TrustSet{TxBase: TxBase{TransactionType: TRUST_SET}} },
	AMENDMENT:            func() Transaction { return &Amendment{TxBase: TxBase{TransactionType: AMENDMENT}} },
	SET_FEE:              func() Transaction { return &SetFee{TxBase: TxBase{TransactionType: SET_FEE}} },
	UNL_MODIFY:           func() Transaction { return &UNLModify{TxBase: TxBase{TransactionType: UNL_MODIFY}} },
	TICKET_CREATE:        func() Transaction { return &TicketCreate{TxBase: TxBase{TransactionType: TICKET_CREATE}} },
	ESCROW_CREATE:        func() Transaction { return &EscrowCreate{TxBase: TxBase{TransactionType: ESCROW_CREATE}} },
	ESCROW_FINISH:        func() Transaction { return &EscrowFinish{TxBase: TxBase{TransactionType: ESCROW_FINISH}} },
	ESCROW_CANCEL:        func() Transaction { return &EscrowCancel{TxBase: TxBase{TransactionType: ESCROW_CANCEL}} },
	SIGNER_LIST_SET:      func() Transaction { return &SignerListSet{TxBase: TxBase{TransactionType: SIGNER_LIST_SET}} },
	PAYCHAN_CREATE:       func() Transaction { return &PaymentChannelCreate{TxBase: TxBase{TransactionType: PAYCHAN_CREATE}} },
	PAYCHAN_FUND:         func() Transaction { return &PaymentChannelFund{TxBase: TxBase{TransactionType: PAYCHAN_FUND}} },
	PAYCHAN_CLAIM:        func() Transaction { return &PaymentChannelClaim{TxBase: TxBase{TransactionType: PAYCHAN_CLAIM}} },
	CHECK_CREATE:         func() Transaction { return &CheckCreate{TxBase: TxBase{TransactionType: CHECK_CREATE}} },
	CHECK_CASH:           func() Transaction { return &CheckCash{TxBase: TxBase{TransactionType: CHECK_CASH}} },
	CHECK_CANCEL:         func() Transaction { return &CheckCancel{TxBase: TxBase{TransactionType: CHECK_CANCEL}} },
	SET_DEPOSIT_PREAUTH:  func() Transaction { return &SetDepositPreAuth{TxBase: TxBase{TransactionType: SET_DEPOSIT_PREAUTH}} },
	NFTOKEN_MINT:         func() Transaction { return &NFTokenMint{TxBase: TxBase{TransactionType: NFTOKEN_MINT}} },
	NFTOKEN_BURN:         func() Transaction { return &NFTokenBurn{TxBase: TxBase{TransactionType: NFTOKEN_BURN}} },
	NFTOKEN_CREATE_OFFER: func() Transaction { return &NFTokenCreateOffer{TxBase: TxBase{TransactionType: NFTOKEN_CREATE_OFFER}} },
	NFTOKEN_CANCEL_OFFER: func() Transaction { return &NFTCancelOffer{TxBase: TxBase{TransactionType: NFTOKEN_CANCEL_OFFER}} },
	NFTOKEN_ACCEPT_OFFER: func() Transaction { return &NFTAcceptOffer{TxBase: TxBase{TransactionType: NFTOKEN_ACCEPT_OFFER}} },
}

var ledgerEntryNames = [...]string{
	ACCOUNT_ROOT:     "AccountRoot",
	DIRECTORY:        "DirectoryNode",
	AMENDMENTS:       "Amendments",
	LEDGER_HASHES:    "LedgerHashes",
	OFFER:            "Offer",
	RIPPLE_STATE:     "RippleState",
	FEE_SETTINGS:     "FeeSettings",
	ESCROW:           "Escrow",
	SIGNER_LIST:      "SignerList",
	TICKET:           "Ticket",
	PAY_CHANNEL:      "PayChannel",
	CHECK:            "Check",
	DEPOSIT_PRE_AUTH: "DepositPreauth",
	NEGATIVE_UNL:     "NegativeUNL",
	NFTOKEN_PAGE:     "NFTokenPage",
	NFTOKEN_OFFER:    "NFTokenOffer",
}

var ledgerEntryTypes = map[string]LedgerEntryType{
	"AccountRoot":    ACCOUNT_ROOT,
	"DirectoryNode":  DIRECTORY,
	"Amendments":     AMENDMENTS,
	"LedgerHashes":   LEDGER_HASHES,
	"Offer":          OFFER,
	"RippleState":    RIPPLE_STATE,
	"FeeSettings":    FEE_SETTINGS,
	"Escrow":         ESCROW,
	"SignerList":     SIGNER_LIST,
	"Ticket":         TICKET,
	"PayChannel":     PAY_CHANNEL,
	"Check":          CHECK,
	"DepositPreauth": DEPOSIT_PRE_AUTH,
	"NegativeUNL":    NEGATIVE_UNL,
	"NFTokenPage":    NFTOKEN_PAGE,
	"NFTokenOffer":   NFTOKEN_OFFER,
}

var txNames = [...]string{
	PAYMENT:              "Payment",
	ACCOUNT_SET:          "AccountSet",
	ACCOUNT_DELETE:       "AccountDelete",
	SET_REGULAR_KEY:      "SetRegularKey",
	OFFER_CREATE:         "OfferCreate",
	OFFER_CANCEL:         "OfferCancel",
	TRUST_SET:            "TrustSet",
	AMENDMENT:            "EnableAmendment",
	SET_FEE:              "SetFee",
	UNL_MODIFY:           "UNLModify",
	TICKET_CREATE:        "TicketCreate",
	ESCROW_CREATE:        "EscrowCreate",
	ESCROW_FINISH:        "EscrowFinish",
	ESCROW_CANCEL:        "EscrowCancel",
	SIGNER_LIST_SET:      "SignerListSet",
	PAYCHAN_CREATE:       "PaymentChannelCreate",
	PAYCHAN_FUND:         "PaymentChannelFund",
	PAYCHAN_CLAIM:        "PaymentChannelClaim",
	CHECK_CREATE:         "CheckCreate",
	CHECK_CASH:           "CheckCash",
	CHECK_CANCEL:         "CheckCancel",
	SET_DEPOSIT_PREAUTH:  "DepositPreauth",
	NFTOKEN_MINT:         "NFTokenMint",
	NFTOKEN_BURN:         "NFTokenBurn",
	NFTOKEN_CREATE_OFFER: "NFTokenCreateOffer",
	NFTOKEN_CANCEL_OFFER: "NFTokenCancelOffer",
	NFTOKEN_ACCEPT_OFFER: "NFTokenAcceptOffer",
}

var txTypes = map[string]TransactionType{
	"Payment":              PAYMENT,
	"AccountSet":           ACCOUNT_SET,
	"AccountDelete":        ACCOUNT_DELETE,
	"SetRegularKey":        SET_REGULAR_KEY,
	"OfferCreate":          OFFER_CREATE,
	"OfferCancel":          OFFER_CANCEL,
	"TrustSet":             TRUST_SET,
	"EnableAmendment":      AMENDMENT,
	"SetFee":               SET_FEE,
	"UNLModify":            UNL_MODIFY,
	"TicketCreate":         TICKET_CREATE,
	"EscrowCreate":         ESCROW_CREATE,
	"EscrowFinish":         ESCROW_FINISH,
	"EscrowCancel":         ESCROW_CANCEL,
	"SignerListSet":        SIGNER_LIST_SET,
	"PaymentChannelCreate": PAYCHAN_CREATE,
	"PaymentChannelFund":   PAYCHAN_FUND,
	"PaymentChannelClaim":  PAYCHAN_CLAIM,
	"CheckCreate":          CHECK_CREATE,
	"CheckCash":            CHECK_CASH,
	"CheckCancel":          CHECK_CANCEL,
	"DepositPreauth":       SET_DEPOSIT_PREAUTH,
	"NFTokenMint":          NFTOKEN_MINT,
	"NFTokenBurn":          NFTOKEN_BURN,
	"NFTokenCreateOffer":   NFTOKEN_CREATE_OFFER,
	"NFTokenCancelOffer":   NFTOKEN_CANCEL_OFFER,
	"NFTokenAcceptOffer":   NFTOKEN_ACCEPT_OFFER,
}

var HashableTypes []string

func init() {
	HashableTypes = append(HashableTypes, NT_TRANSACTION_NODE.String())
	for _, typ := range txNames {
		if len(typ) > 0 {
			HashableTypes = append(HashableTypes, typ)
		}
	}
	HashableTypes = append(HashableTypes, NT_ACCOUNT_NODE.String())
	for _, typ := range ledgerEntryNames {
		if len(typ) > 0 {
			HashableTypes = append(HashableTypes, typ)
		}
	}
}

func (t TransactionType) String() string {
	return txNames[t]
}

func (le LedgerEntryType) String() string {
	return ledgerEntryNames[le]
}

func GetTxFactoryByType(txType string) func() Transaction {
	return TxFactory[txTypes[txType]]
}

func GetLedgerEntryFactoryByType(leType string) func() LedgerEntry {
	return LedgerEntryFactory[ledgerEntryTypes[leType]]
}
