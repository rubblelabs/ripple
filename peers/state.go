package peers

import (
	"encoding/json"
	"fmt"
	"github.com/rubblelabs/ripple/crypto"
	"github.com/rubblelabs/ripple/peers/protocol"
	"sync"
	"time"
)

type VerificationStatus int

const (
	Unverified VerificationStatus = iota
	Verified
	HelloFailed
	Disconnected
)

var verificationStatusMap = map[VerificationStatus]string{
	Unverified:   "Unverified",
	Verified:     "Verified",
	HelloFailed:  "HelloFailed",
	Disconnected: "Disconnected",
}

type State struct {
	Trusted       bool
	Name          string
	MajorVersion  uint32
	MinorVersion  uint32
	PublicKey     crypto.Hash
	CurrentLedger uint32
	MinLedger     uint32
	MaxLedger     uint32
	Event         string
	NodeStatus    string
	Status        VerificationStatus
	Discovered    time.Time
}

type PeerState struct {
	*State
	mu sync.RWMutex
}

func NewPeerState(c *PeerConnection) *PeerState {
	return &PeerState{
		State: &State{
			Trusted:    c.Trusted,
			Status:     Unverified,
			Discovered: time.Now(),
		},
	}
}

func (s *PeerState) UpdateStatus(status VerificationStatus) {
	s.mu.Lock()
	s.Status = status
	s.mu.Unlock()
}

func (s *PeerState) UpdateState(state *protocol.TMStatusChange) {
	s.mu.Lock()
	s.CurrentLedger = state.GetLedgerSeq()
	s.MinLedger = state.GetFirstSeq()
	s.MaxLedger = state.GetLastSeq()
	s.Event = state.GetNewEvent().String()
	s.NodeStatus = state.GetNewStatus().String()
	s.mu.Unlock()
}

func (s *PeerState) ProcessHello(hello *protocol.Hello) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CurrentLedger = hello.GetLedgerIndex()
	s.Name = hello.GetFullVersion()
	s.MajorVersion = hello.GetProtoVersion()
	s.MinorVersion = hello.GetProtoVersionMin()
	var err error
	s.PublicKey, err = crypto.NewRippleHash(string(hello.NodePublic))
	if err != nil {
		s.Status = HelloFailed
		return fmt.Errorf("Bad node public key: %s", hello.NodePublic)
	}
	s.Status = Verified
	return nil
}

func (s *PeerState) GetLedgerRange() (uint32, uint32) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.MinLedger, s.MaxLedger
}

func (s *PeerState) MarshalJSON() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s.State)
}

func (s VerificationStatus) MarshalText() ([]byte, error) {
	return []byte(verificationStatusMap[s]), nil
}
