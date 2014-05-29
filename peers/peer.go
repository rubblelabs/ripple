package peers

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"github.com/donovanhide/ripple/crypto"
	"github.com/donovanhide/ripple/data"
	"github.com/donovanhide/ripple/ledger"
	"github.com/donovanhide/ripple/peers/protocol"
	"github.com/golang/glog"
	"strconv"
	"sync"
	"time"
)

const (
	MAJOR_VERSION = 65538
	MINOR_VERSION = 65538
)

type Dump struct {
	Host     string
	Port     string
	Resolved string
	State    *PeerState
	Stats    *PeerStats
}

type Peer struct {
	*Conn
	*PeerState
	*PeerStats
	Outgoing    chan proto.Message
	synchronous chan proto.Message
}

func NewPeer(c *PeerConnection) (*Peer, error) {
	peer := &Peer{
		PeerState:   NewPeerState(c),
		PeerStats:   NewPeerStats(),
		Outgoing:    make(chan proto.Message, 10),
		synchronous: make(chan proto.Message, 100),
	}
	var err error
	peer.Conn, err = NewConn(c)
	return peer, err
}

func (p *Peer) GetDump() *Dump {
	return &Dump{
		Host:     p.Host,
		Port:     p.Port,
		Resolved: p.Resolved,
		State:    p.PeerState,
		Stats:    p.PeerStats,
	}
}

func (p *Peer) handle(m *Manager, l *ledger.Manager) {
	var ready sync.Once
	incoming := make(chan protocol.ExtendedMessage, 10)
	outgoing := make(chan proto.Message, 10)
	deadline := time.NewTimer(time.Second * 5)
	ping := time.NewTicker(time.Second * 30)
	go p.Conn.run(incoming, outgoing)
	for {
		select {
		case <-ping.C:
			outgoing <- protocol.NewPing()
		case <-deadline.C:
			next := p.takeSynchronous()
			if next != nil {
				p.Send(next)
				outgoing <- next
			}
			glog.Errorf("%s:Deadline hit", p.String())
			deadline.Reset(time.Minute * 2)
		case out := <-p.Outgoing:
			p.Send(out)
			outgoing <- out
		case in, ok := <-incoming:
			if !ok {
				p.UpdateStatus(Disconnected)
				return
			}
			if p.Receive(in) != PassThrough {
				next := p.takeSynchronous()
				if next != nil {
					p.Send(next)
					outgoing <- next
				}
				deadline.Reset(time.Minute * 2)
			}
			glog.V(2).Infof("%s:%s", p.String(), in.Log())
			switch msg := in.(type) {
			case *protocol.TMEndpoints:
				p.handleEndpoints(m, msg)
			case *protocol.Hello:
				p.handleHello(m, msg)
			case *protocol.TMProofWork:
				go p.handleProofOfWork(msg)
			case *protocol.TMStatusChange:
				p.handleStatusChange(msg)
				ready.Do(func() {
					go p.fillQueue(l)
				})
			case *protocol.TMLedgerData:
				go p.handleLedgerData(msg, l)
			case *protocol.TMGetObjectByHash:
				if !msg.GetQuery() {
					go p.handleGetObjectByHashReply(msg, l)
				}
			case *protocol.Ping:
				if msg.IsPing {
					p.Outgoing <- protocol.NewPong()
				}
			}
		}
	}
}

func (p *Peer) takeSynchronous() proto.Message {
	select {
	case s := <-p.synchronous:
		return s
	default:
		return nil
	}
}

func (p *Peer) fillQueue(l *ledger.Manager) {
	for {
		start, end := p.GetLedgerRange()
		request := &ledger.MissingLedgers{
			Request: &data.LedgerRange{
				Start: start,
				End:   end,
				Max:   20,
			},
			Response: make(chan data.LedgerSlice),
		}
		l.Missing <- request
		missing := <-request.Response
		glog.V(1).Infof("%s:Queueing %d-%d %+v", p.String(), start, end, missing)
		if len(missing) < 2 {
			time.Sleep(5 * time.Second)
			continue
		}
		for _, ledger := range missing {
			p.synchronous <- protocol.NewGetLedger(ledger)
		}
	}
}

func (p *Peer) handleHello(m *Manager, hello *protocol.Hello) {
	cookie, err := p.getSessionCookie()
	if err != nil {
		glog.Errorf("%s:Bad cookie: %s", p.String(), err.Error())
		p.UpdateStatus(HelloFailed)
		return
	}
	if !hello.Signature.Verify(hello.PublicKey, cookie) {
		glog.Errorf("%s:Bad signature verification: %X", p.String(), hello.Signature.DER())
		p.UpdateStatus(HelloFailed)
		return
	}
	proof, err := crypto.NewSignature(&m.key.PrivateKey, cookie)
	if err != nil {
		glog.Errorf("%s:Bad signature creation: %X", p.String(), cookie)
		p.UpdateStatus(HelloFailed)
		return
	}
	if err := p.ProcessHello(hello); err != nil {
		glog.Errorf("%s:%s", p.String(), err.Error())
		return
	}
	port, _ := strconv.ParseUint(m.Port, 10, 32)
	p.Outgoing <- &protocol.TMHello{
		FullVersion:     proto.String(m.Name),
		ProtoVersion:    proto.Uint32(MAJOR_VERSION),
		ProtoVersionMin: proto.Uint32(MINOR_VERSION),
		NodePublic:      []byte(m.PublicKey.ToJSON()),
		NodeProof:       proof.DER(),
		Ipv4Port:        proto.Uint32(uint32(port)),
		NetTime:         proto.Uint64(uint64(data.Now())),
		NodePrivate:     proto.Bool(true),
		TestNet:         proto.Bool(false),
	}
	if hello.ProofOfWork != nil {
		go p.handleProofOfWork(hello.ProofOfWork)
	}
}

func (p *Peer) handleProofOfWork(pow *protocol.TMProofWork) {
	glog.Infoln("POW!!!", pow)
	// work := crypto.NewProofOfWork(pow.Challenge, pow.Target, pow.GetIterations())
	// proof, err := work.Solve()
	// if err != nil {
	// 	glog.Errorf("%s:%s", p.String(), err.Error())
	// }
}

func (p *Peer) handleStatusChange(state *protocol.TMStatusChange) {
	p.UpdateState(state)
}

func (p *Peer) handleGetObjectByHashReply(reply *protocol.TMGetObjectByHash, l *ledger.Manager) {
	var nodes []*encoding.InnerNode
	for _, obj := range reply.GetObjects() {
		blob := append(obj.GetData(), obj.GetHash()...)
		node, err := encoding.ParseWire(blob)
		if err != nil {
			glog.Errorf("%s: %s Ledger: %d Blob: %X", p.String(), err.Error(), reply.GetSeq(), blob)
			return
		}
		if tx, ok := node.Value.(data.Transaction); ok {
			tx.SetLedgerSequence(reply.GetSeq())
			var hash data.Hash256
			copy(hash[:], obj.GetData()[len(obj.GetData())-32:])
			tx.SetHash(&hash)
			l.Incoming <- tx
		}
		if node.InnerNode != nil {
			nodes = append(nodes, node.InnerNode)
		}
	}
	if len(nodes) > 0 {
		p.synchronous <- protocol.NewGetObjects(reply.GetSeq(), nodes)
	}
}

func (p *Peer) handleLedgerData(ledgerData *protocol.TMLedgerData, l *ledger.Manager) {
	//msg.AverageLatency = ((msg.AverageLatency * Latency(msg.Successful)) + Latency(m.Time.Sub(msg.Inflight[i].Sent))) / Latency(msg.Successful+1)
	if ledgerData.GetType() != protocol.TMLedgerInfoType_liBASE {
		glog.Infof("%s: Ignoring: %s", ledgerData.Log())
		return
	}
	ledger, err := encoding.ParseLedger(bytes.NewReader(ledgerData.Nodes[0].Nodedata))
	if err != nil {
		glog.Errorf("%s: %s", p.String(), err.Error())
		return
	}
	var hash data.Hash256
	copy(hash[:], ledgerData.GetLedgerHash())
	ledger.SetHash(&hash)
	l.Incoming <- ledger
	// if ledger.TransactionHash.IsZero() {
	// 	return
	// }
	// transactions, err := encoding.ParseUnknownInnerNode(ledgerData.Nodes[2].Nodedata)
	// if err != nil {
	// 	glog.Errorf("%s: %s", p.String(), err.Error())
	// 	return
	// }
	// p.synchronous <- protocol.NewGetObjects(ledger.LedgerSequence, []*encoding.InnerNode{transactions})
}

func (p *Peer) handleEndpoints(m *Manager, msg *protocol.TMEndpoints) {
	for _, endpoint := range msg.Endpoints {
		if endpoint.GetHops() > 0 {
			port := strconv.FormatUint(uint64(endpoint.GetIpv4().GetIpv4Port()), 10)
			m.AddPeer(endpoint.GetIpv4().Host(), port, false, nil)
		}
	}
}
