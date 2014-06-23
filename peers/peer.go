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

var (
	minVersion = protocol.NewNodeVersion(1, 2)
	maxVersion = protocol.NewNodeVersion(1, 2)
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
	sync        ledger.Sync
}

func NewPeer(c *PeerConnection, sync ledger.Sync) (*Peer, error) {
	peer := &Peer{
		PeerState:   NewPeerState(c),
		PeerStats:   NewPeerStats(),
		Outgoing:    make(chan proto.Message, 10),
		synchronous: make(chan proto.Message, 100),
		sync:        sync,
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

func (p *Peer) handle(m *Manager) {
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
					go p.fillQueue()
				})
			case *protocol.TMHaveTransactionSet:
				go p.handleHaveTransactionSet(msg)
			case *protocol.TMProposeSet:
				go p.handleProposeSet(msg)
			case *protocol.TMValidation:
				go p.handleValidation(msg)
			case *protocol.TMTransaction:
				go p.handleTransaction(msg)
			case *protocol.TMLedgerData:
				go p.handleLedgerData(msg)
			case *protocol.TMGetObjectByHash:
				if !msg.GetQuery() {
					go p.handleGetObjectByHashReply(msg)
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

func (p *Peer) fillQueue() {
	for {
		start, end := p.GetLedgerRange()
		r := &data.LedgerRange{
			Start: start,
			End:   end,
			Max:   20,
		}
		work := p.sync.Missing(r)
		glog.V(1).Infof("%s:Queueing %d-%d %+v", p.String(), start, end, work.MissingLedgers)
		if len(work.MissingLedgers) < 2 {
			time.Sleep(5 * time.Second)
			continue
		}
		for _, ledger := range work.MissingLedgers {
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
	pubKey, err := crypto.ParsePublicKeyFromHash(hello.GetNodePublic())
	if err != nil {
		glog.Errorf("Bad public key: %X", hello.GetNodePublic())
		p.UpdateStatus(HelloFailed)
		return
	}
	ok, err := crypto.Verify(pubKey.SerializeUncompressed(), hello.GetNodeProof(), cookie)
	if !ok {
		glog.Errorf("%s:Bad signature: %X public key: %X hash: %X", p.String(), hello.GetNodeProof(), hello.GetNodePublic(), cookie)
		p.UpdateStatus(HelloFailed)
		return
	}
	if err != nil {
		glog.Errorf("%s:Bad signature verification: %s", p.String(), err.Error())
		p.UpdateStatus(HelloFailed)
		return
	}
	proof, err := m.Key.Sign(cookie)
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
		ProtoVersion:    proto.Uint32(uint32(maxVersion)),
		ProtoVersionMin: proto.Uint32(uint32(minVersion)),
		NodePublic:      []byte(m.PublicKey.String()),
		NodeProof:       proof,
		Ipv4Port:        proto.Uint32(uint32(port)),
		NetTime:         proto.Uint64(uint64(data.Now().Uint32())),
		NodePrivate:     proto.Bool(true),
		TestNet:         proto.Bool(false),
	}
	if hello.ProofOfWork != nil {
		go p.handleProofOfWork(hello.ProofOfWork)
	}
}

func (p *Peer) handleProposeSet(proposeSet *protocol.TMProposeSet) {
	proposal := &data.Proposal{
		Sequence:  proposeSet.GetProposeSeq(),
		CloseTime: *data.NewRippleTime(proposeSet.GetCloseTime()),
		Signature: proposeSet.GetSignature(),
	}
	copy(proposal.LedgerHash[:], proposeSet.CurrentTxHash)
	copy(proposal.PreviousLedger[:], proposeSet.GetPreviousledger())
	copy(proposal.PublicKey[:], proposeSet.GetNodePubKey())
	p.sync.Submit([]data.Hashable{proposal})
}

func (p *Peer) handleValidation(validation *protocol.TMValidation) {
	v, err := data.NewDecoder(bytes.NewReader(validation.GetValidation())).Validation()
	if err != nil {
		glog.Errorln(err.Error())
		return
	}
	p.sync.Submit([]data.Hashable{v})
}

func (p *Peer) handleTransaction(tx *protocol.TMTransaction) {
	node, err := data.NewDecoder(bytes.NewReader(tx.GetRawTransaction())).Transaction()
	glog.Infof("%X", tx.GetRawTransaction())
	if err != nil {
		glog.Errorln(err.Error())
		return
	}
	p.sync.Submit([]data.Hashable{node})
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
	p.sync.Current(state.GetLedgerSeq())
}

func (p *Peer) handleHaveTransactionSet(txSet *protocol.TMHaveTransactionSet) {
	// glog.Infof("%s %X", txSet.GetStatus(), txSet.GetHash())
}

func (p *Peer) handleGetObjectByHashReply(reply *protocol.TMGetObjectByHash) {
	var nodes []data.Hashable
	typ := data.NT_ACCOUNT_NODE
	if reply.GetType() == protocol.TMGetObjectByHash_otTRANSACTION_NODE {
		typ = data.NT_TRANSACTION_NODE
	}
	for _, obj := range reply.GetObjects() {
		blob := append(obj.GetData(), obj.GetHash()...)
		node, err := data.NewDecoder(bytes.NewReader(blob)).Wire(typ)
		if err != nil {
			glog.Errorf("%s: %s Ledger: %d Blob: %X", p.String(), err.Error(), reply.GetSeq(), blob)
			return
		}
		glog.Infoln(node)
		nodes = append(nodes, node)
	}
	p.sync.Submit(nodes)
}

func (p *Peer) handleLedgerData(ledgerData *protocol.TMLedgerData) {
	if ledgerData.GetType() != protocol.TMLedgerInfoType_liBASE {
		glog.Infof("%s: Ignoring: %s", ledgerData.Log())
		return
	}
	ledger, err := data.NewDecoder(bytes.NewReader(ledgerData.Nodes[0].Nodedata)).Ledger()
	if err != nil {
		glog.Errorf("%s: %s", p.String(), err.Error())
		return
	}
	glog.Infoln(ledger)
	p.sync.Submit([]data.Hashable{ledger})
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
