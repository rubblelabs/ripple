package peers

import (
	"encoding/json"
	"github.com/donovanhide/ripple/crypto"
	"github.com/donovanhide/ripple/ledger"
	"github.com/golang/glog"
	"net"
)

type Manager struct {
	Name      string
	PublicKey crypto.Hash
	Port      string
	Quit      chan bool
	Status    chan chan []byte
	peers     chan *PeerConnection
	connected chan *Peer
	key       *crypto.RootDeterministicKey
}

func NewManager(key *crypto.RootDeterministicKey, name string, port string) (*Manager, error) {
	public, err := key.PublicNodeKey()
	if err != nil {
		return nil, err
	}
	return &Manager{
		Name:      name,
		Port:      port,
		PublicKey: public,
		Status:    make(chan chan []byte),
		Quit:      make(chan bool),
		peers:     make(chan *PeerConnection, 10),
		connected: make(chan *Peer, 10),
		key:       key,
	}, nil
}

func (m *Manager) Start(l *ledger.Manager) {
	seen := make(map[string]struct{})
	peers := make([]*Peer, 0)
	go Listen(m, m.Port)
	for {
		select {
		case c := <-m.Status:
			var dump []*Dump
			for _, peer := range peers {
				dump = append(dump, peer.GetDump())
			}
			out, err := json.MarshalIndent(dump, "", "\t")
			if err != nil {
				glog.Infoln(err)
				c <- []byte(nil)
			} else {
				c <- out
			}
		case c := <-m.peers:
			if _, ok := seen[c.String()]; !ok {
				seen[c.String()] = struct{}{}
				go m.connectPeer(c, l)
			}
		case peer := <-m.connected:
			peers = append(peers, peer)
		case <-m.Quit:
			return
		}
	}
}

func (m *Manager) connectPeer(c *PeerConnection, l *ledger.Manager) {
	glog.Infof("Peer Manager: New Peer: %s ", c.String())
	peer, err := NewPeer(c)
	if err == nil {
		go peer.handle(m, l)
		glog.Infof("Peer Manager: New Peer: %s successful connection", c.String())
	} else {
		glog.Infof("Peer Manager: New Peer Error: %s", err.Error())
		peer.UpdateStatus(Disconnected)
	}
	m.connected <- peer
}

func (m *Manager) AddPeer(host, port string, trusted bool, conn net.Conn) {
	m.peers <- &PeerConnection{
		Host:    host,
		Port:    port,
		Trusted: trusted,
		Conn:    conn,
	}
}
