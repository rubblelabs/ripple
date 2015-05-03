package data

type TxBase struct {
	TransactionType    TransactionType
	Flags              *TransactionFlag `json:",omitempty"`
	SourceTag          *uint32          `json:",omitempty"`
	Account            Account
	Sequence           uint32
	Fee                Value
	AccountTxnID       *Hash256        `json:",omitempty"`
	SigningPubKey      *PublicKey      `json:",omitempty"`
	TxnSignature       *VariableLength `json:",omitempty"`
	Memos              Memos           `json:",omitempty"`
	PreviousTxnID      *Hash256        `json:",omitempty"`
	LastLedgerSequence *uint32         `json:",omitempty"`
	Hash               Hash256         `json:"hash"`
}

type Payment struct {
	TxBase
	Destination    Account
	Amount         Amount
	SendMax        *Amount  `json:",omitempty"`
	Paths          *PathSet `json:",omitempty"`
	DestinationTag *uint32  `json:",omitempty"`
	InvoiceID      *Hash256 `json:",omitempty"`
}

type AccountSet struct {
	TxBase
	EmailHash     *Hash128        `json:",omitempty"`
	WalletLocator *Hash256        `json:",omitempty"`
	WalletSize    *uint32         `json:",omitempty"`
	MessageKey    *VariableLength `json:",omitempty"`
	Domain        *VariableLength `json:",omitempty"`
	TransferRate  *uint32         `json:",omitempty"`
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
	LimitAmount Amount
	QualityIn   *uint32 `json:",omitempty"`
	QualityOut  *uint32 `json:",omitempty"`
}

type SetFee struct {
	TxBase
	BaseFee           Uint64Hex
	ReferenceFeeUnits uint32
	ReserveBase       uint32
	ReserveIncrement  uint32
}

type Amendment struct {
	TxBase
	Amendment Hash256
}

func (t *TxBase) GetBase() *TxBase                    { return t }
func (t *TxBase) GetType() string                     { return txNames[t.TransactionType] }
func (t *TxBase) GetTransactionType() TransactionType { return t.TransactionType }
func (t *TxBase) Prefix() HashPrefix                  { return HP_TRANSACTION_ID }
func (t *TxBase) GetPublicKey() *PublicKey            { return t.SigningPubKey }
func (t *TxBase) GetSignature() *VariableLength       { return t.TxnSignature }
func (t *TxBase) SigningPrefix() HashPrefix           { return HP_TRANSACTION_SIGN }
func (t *TxBase) PathSet() PathSet                    { return PathSet(nil) }
func (t *TxBase) GetHash() *Hash256                   { return &t.Hash }

func (t *TxBase) InitialiseForSigning() {
	if t.SigningPubKey == nil {
		t.SigningPubKey = new(PublicKey)
	}
	if t.TxnSignature == nil {
		t.TxnSignature = new(VariableLength)
	}
}

func (o *OfferCreate) Ratio() *Value {
	return o.TakerPays.Ratio(o.TakerGets)
}

func (p *Payment) PathSet() PathSet {
	if p.Paths == nil {
		return PathSet(nil)
	}
	return *p.Paths
}
