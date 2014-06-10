package data

type LedgerEntrySlice []LedgerEntry

type Index uint64

type leBase struct {
	hashable
	LedgerEntryType LedgerEntryType
	LedgerIndex     *Hash256 `json:",omitempty"`
}

type AccountRootFields struct {
	Flags             *LedgerEntryFlag `json:",omitempty"`
	Account           *Account         `json:",omitempty"`
	Sequence          *uint32          `json:",omitempty"`
	Balance           *Value           `json:",omitempty"`
	OwnerCount        *uint32          `json:",omitempty"`
	PreviousTxnID     *Hash256         `json:",omitempty"`
	PreviousTxnLgrSeq *uint32          `json:",omitempty"`
	AccountTxnID      *Hash256         `json:",omitempty"`
	RegularKey        *RegularKey      `json:",omitempty"`
	EmailHash         *Hash128         `json:",omitempty"`
	WalletLocator     *Hash256         `json:",omitempty"`
	WalletSize        *uint32          `json:",omitempty"`
	MessageKey        *PublicKey       `json:",omitempty"`
	TransferRate      *uint32          `json:",omitempty"`
	Domain            *VariableLength  `json:",omitempty"`
	Signers           *VariableLength  `json:",omitempty"`
}

type AccountRoot struct {
	leBase
	AccountRootFields
}

type RippleStateFields struct {
	Flags             *LedgerEntryFlag `json:",omitempty"`
	LowLimit          *Amount          `json:",omitempty"`
	HighLimit         *Amount          `json:",omitempty"`
	PreviousTxnID     *Hash256         `json:",omitempty"`
	PreviousTxnLgrSeq *uint32          `json:",omitempty"`
	Balance           *Amount          `json:",omitempty"`
	LowNode           *Index           `json:",omitempty"`
	HighNode          *Index           `json:",omitempty"`
	LowQualityIn      *uint32          `json:",omitempty"`
	LowQualityOut     *uint32          `json:",omitempty"`
	HighQualityIn     *uint32          `json:",omitempty"`
	HighQualityOut    *uint32          `json:",omitempty"`
}

type RippleState struct {
	leBase
	RippleStateFields
}

type OfferFields struct {
	Flags             *LedgerEntryFlag `json:",omitempty"`
	Account           *Account         `json:",omitempty"`
	Sequence          *uint32          `json:",omitempty"`
	TakerPays         *Amount          `json:",omitempty"`
	TakerGets         *Amount          `json:",omitempty"`
	BookDirectory     *Hash256         `json:",omitempty"`
	BookNode          *Index           `json:",omitempty"`
	OwnerNode         *Index           `json:",omitempty"`
	PreviousTxnID     *Hash256         `json:",omitempty"`
	PreviousTxnLgrSeq *uint32          `json:",omitempty"`
	Expiration        *RippleTime      `json:",omitempty"`
}

type Offer struct {
	leBase
	OfferFields
}

type DirectoryFields struct {
	Flags             *LedgerEntryFlag `json:",omitempty"`
	RootIndex         *Hash256         `json:",omitempty"`
	Indexes           *Vector256       `json:",omitempty"`
	Owner             *Account         `json:",omitempty"`
	TakerPaysCurrency *Hash160         `json:",omitempty"`
	TakerPaysIssuer   *Hash160         `json:",omitempty"`
	TakerGetsCurrency *Hash160         `json:",omitempty"`
	TakerGetsIssuer   *Hash160         `json:",omitempty"`
	ExchangeRate      *Index           `json:",omitempty"`
	IndexNext         *Index           `json:",omitempty"`
	IndexPrevious     *Index           `json:",omitempty"`
}

type Directory struct {
	leBase
	DirectoryFields
}

type LedgerHashesFields struct {
	Flags               *LedgerEntryFlag `json:",omitempty"`
	FirstLedgerSequence uint32
	LastLedgerSequence  uint32
	Hashes              Vector256
}

type LedgerHashes struct {
	leBase
	LedgerHashesFields
}

type AmendmentsFields struct {
	Flags      *LedgerEntryFlag `json:",omitempty"`
	Amendments Hash256
}

type Amendments struct {
	leBase
	AmendmentsFields
}

type FeeSettingFields struct {
	Flags             *LedgerEntryFlag `json:",omitempty"`
	BaseFee           uint64
	ReferenceFeeUnits uint32
	ReserveBase       uint32
	ReserveIncrement  uint32
}

type FeeSetting struct {
	leBase
	FeeSettingFields
}

func (le *leBase) GetType() string {
	return ledgerEntryNames[le.LedgerEntryType]
}

func (le *leBase) GetLedgerEntryType() LedgerEntryType {
	return le.LedgerEntryType
}
