package tx

import (
	"github.com/maybeTomorrow/ripple/data"
	"github.com/maybeTomorrow/ripple/websockets"
	"math"
)

const LEDGER_OFFSET = 20
const RIPPLE_EPOCH_DIFF int64 = 0x386d4380

func RippleTimeToUnixTime(t uint32) int64 {
	return int64(t) + RIPPLE_EPOCH_DIFF
}
func UnixTimeToRippleTime(t int64) *uint32 {
	rs := uint32(t - RIPPLE_EPOCH_DIFF)
	return &rs
}

func CalFee(tx data.Transaction, remote *websockets.Remote, cushion float64) data.Value {

	rs, _ := remote.ServerState()
	baseFee := rs.State.ValidatedLedger.BaseFee

	baseFee = baseFee * rs.State.LoadFactor / rs.State.LoadBase
	if cushion == 0 {
		cushion = 1.2
	}
	fee := uint32(math.Round(baseFee * cushion))

	lls := uint32(rs.State.ValidatedLedger.Seq + LEDGER_OFFSET)
	tx.GetBase().LastLedgerSequence = &lls

	if tx.GetTransactionType() == data.ESCROW_FINISH {

	}
	if tx.GetTransactionType() == data.ACCOUNT_DELETE {
		fee = fetchAccountDeleteFee(remote)
	}
	v, _ := data.NewNativeValue(int64(fee))
	return *v

}

func fetchAccountDeleteFee(remote *websockets.Remote) uint32 {
	rs, _ := remote.ServerState()
	if rs.Status == "success" {
		return uint32(math.Round(rs.State.ValidatedLedger.ReserveInc * 1e6))
	}
	return 0
}
