package websockets

import (
	"encoding/json"
	"github.com/rubblelabs/ripple/data"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MessagesSuite struct{}

var _ = Suite(&MessagesSuite{})

func readResponseFile(c *C, msg interface{}, path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		c.Error(err)
	}

	if err = json.Unmarshal(b, msg); err != nil {
		c.Error(err)
	}
}

func (s *MessagesSuite) TestLedgerResponse(c *C) {
	msg := &LedgerCommand{}
	readResponseFile(c, msg, "testdata/ledger.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")

	// Result fields
	c.Assert(msg.Result.Ledger.LedgerSequence, Equals, uint32(6917762))
	c.Assert(msg.Result.Ledger.Accepted, Equals, true)
	c.Assert(msg.Result.Ledger.CloseTime.String(), Equals, "2014-May-30 13:11:50")
	c.Assert(msg.Result.Ledger.Closed, Equals, true)
	c.Assert(msg.Result.Ledger.Hash.String(), Equals, "0C5C5B39EA40D40ACA6EB47E50B2B85FD516D1A2BA67BA3E050349D3EF3632A4")
	c.Assert(msg.Result.Ledger.PreviousLedger.String(), Equals, "F8F0363803C30E659AA24D6A62A6512BA24BEA5AC52A29731ABA1E2D80796E8B")
	c.Assert(msg.Result.Ledger.TotalXRP, Equals, uint64(99999990098968782))
	c.Assert(msg.Result.Ledger.AccountHash.String(), Equals, "46D3E36FE845B9A18293F4C0F134D7DAFB06D4D9A1C7E4CB03F8B293CCA45FA0")
	c.Assert(msg.Result.Ledger.TransactionHash.String(), Equals, "757CCB586D44F3C58E366EC7618988C0596277D3D5D0B412E49563B5EEDF04FF")

	c.Assert(msg.Result.Ledger.Transactions, HasLen, 7)
	tx0 := msg.Result.Ledger.Transactions[0]
	c.Assert(tx0.GetHash().String(), Equals, "2D0CE11154B655A2BFE7F3F857AAC344622EC7DAB11B1EBD920DCDB00E8646FF")
	c.Assert(tx0.MetaData.AffectedNodes, HasLen, 4)
}

func (s *MessagesSuite) TestTxResponse(c *C) {
	msg := &TxCommand{}
	readResponseFile(c, msg, "testdata/tx.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")

	// Result fields
	c.Assert(msg.Result.Validated, Equals, true)
	c.Assert(msg.Result.MetaData.AffectedNodes, HasLen, 4)
	c.Assert(msg.Result.MetaData.TransactionResult.String(), Equals, "tesSUCCESS")

	offer := msg.Result.Transaction.(*data.OfferCreate)
	c.Assert(msg.Result.GetHash().String(), Equals, "2D0CE11154B655A2BFE7F3F857AAC344622EC7DAB11B1EBD920DCDB00E8646FF")
	c.Assert(offer.GetType(), Equals, "OfferCreate")
	c.Assert(offer.Account.String(), Equals, "rwpxNWdpKu2QVgrh5LQXEygYLshhgnRL1Y")
	c.Assert(offer.Fee.String(), Equals, "0.00001")
	c.Assert(offer.SigningPubKey.String(), Equals, "02BD6F0CFD0182F2F408512286A0D935C58FF41169DAC7E721D159D711695DFF85")
	c.Assert(offer.TxnSignature.String(), Equals, "30440220216D42DF672C1CC7EF0CA9C7840838A2AF5FEDD4DEFCBA770C763D7509703C8702203C8D831BFF8A8BC2CC993BECB4E6C7BE1EA9D394AB7CE7C6F7542B6CDA781467")
	c.Assert(offer.Sequence, Equals, uint32(1681497))
}

func (s *MessagesSuite) TestAccountTxResponse(c *C) {
	msg := &AccountTxCommand{}
	readResponseFile(c, msg, "testdata/account_tx.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")

	c.Assert(len(msg.Result.Transactions), Equals, 2)
	offer := msg.Result.Transactions[1].Transaction.(*data.OfferCreate)
	c.Assert(offer.TakerPays.String(), Equals, "0.034800328/BTC/rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B")
}

func (s *MessagesSuite) TestLedgerDataResponse(c *C) {
	msg := &LedgerDataCommand{}
	readResponseFile(c, msg, "testdata/ledger_data.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")

	c.Assert(msg.Result.LedgerSequence, Equals, uint32(6281820))
	c.Assert(msg.Result.Hash.String(), Equals, "83CC350B1CDD9792D47F60D3DBB7673518FD6E71821070673E6EAE65DE69086B")
	c.Assert(msg.Result.Marker.String(), Equals, "02DE1A2AD4332A1AF01C59F16E45218FA70E5792BD963B6D7ACF188D6D150607")
	c.Assert(len(msg.Result.State), Equals, 2048)
	c.Assert(msg.Result.State[0].GetType(), Equals, "AccountRoot")
}
