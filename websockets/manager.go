package websockets

import (
	"github.com/golang/glog"
)

var URIs []string = []string{
	"wss://s-west.ripple.com:443",
	"wss://s-east.ripple.com:443",
	"wss://s1.ripple.com:443",
}

type Manager struct {
	Outgoing chan interface{}
	Incoming chan interface{}
	remote   *Remote
}

func NewManager() (*Manager, error) {
	m := &Manager{
		Outgoing: make(chan interface{}, 10),
		Incoming: make(chan interface{}, 10),
	}
	return m, nil
}

func (m *Manager) Run() {
	uriIndex := 0

	for {
		err := m.handleConnection(URIs[uriIndex])
		glog.Errorln(err.Error())

		// Increment to next URI
		uriIndex = (uriIndex + 1) % len(URIs)
	}
}

func (m *Manager) handleConnection(uri string) (err error) {
	m.remote, err = NewRemote(uri)
	if err != nil {
		panic(err.Error())
	}

	return nil
}
