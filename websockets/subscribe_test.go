package websockets

import (
	"encoding/json"
	"github.com/donovanhide/ripple/data"
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

func (s *MessagesSuite) TestLedgerSubscribeResponse(c *C) {
	msg := SubscribeLedgers()
	readResponseFile(c, msg, "testdata/subscribe_ledger.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")
	c.Assert(msg.Id, Equals, uint64(3))

	// Result fields
	c.Assert(msg.Result.FeeBase, Equals, uint64(10))
	c.Assert(msg.Result.FeeRef, Equals, uint64(10))
	c.Assert(msg.Result.LedgerSequence, Equals, uint32(6959228))
	c.Assert(msg.Result.LedgerHash, Equals, "E23869F043A46C2735BCA40781A674C5F24460BAC26C6B7475550493A9180200")
	c.Assert(msg.Result.LedgerTime.String(), Equals, "2014-06-01 20:56:40")
	c.Assert(msg.Result.ReserveBase, Equals, uint64(20000000))
	c.Assert(msg.Result.ReserveIncrement, Equals, uint64(5000000))
	c.Assert(msg.Result.ValidatedLedgers, Equals, "32570-6959228")
	c.Assert(msg.Result.TxnCount, Equals, uint32(0))
}

func (s *MessagesSuite) TestLedgerStreamMsg(c *C) {
	msg := streamMessageFactory["ledgerClosed"]().(*LedgerStreamMsg)
	readResponseFile(c, msg, "testdata/ledger_stream.json")

	c.Assert(msg.FeeBase, Equals, uint64(10))
	c.Assert(msg.FeeRef, Equals, uint64(10))
	c.Assert(msg.LedgerSequence, Equals, uint32(6959229))
	c.Assert(msg.LedgerHash, Equals, "21EB30937A47EA6B71B63183806FFE9308CCB786137AA00FFB32A7094C6426FA")
	c.Assert(msg.LedgerTime.String(), Equals, "2014-06-01 20:56:40")
	c.Assert(msg.ReserveBase, Equals, uint64(20000000))
	c.Assert(msg.ReserveIncrement, Equals, uint64(5000000))
	c.Assert(msg.ValidatedLedgers, Equals, "32570-6959229")
	c.Assert(msg.TxnCount, Equals, uint32(1))
}

func (s *MessagesSuite) TestTransactionSubscribeResponse(c *C) {
	msg := SubscribeTransactions()
	readResponseFile(c, msg, "testdata/subscribe_transactions.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")
	c.Assert(msg.Id, Equals, uint64(3))
}

func (s *MessagesSuite) TestTransactionStreamMsg(c *C) {
	msg := streamMessageFactory["transaction"]().(*TransactionStreamMsg)
	readResponseFile(c, msg, "testdata/transactions_stream.json")

	c.Assert(msg.EngineResult, Equals, "tesSUCCESS")
	c.Assert(msg.EngineResultCode, Equals, 0)
	c.Assert(msg.EngineResultMessage, Equals, "The transaction was applied.")
	c.Assert(msg.LedgerHash, Equals, "9B0E9D19E8246BA9B224078B73158ED8970B90DBFAAA68D73A2E0E2899B5AF5A")
	c.Assert(msg.LedgerSequence, Equals, uint32(6959249))
	c.Assert(msg.Status, Equals, "closed")
	c.Assert(msg.Validated, Equals, true)

	c.Assert(msg.Transaction.(*data.OfferCreate).GetType(), Equals, "OfferCreate")
	c.Assert(msg.Transaction.(*data.OfferCreate).GetAccount(), Equals, "rPEZyTnSyQyXBCwMVYyaafSVPL8oMtfG6a")
	c.Assert(msg.Transaction.(*data.OfferCreate).Fee.String(), Equals, "0.00005")
	//FIXME(luke): Hash is not unmarshaled from json
	//c.Assert(msg.Transaction.(*data.OfferCreate).Hash().String(), Equals, "25174B56C40B090D4AFCDAC3F07DCCF8A49A096D62CE1CE6864A8624F790F980")
	c.Assert(msg.Transaction.(*data.OfferCreate).SigningPubKey.String(), Equals, "0309AEAA170F651170F85C85237CD25CD4200CF91C1C05A9B8A19E72912C2254DF")
	c.Assert(msg.Transaction.(*data.OfferCreate).TxnSignature.String(), Equals, "304402201480DBC8253B2E5CCB24001C6E6A0AE73C8FC8D6237B0AA1A5B1CADA92306070022013B02C3CE6E7AFD5F8F348BC40975D15056D414BBC11AD2EA04A65496482212E")
	c.Assert(msg.Transaction.(*data.OfferCreate).Sequence, Equals, uint32(753273))

	c.Assert(*msg.Transaction.(*data.OfferCreate).OfferSequence, Equals, uint32(753240))
	c.Assert(msg.Transaction.(*data.OfferCreate).TakerGets.String(), Equals, "6400.064/XRP")
	c.Assert(msg.Transaction.(*data.OfferCreate).TakerPays.String(), Equals, "174.72/CNY/razqQKzJRdB4UxFPWf5NEpEG3WMkmwgcXA")
}

func (s *MessagesSuite) TestServerSubscribeResponse(c *C) {
	msg := SubscribeServer()
	readResponseFile(c, msg, "testdata/subscribe_server.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")
	c.Assert(msg.Id, Equals, uint64(3))

	// Result fields
	c.Assert(msg.Result.Status, Equals, "full")
	c.Assert(msg.Result.LoadBase, Equals, 256)
	c.Assert(msg.Result.LoadFactor, Equals, 256)
}

func (s *MessagesSuite) TestServerStreamMsg(c *C) {
	msg := streamMessageFactory["serverStatus"]().(*ServerStreamMsg)
	readResponseFile(c, msg, "testdata/server_stream.json")

	c.Assert(msg.Status, Equals, "syncing")
	c.Assert(msg.LoadBase, Equals, 256)
	c.Assert(msg.LoadFactor, Equals, 256)
}
