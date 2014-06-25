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

func (v *Validation) GetType() string               { return "Validation" }
func (v *Validation) GetPublicKey() *PublicKey      { return &v.SigningPubKey }
func (v *Validation) GetSignature() *VariableLength { return &v.Signature }

func (v Validation) SigningHash() (Hash256, error) {
	if err := NewEncoder().Validation(&v, true); err != nil {
		return zero256, err
	}
	return hashValues([]interface{}{
		HP_VALIDATION,
		v.Raw(),
	})
}

func (v Validation) SuppressionId() (Hash256, error) {
	return hashValues([]interface{}{
		v.Raw(),
	})
}
