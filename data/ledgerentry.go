package data

type LedgerEntrySlice []LedgerEntry

type leBase struct {
	LedgerEntryType   LedgerEntryType
	LedgerIndex       *Hash256 `json:"index,omitempty"`
	PreviousTxnID     *Hash256 `json:",omitempty"`
	PreviousTxnLgrSeq *uint32  `json:",omitempty"`
	Hash              Hash256  `json:"-"`
	Id                Hash256  `json:"-"`
}

type AccountRoot struct {
	leBase
	Flags         *LedgerEntryFlag `json:",omitempty"`
	Account       *Account         `json:",omitempty"`
	Sequence      *uint32          `json:",omitempty"`
	Balance       *Value           `json:",omitempty"`
	OwnerCount    *uint32          `json:",omitempty"`
	AccountTxnID  *Hash256         `json:",omitempty"`
	RegularKey    *RegularKey      `json:",omitempty"`
	EmailHash     *Hash128         `json:",omitempty"`
	WalletLocator *Hash256         `json:",omitempty"`
	WalletSize    *uint32          `json:",omitempty"`
	MessageKey    *VariableLength  `json:",omitempty"`
	TransferRate  *uint32          `json:",omitempty"`
	Domain        *VariableLength  `json:",omitempty"`
	Signers       *VariableLength  `json:",omitempty"`
}

type RippleState struct {
	leBase
	Flags          *LedgerEntryFlag `json:",omitempty"`
	LowLimit       *Amount          `json:",omitempty"`
	HighLimit      *Amount          `json:",omitempty"`
	Balance        *Amount          `json:",omitempty"`
	LowNode        *NodeIndex       `json:",omitempty"`
	HighNode       *NodeIndex       `json:",omitempty"`
	LowQualityIn   *uint32          `json:",omitempty"`
	LowQualityOut  *uint32          `json:",omitempty"`
	HighQualityIn  *uint32          `json:",omitempty"`
	HighQualityOut *uint32          `json:",omitempty"`
}

type Offer struct {
	leBase
	Flags         *LedgerEntryFlag `json:",omitempty"`
	Account       *Account         `json:",omitempty"`
	Sequence      *uint32          `json:",omitempty"`
	TakerPays     *Amount          `json:",omitempty"`
	TakerGets     *Amount          `json:",omitempty"`
	BookDirectory *Hash256         `json:",omitempty"`
	BookNode      *NodeIndex       `json:",omitempty"`
	OwnerNode     *NodeIndex       `json:",omitempty"`
	Expiration    *uint32          `json:",omitempty"`
}

type Directory struct {
	leBase
	Flags             *LedgerEntryFlag `json:",omitempty"`
	RootIndex         *Hash256         `json:",omitempty"`
	Indexes           *Vector256       `json:",omitempty"`
	Owner             *Account         `json:",omitempty"`
	TakerPaysCurrency *Hash160         `json:",omitempty"`
	TakerPaysIssuer   *Hash160         `json:",omitempty"`
	TakerGetsCurrency *Hash160         `json:",omitempty"`
	TakerGetsIssuer   *Hash160         `json:",omitempty"`
	ExchangeRate      *ExchangeRate    `json:",omitempty"`
	IndexNext         *NodeIndex       `json:",omitempty"`
	IndexPrevious     *NodeIndex       `json:",omitempty"`
}

type LedgerHashes struct {
	leBase
	Flags               *LedgerEntryFlag `json:",omitempty"`
	FirstLedgerSequence *uint32          `json:",omitempty"`
	LastLedgerSequence  *uint32          `json:",omitempty"`
	Hashes              *Vector256       `json:",omitempty"`
}

type Amendments struct {
	leBase
	Flags      *LedgerEntryFlag `json:",omitempty"`
	Amendments *Hash256         `json:",omitempty"`
}

type FeeSettings struct {
	leBase
	Flags             *LedgerEntryFlag `json:",omitempty"`
	BaseFee           *Uint64Hex       `json:",omitempty"`
	ReferenceFeeUnits *uint32          `json:",omitempty"`
	ReserveBase       *uint32          `json:",omitempty"`
	ReserveIncrement  *uint32          `json:",omitempty"`
}

func (le *leBase) GetType() string                     { return ledgerEntryNames[le.LedgerEntryType] }
func (le *leBase) GetLedgerEntryType() LedgerEntryType { return le.LedgerEntryType }
func (le *leBase) Prefix() HashPrefix                  { return HP_LEAF_NODE }
func (le *leBase) NodeType() NodeType                  { return NT_ACCOUNT_NODE }
func (le *leBase) Ledger() uint32                      { return 0 }
func (le *leBase) GetHash() *Hash256                   { return &le.Hash }
func (le *leBase) NodeId() *Hash256                    { return &le.Id }
func (le *leBase) GetLedgerIndex() *Hash256            { return le.LedgerIndex }
func (le *leBase) GetPreviousTxnId() *Hash256          { return le.PreviousTxnID }

func (o *Offer) Ratio() *Value {
	return o.TakerPays.Ratio(*o.TakerGets)
}
