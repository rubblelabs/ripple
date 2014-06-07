package data

type TxBase struct {
	hashable
	TransactionType    TransactionType
	Flags              *uint32 `json:",omitempty"`
	SourceTag          *uint32 `json:",omitempty"`
	Account            Account
	Sequence           uint32
	Fee                Value
	SigningPubKey      *PublicKey      `json:",omitempty"`
	TxnSignature       *VariableLength `json:",omitempty"`
	Memos              Memos           `json:",omitempty"`
	PreviousTxnID      *Hash256        `json:",omitempty"`
	LastLedgerSequence *uint32         `json:",omitempty"`
}

type Payment struct {
	TxBase
	Destination    Account
	Amount         Amount
	SendMax        *Amount  `json:",omitempty"`
	Paths          *Paths   `json:",omitempty"`
	DestinationTag *uint32  `json:",omitempty"`
	InvoiceID      *Hash256 `json:",omitempty"`
}

type AccountSet struct {
	TxBase
	EmailHash     *Hash128        `json:",omitempty"`
	WalletLocator *Hash256        `json:",omitempty"`
	WalletSize    *uint32         `json:",omitempty"`
	MessageKey    *PublicKey      `json:",omitempty"`
	Domain        *VariableLength `json:",omitempty"`
	TransferRate  *uint32         `json:",string"`
	SetFlag       *uint32         `json:",omitempty"`
	ClearFlag     *uint32         `json:",omitempty"`
}

type SetRegularKey struct {
	TxBase
	RegularKey *RegularKey `json:",omitempty"`
}

type OfferCreate struct {
	TxBase
	OfferSequence *uint32 `json:",omitempty"`
	TakerPays     Amount
	TakerGets     Amount
	Expiration    *uint32 `json:",omitempty"`
}

type OfferCancel struct {
	TxBase
	OfferSequence uint32
}

type TrustSet struct {
	TxBase
	LimitAmount *Amount `json:",omitempty"`
	QualityIn   *uint32 `json:",omitempty"`
	QualityOut  *uint32 `json:",omitempty"`
}

type SetFee struct {
	TxBase
	BaseFee           uint64
	ReferenceFeeUnits uint32
	ReserveBase       uint32
	ReserveIncrement  uint32
}

type Amendment struct {
	TxBase
	Amendment Hash256
}

func (t TransactionType) String() string {
	return txNames[t]
}

func (t *TxBase) GetBase() *TxBase {
	return t
}

func (t *TxBase) GetTransactionType() TransactionType {
	return t.TransactionType
}

func (t *TxBase) GetType() string {
	return txNames[t.TransactionType]
}

func (t *TxBase) GetAccount() string {
	if a, err := t.Account.MarshalText(); err == nil {
		return string(a)
	}
	return ""
}

func (t *TransactionWithMetaData) GetAffectedNodes() []NodeEffect {
	return t.MetaData.AffectedNodes
}

func (t *TransactionWithMetaData) GetType() string {
	return "TransactionWithMetadata"
}
