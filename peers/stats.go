package peers

import (
	"code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"fmt"
	"github.com/donovanhide/ripple/peers/protocol"
	metrics "github.com/rcrowley/go-metrics"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MessageStatus int

const (
	PassThrough MessageStatus = iota
	Expected
	Unexpected
)

type Latency struct {
	Count uint64
	Total time.Duration
	Last  time.Time
	Min   time.Duration
	Max   time.Duration
}

type LatencyDump struct {
	Average, Min, Max float64
}

type Stats struct {
	Sent       map[string]uint64
	Received   map[string]uint64
	Latencies  map[string]*Latency
	Unexpected uint64
	InFlight   string
}

type PeerStats struct {
	*Stats
	mu sync.RWMutex
}

func NewPeerStats() *PeerStats {
	return &PeerStats{
		Stats: &Stats{
			Sent:      make(map[string]uint64),
			Received:  make(map[string]uint64),
			Latencies: make(map[string]*Latency),
		},
	}
}

func shortName(msg proto.Message) string {
	return strings.TrimPrefix(reflect.TypeOf(msg).String(), "*protocol.")
}

func latencyName(msg proto.Message, inbound bool) (string, string) {
	switch v := msg.(type) {
	case *protocol.TMLedgerData:
		if inbound {
			return "LedgerData", strconv.FormatUint(uint64(v.GetLedgerSeq()), 10)
		}
	case *protocol.TMGetLedger:
		if !inbound {
			return "LedgerData", strconv.FormatUint(uint64(v.GetLedgerSeq()), 10)
		}
	case *protocol.TMGetObjectByHash:
		if (inbound && !v.GetQuery()) || (!inbound && v.GetQuery()) {
			return "GetObjectByHash", fmt.Sprintf("%d:%d", v.GetSeq(), len(v.GetObjects()))
		}
	}
	return "", ""
}

func (s *PeerStats) Send(msg proto.Message) {
	name := shortName(msg)
	latencyName, latencyId := latencyName(msg, false)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Sent[name]++
	if latencyName != "" {
		s.InFlight = latencyName + ":" + latencyId
		if l, ok := s.Latencies[latencyName]; !ok {
			s.Latencies[latencyName] = &Latency{
				Last: time.Now(),
				Min:  time.Hour * 24,
			}
		} else {
			l.Last = time.Now()
		}
	}
}

func (s *PeerStats) Receive(msg proto.Message) MessageStatus {
	status := PassThrough
	name := shortName(msg)
	latencyName, latencyId := latencyName(msg, true)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Received[name]++
	if latencyName != "" {
		status = Expected
		if s.InFlight != latencyName+":"+latencyId {
			s.Unexpected++
			status = Unexpected
		}
		s.InFlight = ""
		// metrics.GetOrRegisterCounter(latencyName, nil).Inc(1)
		if l, ok := s.Latencies[latencyName]; ok {
			metrics.GetOrRegisterTimer(latencyName, nil).UpdateSince(l.Last)
			diff := time.Now().Sub(l.Last)
			l.Total += diff
			if diff < l.Min {
				l.Min = diff
			}
			if diff > l.Max {
				l.Max = diff
			}
			l.Count++

		} else {
			panic(name)
		}
	}
	return status
}

func (s *PeerStats) MarshalJSON() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s.Stats)
}

func (l *Latency) MarshalJSON() ([]byte, error) {
	average := l.Total.Seconds() / float64(l.Count)
	if l.Count == 0 {
		average = 0
	}
	return json.Marshal(LatencyDump{
		Average: average,
		Min:     l.Min.Seconds(),
		Max:     l.Max.Seconds(),
	})
}
