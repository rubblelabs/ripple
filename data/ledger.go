package data

type LedgerHeader struct {
	LedgerSequence  uint32  `json:"ledger_index"`
	TotalXRP        uint64  `json:"total_coins"`
	PreviousLedger  Hash256 `json:"parent_hash"`
	TransactionHash Hash256 `json:"transaction_hash"`
	StateHash       Hash256 `json:"account_hash"`
	ParentCloseTime uint32
	CloseTime       uint32 `json:"close_time"`
	CloseResolution uint8  `json:"close_time_resolution"`
	CloseFlags      uint8
}

type Ledger struct {
	hashable
	LedgerHeader
	Fees uint64
}

func NewEmptyLedger(sequence uint32) *Ledger {
	return &Ledger{
		LedgerHeader: LedgerHeader{
			LedgerSequence: sequence,
		},
	}
}

func (l *Ledger) GetType() string {
	return "LedgerMaster"
}
