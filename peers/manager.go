package peers

import (
	"encoding/json"
	"fmt"
	"github.com/donovanhide/ripple/crypto"
	"github.com/donovanhide/ripple/ledger"
	"github.com/golang/glog"
	"net"
	"strings"
)

type Config struct {
	Key      crypto.Key
	Name     string
	Port     string
	Sync     ledger.Sync
	MaxPeers int
	Trusted  string
}

type Manager struct {
	*Config
	PublicKey crypto.Hash
	Quit      chan bool
	Status    chan chan []byte
	peers     chan *PeerConnection
	connected chan *Peer
}

func NewManager(config *Config) (*Manager, error) {
	mgr := &Manager{
		Config:    config,
		Status:    make(chan chan []byte),
		Quit:      make(chan bool),
		peers:     make(chan *PeerConnection, 10),
		connected: make(chan *Peer, 10),
	}
	var err error
	mgr.PublicKey, err = crypto.NewRipplePublicNode(mgr.Key.PublicCompressed())
	if err != nil {
		return nil, err
	}
	go mgr.run()
	for _, address := range strings.Split(mgr.Trusted, ",") {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return nil, fmt.Errorf("Bad trusted peer: %s Part: %s", config.Trusted, address)
		}
		mgr.AddPeer(host, port, true, nil)
	}
	return mgr, nil
}

func (m *Manager) run() {
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
				if len(peers) < m.MaxPeers {
					go m.connectPeer(c)
				}
			}
		case peer := <-m.connected:
			peers = append(peers, peer)
		case <-m.Quit:
			return
		}
	}
}

func (m *Manager) connectPeer(c *PeerConnection) {
	glog.Infof("Peer Manager: New Peer: %s ", c.String())
	peer, err := NewPeer(c, m.Sync)
	if err == nil {
		go peer.handle(m)
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
