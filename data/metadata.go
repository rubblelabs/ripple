package data

type Fields struct {
	Account             *Account         `json:",omitempty"`
	AccountTxnID        *Hash256         `json:",omitempty"`
	Amendments          *Hash256         `json:",omitempty"`
	Balance             *Amount          `json:",omitempty"`
	BaseFee             *Uint64Hex       `json:",omitempty"`
	BookDirectory       *Hash256         `json:",omitempty"`
	BookNode            *NodeIndex       `json:",omitempty"`
	Domain              *VariableLength  `json:",omitempty"`
	EmailHash           *Hash128         `json:",omitempty"`
	ExchangeRate        *NodeIndex       `json:",omitempty"`
	Expiration          *uint32          `json:",omitempty"`
	FirstLedgerSequence *uint32          `json:",omitempty"`
	Flags               *LedgerEntryFlag `json:",omitempty"`
	Hashes              *Vector256       `json:",omitempty"`
	HighLimit           *Amount          `json:",omitempty"`
	HighNode            *NodeIndex       `json:",omitempty"`
	HighQualityIn       *uint32          `json:",omitempty"`
	HighQualityOut      *uint32          `json:",omitempty"`
	Indexes             *Vector256       `json:",omitempty"`
	IndexNext           *NodeIndex       `json:",omitempty"`
	IndexPrevious       *NodeIndex       `json:",omitempty"`
	LastLedgerSequence  *uint32          `json:",omitempty"`
	LowLimit            *Amount          `json:",omitempty"`
	LowNode             *NodeIndex       `json:",omitempty"`
	LowQualityIn        *uint32          `json:",omitempty"`
	LowQualityOut       *uint32          `json:",omitempty"`
	MessageKey          *VariableLength  `json:",omitempty"`
	Owner               *Account         `json:",omitempty"`
	OwnerCount          *uint32          `json:",omitempty"`
	OwnerNode           *NodeIndex       `json:",omitempty"`
	PreviousTxnID       *Hash256         `json:",omitempty"`
	PreviousTxnLgrSeq   *uint32          `json:",omitempty"`
	ReferenceFeeUnits   *uint32          `json:",omitempty"`
	RegularKey          *RegularKey      `json:",omitempty"`
	ReserveBase         *uint32          `json:",omitempty"`
	ReserveIncrement    *uint32          `json:",omitempty"`
	RootIndex           *Hash256         `json:",omitempty"`
	Sequence            *uint32          `json:",omitempty"`
	Signers             *VariableLength  `json:",omitempty"`
	TakerGets           *Amount          `json:",omitempty"`
	TakerGetsCurrency   *Hash160         `json:",omitempty"`
	TakerGetsIssuer     *Hash160         `json:",omitempty"`
	TakerPays           *Amount          `json:",omitempty"`
	TakerPaysCurrency   *Hash160         `json:",omitempty"`
	TakerPaysIssuer     *Hash160         `json:",omitempty"`
	TransferRate        *uint32          `json:",omitempty"`
	WalletLocator       *Hash256         `json:",omitempty"`
	WalletSize          *uint32          `json:",omitempty"`
}

type AffectedNode struct {
	FinalFields       *Fields         `json:",omitempty"`
	LedgerEntryType   LedgerEntryType `json:",omitempty"`
	LedgerIndex       *Hash256        `json:",omitempty"`
	PreviousFields    *Fields         `json:",omitempty"`
	NewFields         *Fields         `json:",omitempty"`
	PreviousTxnID     *Hash256        `json:",omitempty"`
	PreviousTxnLgrSeq *uint32         `json:",omitempty"`
}

type NodeEffect struct {
	ModifiedNode *AffectedNode `json:",omitempty"`
	CreatedNode  *AffectedNode `json:",omitempty"`
	DeletedNode  *AffectedNode `json:",omitempty"`
}

type NodeEffects []NodeEffect

type MetaData struct {
	AffectedNodes     NodeEffects
	TransactionIndex  uint32
	TransactionResult TransactionResult
	DeliveredAmount   *Amount `json:",omitempty"`
}

type TransactionSlice []*TransactionWithMetaData

func (s TransactionSlice) Len() int      { return len(s) }
func (s TransactionSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TransactionSlice) Less(i, j int) bool {
	if s[i].LedgerSequence == s[j].LedgerSequence {
		return s[i].MetaData.TransactionIndex < s[j].MetaData.TransactionIndex
	}
	return s[i].LedgerSequence < s[j].LedgerSequence
}

type TransactionWithMetaData struct {
	Transaction
	MetaData       MetaData `json:"meta"`
	LedgerSequence uint32   `json:"ledger_index"`
	Id             Hash256  `json:"-"`
}

func (t TransactionWithMetaData) GetType() string    { return "TransactionWithMetadata" }
func (t TransactionWithMetaData) Prefix() HashPrefix { return HP_TRANSACTION_NODE }
func (t TransactionWithMetaData) NodeType() NodeType { return NT_TRANSACTION_NODE }
func (t TransactionWithMetaData) Ledger() uint32     { return t.LedgerSequence }
func (t TransactionWithMetaData) NodeId() *Hash256   { return &t.Id }

func NewTransactionWithMetadata(typ TransactionType) *TransactionWithMetaData {
	return &TransactionWithMetaData{Transaction: TxFactory[typ]()}
}
