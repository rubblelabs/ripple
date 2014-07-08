package websockets

import (
	"encoding/json"
	"fmt"
	"github.com/rubblelabs/ripple/data"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"net"
	"net/url"
	"reflect"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Remote struct {
	Outgoing chan Syncer
	Incoming chan interface{}
	ws       *websocket.Conn
}

func NewRemote(endpoint string) (*Remote, error) {
	glog.Infoln(endpoint)
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	c, err := net.DialTimeout("tcp", u.Host, time.Second*5)
	if err != nil {
		return nil, err
	}
	ws, _, err := websocket.NewClient(c, u, nil, 1024, 1024)
	if err != nil {
		return nil, err
	}
	return &Remote{
		Outgoing: make(chan Syncer, 10),
		Incoming: make(chan interface{}, 10),
		ws:       ws,
	}, nil
}

func (r *Remote) Close() {
	close(r.Outgoing)
}

func (r *Remote) Run() {
	outbound := make(chan interface{})
	inbound := make(chan []byte)
	pending := make(map[uint64]Syncer)

	defer func() {
		close(outbound)
		close(r.Incoming)
		for _, c := range pending {
			c.Fail("Connection Closed")
		}
	}()

	// Spawn read/write goroutines
	go func() {
		defer r.ws.Close()
		r.writePump(outbound)
	}()
	go r.readPump(inbound)

	// Main run loop
	var response Command
	for {
		select {
		case command, ok := <-r.Outgoing:
			if !ok {
				glog.Infoln("Closing")
				return
			}
			outbound <- command
			id := reflect.ValueOf(command).Elem().FieldByName("Id").Uint()
			pending[id] = command

		case in, ok := <-inbound:
			if !ok {
				return
			}

			if err := json.Unmarshal(in, &response); err != nil {
				glog.Errorln(err.Error())
				continue
			}
			// Stream message
			factory, ok := streamMessageFactory[response.Type]
			if ok {
				cmd := factory()
				if err := json.Unmarshal(in, &cmd); err != nil {
					glog.Errorln(err.Error(), string(in))
					continue
				}
				r.Incoming <- cmd
				continue
			}

			// Command response message
			cmd, ok := pending[response.Id]
			if !ok {
				glog.Errorf("Unexpected message: %+v", response)
				continue
			}
			delete(pending, response.Id)
			if err := json.Unmarshal(in, &cmd); err != nil {
				glog.Errorln(err.Error())
				continue
			}
			cmd.Done()
		}
	}
}

// Synchronously get a single transaction
func (r *Remote) Tx(hash data.Hash256) (*TxResult, error) {
	cmd := &TxCommand{
		Command:     newCommand("tx"),
		Transaction: hash,
	}
	r.Outgoing <- cmd
	<-cmd.Ready
	if cmd.CommandError != nil {
		return nil, cmd.CommandError
	}
	return cmd.Result, nil
}

func (r *Remote) accountTx(account data.Account, c chan *data.TransactionWithMetaData, pageSize int) {
	defer close(c)
	cmd := newAccountTxCommand(account, pageSize, nil)
	for ; ; cmd = newAccountTxCommand(account, pageSize, cmd.Result.Marker) {
		r.Outgoing <- cmd
		<-cmd.Ready
		if cmd.CommandError != nil {
			glog.Errorln(cmd.Error())
			return
		}
		for _, tx := range cmd.Result.Transactions {
			c <- tx
		}
		if cmd.Result.Marker == nil {
			return
		}
	}
}

// Asynchronously retrieve all transactions for an account
func (r *Remote) AccountTx(account data.Account, pageSize int) chan *data.TransactionWithMetaData {
	c := make(chan *data.TransactionWithMetaData)
	go r.accountTx(account, c, pageSize)
	return c
}

// Synchronously submit a single transaction
func (r *Remote) Submit(tx data.Transaction) (*SubmitResult, error) {
	_, raw, err := data.Raw(tx)
	if err != nil {
		return nil, err
	}
	cmd := &SubmitCommand{
		Command: newCommand("submit"),
		TxBlob:  fmt.Sprintf("%X", raw),
	}
	r.Outgoing <- cmd
	<-cmd.Ready
	if cmd.CommandError != nil {
		return nil, cmd.CommandError
	}
	return cmd.Result, nil
}

// Synchronously gets ledger entries
func (r *Remote) LedgerData(ledger interface{}, marker *data.Hash256) (*LedgerDataResult, error) {
	cmd := &LedgerDataCommand{
		Command: newCommand("ledger_data"),
		Ledger:  ledger,
		Marker:  marker,
	}
	r.Outgoing <- cmd
	<-cmd.Ready
	if cmd.CommandError != nil {
		return nil, cmd.CommandError
	}
	return cmd.Result, nil
}

// Synchronously gets a single ledger
func (r *Remote) Ledger(ledger interface{}, transactions bool) (*LedgerResult, error) {
	cmd := &LedgerCommand{
		Command:      newCommand("ledger"),
		LedgerIndex:  ledger,
		Transactions: transactions,
		Expand:       true,
	}
	r.Outgoing <- cmd
	<-cmd.Ready
	if cmd.CommandError != nil {
		return nil, cmd.CommandError
	}
	return cmd.Result, nil
}

// Synchronously requests paths
func (r *Remote) RipplePathFind(src, dest data.Account, amount data.Amount, srcCurr *[]data.Currency) (*RipplePathFindResult, error) {
	cmd := &RipplePathFindCommand{
		Command:       newCommand("ledger"),
		SrcAccount:    src,
		SrcCurrencies: srcCurr,
		DestAccount:   dest,
		DestAmount:    amount,
	}
	r.Outgoing <- cmd
	<-cmd.Ready
	if cmd.CommandError != nil {
		return nil, cmd.CommandError
	}
	return cmd.Result, nil
}

// Synchronously subscribe to streams and receive a confirmation message
// Streams are recived asynchronously over the Incoming channel
func (r *Remote) Subscribe(ledger, transactions, server bool) (*SubscribeResult, error) {
	streams := []string{}
	if ledger {
		streams = append(streams, "ledger")
	}
	if transactions {
		streams = append(streams, "transactions")
	}
	if server {
		streams = append(streams, "server")
	}
	cmd := &SubscribeCommand{
		Command: newCommand("subscribe"),
		Streams: streams,
	}
	r.Outgoing <- cmd
	<-cmd.Ready
	if cmd.CommandError != nil {
		return nil, cmd.CommandError
	}
	return cmd.Result, nil

}

// Reads from the websocket and sends to inbound channel
// Expects to receive PONGs at specified interval, or kills the session
func (r *Remote) readPump(inbound chan []byte) {
	r.ws.SetReadDeadline(time.Now().Add(pongWait))
	r.ws.SetPongHandler(func(string) error { r.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := r.ws.ReadMessage()
		if err != nil {
			glog.Errorln(err)
			close(inbound)
			return
		}
		glog.V(2).Infoln(dump(message))
		r.ws.SetReadDeadline(time.Now().Add(pongWait))
		inbound <- message
	}
}

// Consumes from the outbound channel and sends them over the websocket.
// Also sends PING messages at specified interval
func (r *Remote) writePump(outbound chan interface{}) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {

		// An outbound message is available to send
		case message, ok := <-outbound:
			if !ok {
				r.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			b, err := json.Marshal(message)
			if err != nil {
				// Outbound message cannot be JSON serialized (log it and continue)
				glog.Errorln(err)
				continue
			}

			glog.V(2).Infoln(dump(b))
			if err := r.ws.WriteMessage(websocket.TextMessage, b); err != nil {
				glog.Errorln(err)
				return
			}

		// Time to send a ping
		case <-ticker.C:
			if err := r.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				glog.Errorln(err)
				return
			}
		}
	}
}

func dump(b []byte) string {
	var v map[string]interface{}
	json.Unmarshal(b, &v)
	out, _ := json.MarshalIndent(v, "", "  ")
	return string(out)
}
