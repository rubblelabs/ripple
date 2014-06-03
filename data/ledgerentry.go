package data

type LeCommon struct {
	Flags             *uint32  `json:",omitempty"`
	LedgerIndex       *Hash256 `json:",omitempty"`
	PreviousTxnID     *Hash256 `json:",omitempty"`
	PreviousTxnLgrSeq *uint32  `json:",omitempty"`
}

type LeBase struct {
	hashable
	LedgerEntryType LedgerEntryType `json:",omitempty"`
}

type AccountRootFields struct {
	LeCommon
	Account       *Account        `json:",omitempty"`
	Sequence      *uint32         `json:",omitempty"`
	Balance       *Value          `json:",omitempty"`
	AccountTxnID  *Hash256        `json:",omitempty"`
	OwnerCount    *uint32         `json:",omitempty"`
	RegularKey    *RegularKey     `json:",omitempty"`
	EmailHash     *Hash128        `json:",omitempty"`
	WalletLocator *Hash256        `json:",omitempty"`
	WalletSize    *uint32         `json:",omitempty"`
	MessageKey    *PublicKey      `json:",omitempty"`
	TransferRate  *uint32         `json:",omitempty"`
	Domain        *VariableLength `json:",omitempty"`
	Signers       *VariableLength `json:",omitempty"`
}

type AccountRoot struct {
	LeBase
	AccountRootFields
}

type RippleStateFields struct {
	LeCommon
	LowLimit       *Amount `json:",omitempty"`
	HighLimit      *Amount `json:",omitempty"`
	Balance        *Amount `json:",omitempty"`
	LowNode        *uint64 `json:",omitempty"`
	HighNode       *uint64 `json:",omitempty"`
	LowQualityIn   *uint32 `json:",omitempty"`
	LowQualityOut  *uint32 `json:",omitempty"`
	HighQualityIn  *uint32 `json:",omitempty"`
	HighQualityOut *uint32 `json:",omitempty"`
}

type RippleState struct {
	LeBase
	RippleStateFields
}

type OfferFields struct {
	LeCommon
	Account       *Account `json:",omitempty"`
	Sequence      *uint32  `json:",omitempty"`
	TakerPays     *Amount  `json:",omitempty"`
	TakerGets     *Amount  `json:",omitempty"`
	BookDirectory *Hash256 `json:",omitempty"`
	BookNode      *uint64  `json:",omitempty"`
	OwnerNode     *uint64  `json:",omitempty"`
	Expiration    *uint32  `json:",omitempty"`
}

type Offer struct {
	LeBase
	OfferFields
}

type DirectoryFields struct {
	Flags             *uint32    `json:",omitempty"`
	RootIndex         *Hash256   `json:",omitempty"`
	Indexes           *Vector256 `json:",omitempty"`
	Owner             *Account   `json:",omitempty"`
	TakerPaysCurrency *Hash160   `json:",omitempty"`
	TakerPaysIssuer   *Hash160   `json:",omitempty"`
	TakerGetsCurrency *Hash160   `json:",omitempty"`
	TakerGetsIssuer   *Hash160   `json:",omitempty"`
	ExchangeRate      *uint64    `json:",omitempty"`
	IndexNext         *uint64    `json:",omitempty"`
	IndexPrevious     *uint64    `json:",omitempty"`
}

type Directory struct {
	LeBase
	DirectoryFields
}

type LedgerHashesFields struct {
	LeCommon
	FirstLedgerSequence uint32
	LastLedgerSequence  uint32
	Hashes              Vector256
}

type LedgerHashes struct {
	LeBase
	LedgerHashesFields
}

type AmendmentsFields struct {
	LeCommon
	Amendments Hash256
}

type Amendments struct {
	LeBase
	AmendmentsFields
}

type FeeSettingFields struct {
	LeCommon
	BaseFee           uint64
	ReferenceFeeUnits uint32
	ReserveBase       uint32
	ReserveIncrement  uint32
}

type FeeSetting struct {
	LeBase
	FeeSettingFields
}

func (le *LeBase) GetType() string {
	return ledgerEntryNames[le.LedgerEntryType]
}

func (le *LeBase) GetLedgerEntryType() LedgerEntryType {
	return le.LedgerEntryType
}
