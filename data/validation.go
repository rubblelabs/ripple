package data

type Validation struct {
	hashable
	Flags            uint32
	LedgerHash       Hash256
	LedgerSequence   uint32
	CloseTime        uint32
	LoadFee          uint32
	Amendments       Vector256
	BaseFee          uint64
	ReserveBase      uint32
	ReserveIncrement uint32
	SigningTime      uint32
	SigningPubKey    PublicKey
	Signature        VariableLength
}

func (*Validation) GetType() string {
	return "Validation"
}
