package data

// Horrible look up tables
// Could all this be one big map?

type LedgerEntryType uint16
type TransactionType uint16

const (
	SIGNER_LIST   LedgerEntryType = 0x53 // 'S'
	TICKET        LedgerEntryType = 0x54 // 'T'
	ACCOUNT_ROOT  LedgerEntryType = 0x61 // 'a'
	DIRECTORY     LedgerEntryType = 0x64 // 'd'
	AMENDMENTS    LedgerEntryType = 0x66 // 'f'
	LEDGER_HASHES LedgerEntryType = 0x68 // 'h'
	OFFER         LedgerEntryType = 0x6f // 'o'
	RIPPLE_STATE  LedgerEntryType = 0x72 // 'r'
	FEE_SETTINGS  LedgerEntryType = 0x73 // 's'
	SUS_PAY       LedgerEntryType = 0x75 // 'u'
	PAY_CHANNEL   LedgerEntryType = 0x78 // 'x'

	PAYMENT         TransactionType = 0
	SUS_PAY_CREATE  TransactionType = 1
	SUS_PAY_FINISH  TransactionType = 2
	ACCOUNT_SET     TransactionType = 3
	SUS_PAY_CANCEL  TransactionType = 4
	SET_REGULAR_KEY TransactionType = 5
	OFFER_CREATE    TransactionType = 7
	OFFER_CANCEL    TransactionType = 8
	TICKET_CREATE   TransactionType = 10
	TICKET_CANCEL   TransactionType = 11
	SIGNER_LIST_SET TransactionType = 12
	TRUST_SET       TransactionType = 20
	AMENDMENT       TransactionType = 100
	SET_FEE         TransactionType = 101
)

var LedgerFactory = [...]func() Hashable{
	func() Hashable { return &Ledger{} },
}

var LedgerEntryFactory = [...]func() LedgerEntry{
	ACCOUNT_ROOT:  func() LedgerEntry { return &AccountRoot{leBase: leBase{LedgerEntryType: ACCOUNT_ROOT}} },
	DIRECTORY:     func() LedgerEntry { return &Directory{leBase: leBase{LedgerEntryType: DIRECTORY}} },
	AMENDMENTS:    func() LedgerEntry { return &Amendments{leBase: leBase{LedgerEntryType: AMENDMENTS}} },
	LEDGER_HASHES: func() LedgerEntry { return &LedgerHashes{leBase: leBase{LedgerEntryType: LEDGER_HASHES}} },
	OFFER:         func() LedgerEntry { return &Offer{leBase: leBase{LedgerEntryType: OFFER}} },
	RIPPLE_STATE:  func() LedgerEntry { return &RippleState{leBase: leBase{LedgerEntryType: RIPPLE_STATE}} },
	FEE_SETTINGS:  func() LedgerEntry { return &FeeSettings{leBase: leBase{LedgerEntryType: FEE_SETTINGS}} },
	SUS_PAY:       func() LedgerEntry { return &SuspendedPayment{leBase: leBase{LedgerEntryType: SUS_PAY}} },
	SIGNER_LIST:   func() LedgerEntry { return &SignerList{leBase: leBase{LedgerEntryType: SIGNER_LIST}} },
	TICKET:        func() LedgerEntry { return &Ticket{leBase: leBase{LedgerEntryType: TICKET}} },
	PAY_CHANNEL:   func() LedgerEntry { return &PayChannel{leBase: leBase{LedgerEntryType: PAY_CHANNEL}} },
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
	SUS_PAY_CREATE:  func() Transaction { return &SuspendedPaymentCreate{TxBase: TxBase{TransactionType: SUS_PAY_CREATE}} },
	SUS_PAY_FINISH:  func() Transaction { return &SuspendedPaymentFinish{TxBase: TxBase{TransactionType: SUS_PAY_FINISH}} },
	SUS_PAY_CANCEL:  func() Transaction { return &SuspendedPaymentCancel{TxBase: TxBase{TransactionType: SUS_PAY_CANCEL}} },
	SIGNER_LIST_SET: func() Transaction { return &SignerListSet{TxBase: TxBase{TransactionType: SIGNER_LIST_SET}} },
}

var ledgerEntryNames = [...]string{
	ACCOUNT_ROOT:  "AccountRoot",
	DIRECTORY:     "DirectoryNode",
	AMENDMENTS:    "Amendments",
	LEDGER_HASHES: "LedgerHashes",
	OFFER:         "Offer",
	RIPPLE_STATE:  "RippleState",
	FEE_SETTINGS:  "FeeSettings",
	SUS_PAY:       "SuspendedPayment",
	SIGNER_LIST:   "SignerList",
	TICKET:        "Ticket",
	PAY_CHANNEL:   "PaymentChannel",
}

var ledgerEntryTypes = map[string]LedgerEntryType{
	"AccountRoot":      ACCOUNT_ROOT,
	"DirectoryNode":    DIRECTORY,
	"Amendments":       AMENDMENTS,
	"LedgerHashes":     LEDGER_HASHES,
	"Offer":            OFFER,
	"RippleState":      RIPPLE_STATE,
	"FeeSettings":      FEE_SETTINGS,
	"SuspendedPayment": SUS_PAY,
	"SignerList":       SIGNER_LIST,
	"Ticket":           TICKET,
	"PaymentChannel":   PAY_CHANNEL,
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
	SUS_PAY_CREATE:  "SuspendedPaymentCreate",
	SUS_PAY_FINISH:  "SuspendedPaymentFinish",
	SUS_PAY_CANCEL:  "SuspendedPaymentCancel",
	SIGNER_LIST_SET: "SignerListSet",
}

var txTypes = map[string]TransactionType{
	"Payment":                PAYMENT,
	"AccountSet":             ACCOUNT_SET,
	"SetRegularKey":          SET_REGULAR_KEY,
	"OfferCreate":            OFFER_CREATE,
	"OfferCancel":            OFFER_CANCEL,
	"TrustSet":               TRUST_SET,
	"Amendment":              AMENDMENT,
	"SetFee":                 SET_FEE,
	"SuspendedPaymentCreate": SUS_PAY_CREATE,
	"SuspendedPaymentFinish": SUS_PAY_FINISH,
	"SuspendedPaymentCancel": SUS_PAY_CANCEL,
	"SignerListSet":          SIGNER_LIST_SET,
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
