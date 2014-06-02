package data

type TransactionResult uint8

const (
	tesSUCCESS               TransactionResult = 0
	tecCLAIM                 TransactionResult = 100
	tecPATH_PARTIAL          TransactionResult = 101
	tecUNFUNDED_ADD          TransactionResult = 102
	tecUNFUNDED_OFFER        TransactionResult = 103
	tecUNFUNDED_PAYMENT      TransactionResult = 104
	tecFAILED_PROCESSING     TransactionResult = 105
	tecDIR_FULL              TransactionResult = 121
	tecINSUF_RESERVE_LINE    TransactionResult = 122
	tecINSUF_RESERVE_OFFER   TransactionResult = 123
	tecNO_DST                TransactionResult = 124
	tecNO_DST_INSUF_XRP      TransactionResult = 125
	tecNO_LINE_INSUF_RESERVE TransactionResult = 126
	tecNO_LINE_REDUNDANT     TransactionResult = 127
	tecPATH_DRY              TransactionResult = 128
	tecUNFUNDED              TransactionResult = 129
	tecMASTER_DISABLED       TransactionResult = 130
	tecNO_REGULAR_KEY        TransactionResult = 131
	tecOWNERS                TransactionResult = 132
)

type TxBase struct {
	hashable
	TransactionType    TransactionType
	Flags              *uint32 `json:",omitempty"`
	SourceTag          *uint32 `json:",omitempty"`
	Account            Account
	Sequence           uint32
	Fee                Value
	SigningPubKey      PublicKey
	TxnSignature       VariableLength `json:",omitempty"`
	Memos              Memos          `json:",omitempty"`
	PreviousTxnID      *Hash256       `json:",omitempty"`
	LastLedgerSequence *uint32        `json:",omitempty"`
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
	OfferSequence *uint32
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

func (t TransactionType) MarshalText() ([]byte, error) {
	return []byte(txNames[t]), nil
}

func (v *TransactionType) UnmarshalText(text []byte) (err error) {
	//FIXME: Currently a NoOp because TransactionType is already correctly
	//set if unmarshaling into a transaction created with GetTxFactoryByType

	return nil
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
