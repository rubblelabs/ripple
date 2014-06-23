package data

type Validation struct {
	hashable
	Flags          uint32
	LedgerHash     Hash256
	LedgerSequence uint32
	Amendments     Vector256
	SigningTime    RippleTime
	SigningPubKey  PublicKey
	Signature      VariableLength
	// Fields below are part of serialization format
	// but never witnessed in the wild
	CloseTime        *uint32
	LoadFee          *uint32
	BaseFee          *uint64
	ReserveBase      *uint32
	ReserveIncrement *uint32
}

func (v *Validation) GetType() string {
	return "Validation"
}
