package data

type TransactionFlag uint32
type LedgerEntryFlag uint32

// Transaction Flags
const (
	//Universal flags
	TxCanonicalSignature TransactionFlag = 0x80000000

	// Payment flags
	TxNoDirectRipple TransactionFlag = 0x00010000
	TxPartialPayment TransactionFlag = 0x00020000
	TxLimitQuality   TransactionFlag = 0x00040000
	TxCircle         TransactionFlag = 0x00080000 // Not implemented

	// AccountSet flags
	TxSetRequireDest   TransactionFlag = 0x00000001
	TxSetRequireAuth   TransactionFlag = 0x00000002
	TxSetDisallowXRP   TransactionFlag = 0x00000003
	TxSetDisableMaster TransactionFlag = 0x00000004
	TxRequireDestTag   TransactionFlag = 0x00010000
	TxOptionalDestTag  TransactionFlag = 0x00020000
	TxRequireAuth      TransactionFlag = 0x00040000
	TxDisallowXRP      TransactionFlag = 0x00100000
	TxAllowXRP         TransactionFlag = 0x00200000

	// OfferCreate flags
	TxPassive           TransactionFlag = 0x00010000
	TxImmediateOrCancel TransactionFlag = 0x00020000
	TxFillOrKill        TransactionFlag = 0x00040000
	TxSell              TransactionFlag = 0x00080000

	// TrustSet flags
	TxSetAuth       TransactionFlag = 0x00010000
	TxSetNoRipple   TransactionFlag = 0x00020000
	TxClearNoRipple TransactionFlag = 0x00040000
)

// Ledger entry flags
const (
	// AccountRoot flags
	LsPasswordSpent  LedgerEntryFlag = 0x00010000
	LsRequireDestTag LedgerEntryFlag = 0x00020000
	LsRequireAuth    LedgerEntryFlag = 0x00040000
	LsDisallowXRP    LedgerEntryFlag = 0x00080000
	LsDisableMaster  LedgerEntryFlag = 0x00100000

	// Offer flags
	LsPassive LedgerEntryFlag = 0x00010000
	LsSell    LedgerEntryFlag = 0x00020000

	// RippleState flags
	LsLowReserve   LedgerEntryFlag = 0x00010000
	LsHighReserve  LedgerEntryFlag = 0x00020000
	LsLowAuth      LedgerEntryFlag = 0x00040000
	LsHighAuth     LedgerEntryFlag = 0x00080000
	LsLowNoRipple  LedgerEntryFlag = 0x00100000
	LsHighNoRipple LedgerEntryFlag = 0x00200000
)
