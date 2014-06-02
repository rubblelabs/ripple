package websockets

import (
	"encoding/json"
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
	Outgoing chan interface{}
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
		Outgoing: make(chan interface{}, 10),
		Incoming: make(chan interface{}, 10),
		ws:       ws,
	}, nil
}

func (r *Remote) Run() {
	outbound := make(chan interface{})
	inbound := make(chan []byte)
	pending := make(map[uint64]interface{})
	go func() {
		defer r.ws.Close()
		go r.readPump(inbound)
		r.writePump(outbound)
	}()
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
					glog.Errorln(err.Error())
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
		}
	}
}

func (r *Remote) readPump(inbound chan []byte) {
	r.ws.SetReadDeadline(time.Now().Add(pongWait))
	r.ws.SetPongHandler(func(string) error { r.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := r.ws.ReadMessage()
		if err != nil {
			glog.Errorln(err)
			break
		}
		glog.V(2).Infoln(string(message))
		inbound <- message
	}
}

func (r *Remote) writePump(outbound chan interface{}) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-outbound:
			if !ok {
				r.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			b, err := json.Marshal(message)
			if err != nil {
				glog.Errorln(err)
				return
			}
			glog.V(2).Infoln(string(b))
			if err := r.ws.WriteMessage(websocket.TextMessage, b); err != nil {
				glog.Errorln(err)
				return
			}
		case <-ticker.C:
			if err := r.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				glog.Errorln(err)
				return
			}
		}
	}
}
