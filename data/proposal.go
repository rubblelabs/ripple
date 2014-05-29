package data

type Proposal struct {
	hashable
	LedgerHash     Hash256
	PreviousLedger Hash256
	Sequence       uint32
	CloseTime      uint32
	PublicKey      PublicKey
	Signature      VariableLength
}

func (p *Proposal) GetType() string {
	return "Proposal"
}
