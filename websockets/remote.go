package websockets

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"launchpad.net/tomb"
	"net"
	"net/url"
	"reflect"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 5 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 15 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Remote struct {
	Outgoing chan interface{}
	Incoming chan interface{}
	ws       *websocket.Conn
	t        tomb.Tomb
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
		Outgoing: make(chan interface{}, 10),
		Incoming: make(chan interface{}, 10),
		ws:       ws,
	}, nil
}

func (r *Remote) Run() {
	outbound := make(chan interface{})
	inbound := make(chan []byte)
	pending := make(map[uint64]interface{})

	defer func() {
		close(outbound)
		close(inbound)
		close(r.Outgoing)
		close(r.Incoming)
		r.t.Done()
	}()

	// Spawn read/write goroutines
	go func() {
		defer r.ws.Close()
		r.writePump(outbound)
	}()
	go r.readPump(inbound)

	// Main run loop
	var response Response
	for {
		select {
		case command := <-r.Outgoing:
			outbound <- command
			id := reflect.ValueOf(command).Elem().FieldByName("Id").Uint()
			pending[id] = command

		case in := <-inbound:
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
			r.Incoming <- cmd

		case <-r.t.Dying():
			return
		}
	}
}

// Waits for the session to close and returns the error (if any)
func (r *Remote) Wait() error {
	return r.t.Wait()
}

// Reads from the websocket and sends to inbound channel
// Expects to receive PONGs at specified interval, or kills the session
func (r *Remote) readPump(inbound chan []byte) {
	r.ws.SetReadDeadline(time.Now().Add(pongWait))
	r.ws.SetPongHandler(func(string) error { r.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := r.ws.ReadMessage()
		if err != nil {
			r.t.Kill(err)
			return
		}
		glog.V(2).Infoln(string(message))
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
				r.t.Kill(fmt.Errorf("outbound channel closed"))
				return
			}

			b, err := json.Marshal(message)
			if err != nil {
				// Outbound message cannot be JSON serialized (log it and continue)
				glog.Errorln(err)
				continue
			}

			glog.V(2).Infoln(string(b))
			if err := r.ws.WriteMessage(websocket.TextMessage, b); err != nil {
				r.t.Kill(err)
				return
			}

		// Time to send a ping
		case <-ticker.C:
			if err := r.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				r.t.Kill(err)
				return
			}

		// The session is shutting down
		case <-r.t.Dying():
			r.ws.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
	}
}
