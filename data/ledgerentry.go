package data

type LedgerEntrySlice []LedgerEntry

type leBase struct {
	LedgerEntryType LedgerEntryType
	Flags           *LedgerEntryFlag `json:",omitempty"`
	LedgerIndex     *Hash256         `json:"index,omitempty"`
	Hash            Hash256          `json:"-"`
	Id              Hash256          `json:"-"`
}

type AccountRoot struct {
	leBase
	Account           Account
	Sequence          uint32
	Balance           Value
	PreviousTxnID     Hash256
	PreviousTxnLgrSeq uint32
	OwnerCount        *uint32         `json:",omitempty"`
	AccountTxnID      *Hash256        `json:",omitempty"`
	RegularKey        *RegularKey     `json:",omitempty"`
	EmailHash         *Hash128        `json:",omitempty"`
	WalletLocator     *Hash256        `json:",omitempty"`
	WalletSize        *uint32         `json:",omitempty"`
	MessageKey        *VariableLength `json:",omitempty"`
	TransferRate      *uint32         `json:",omitempty"`
	Domain            *VariableLength `json:",omitempty"`
	Signers           *VariableLength `json:",omitempty"`
}

type RippleState struct {
	leBase
	LowLimit          Amount
	HighLimit         Amount
	PreviousTxnID     *Hash256   `json:",omitempty"`
	PreviousTxnLgrSeq *uint32    `json:",omitempty"`
	Balance           Amount     `json:",omitempty"`
	LowNode           *NodeIndex `json:",omitempty"`
	HighNode          *NodeIndex `json:",omitempty"`
	LowQualityIn      *uint32    `json:",omitempty"`
	LowQualityOut     *uint32    `json:",omitempty"`
	HighQualityIn     *uint32    `json:",omitempty"`
	HighQualityOut    *uint32    `json:",omitempty"`
}

type Offer struct {
	leBase
	Account           Account
	Sequence          uint32
	TakerPays         Amount
	TakerGets         Amount
	BookDirectory     Hash256
	BookNode          *NodeIndex `json:",omitempty"`
	OwnerNode         *NodeIndex `json:",omitempty"`
	PreviousTxnID     *Hash256   `json:",omitempty"`
	PreviousTxnLgrSeq *uint32    `json:",omitempty"`
	Expiration        *uint32    `json:",omitempty"`
}

type Directory struct {
	leBase
	RootIndex         Hash256
	Indexes           Vector256
	Owner             *Account   `json:",omitempty"`
	TakerPaysCurrency *Hash160   `json:",omitempty"`
	TakerPaysIssuer   *Hash160   `json:",omitempty"`
	TakerGetsCurrency *Hash160   `json:",omitempty"`
	TakerGetsIssuer   *Hash160   `json:",omitempty"`
	ExchangeRate      *NodeIndex `json:",omitempty"`
	IndexNext         *NodeIndex `json:",omitempty"`
	IndexPrevious     *NodeIndex `json:",omitempty"`
}

type LedgerHashes struct {
	leBase
	FirstLedgerSequence uint32
	LastLedgerSequence  uint32
	Hashes              Vector256
}

type Amendments struct {
	leBase
	Amendments Hash256
}

type FeeSettings struct {
	leBase
	BaseFee           Uint64Hex
	ReferenceFeeUnits uint32
	ReserveBase       uint32
	ReserveIncrement  uint32
}

func (le *leBase) GetType() string                     { return ledgerEntryNames[le.LedgerEntryType] }
func (le *leBase) GetLedgerEntryType() LedgerEntryType { return le.LedgerEntryType }
func (le *leBase) Prefix() HashPrefix                  { return HP_LEAF_NODE }
func (le *leBase) NodeType() NodeType                  { return NT_ACCOUNT_NODE }
func (le *leBase) Ledger() uint32                      { return 0 }
func (le *leBase) GetHash() *Hash256                   { return &le.Hash }
func (le *leBase) NodeId() *Hash256                    { return &le.Id }
func (le *leBase) GetLedgerIndex() *Hash256            { return le.LedgerIndex }

func (o *Offer) Ratio() *Value {
	return o.TakerPays.Ratio(o.TakerGets)
}
