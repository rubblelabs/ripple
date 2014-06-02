package data

type Validation struct {
	hashable
	Flags            uint32
	LedgerHash       Hash256
	LedgerSequence   uint32
	CloseTime        *uint32 // Should exist?
	LoadFee          *uint32 // Should exist?
	Amendments       Vector256
	BaseFee          *uint64 // Should exist?
	ReserveBase      *uint32 // Should exist?
	ReserveIncrement *uint32 // Should exist?
	SigningTime      uint32
	SigningPubKey    PublicKey
	Signature        VariableLength
}

func (v *Validation) GetType() string {
	return "Validation"
}
