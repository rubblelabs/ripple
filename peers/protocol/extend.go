package protocol

import (
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"fmt"
	"github.com/conformal/btcec"
	"github.com/donovanhide/ripple/crypto"
	"github.com/donovanhide/ripple/data"
	"net"
)

// Simple factories

func NewGetPeers() *TMGetPeers {
	return &TMGetPeers{
		DoWeNeedThis: proto.Uint32(0),
	}
}

func NewPing() *TMPing {
	return &TMPing{Type: TMPing_ptPING.Enum()}
}

func NewPong() *TMPing {
	return &TMPing{Type: TMPing_ptPONG.Enum()}
}

func NewGetLedger(sequence uint32) *TMGetLedger {
	return &TMGetLedger{
		Itype:     TMLedgerInfoType_liBASE.Enum(),
		LedgerSeq: proto.Uint32(sequence),
	}
}

func NewGetClosedLedger() *TMGetLedger {
	return &TMGetLedger{
		Itype: TMLedgerInfoType_liBASE.Enum(),
		Ltype: TMLedgerType_ltCLOSED.Enum(),
	}
}

func NewGetLedgerTransactions(sequence uint32, nodeids [][]byte) *TMGetLedger {
	return &TMGetLedger{
		Itype:     TMLedgerInfoType_liTX_NODE.Enum(),
		LedgerSeq: proto.Uint32(sequence),
		NodeIDs:   nodeids,
	}
}

func NewGetObjects(sequence uint32, nodes []*data.InnerNode) *TMGetObjectByHash {
	var objects []*TMIndexedObject
	for _, node := range nodes {
		for _, hash := range node.Hashes() {
			objects = append(objects, &TMIndexedObject{Hash: hash})
		}
	}
	return &TMGetObjectByHash{
		Type:    TMGetObjectByHash_otTRANSACTION.Enum(),
		Query:   proto.Bool(true),
		Objects: objects,
		Seq:     proto.Uint32(sequence),
	}
}

func (endpoint *TMIPv4Endpoint) Host() string {
	ip := make(net.IP, 4)
	binary.LittleEndian.PutUint32(ip, endpoint.GetIpv4())
	return ip.String()
}

// Simple utility methods
func (endpoint *TMIPv4Endpoint) Address() string {
	return fmt.Sprintf("%s:%d", endpoint.Host(), endpoint.GetIpv4Port())
}

// TMPing extension
type Ping struct {
	*TMPing
	IsPing bool
}

func (m *TMPing) Extend() (ExtendedMessage, error) {
	return &Ping{
		TMPing: m,
		IsPing: m.GetType() == TMPing_ptPING,
	}, nil
}

// TMHello extension
type Hello struct {
	*TMHello
	Version    string
	MinVersion string
	PublicKey  []byte
	Signature  []byte
}

func (m *TMHello) Extend() (ExtendedMessage, error) {
	key, err := crypto.ParsePublicKeyFromHash(m.NodePublic)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.ParseSignature(m.NodeProof)
	if err != nil {
		return nil, err
	}
	return &Hello{
		TMHello:    m,
		Version:    fmt.Sprintf("%d.%d", m.GetProtoVersion()>>16, m.GetProtoVersion()&0xFFFF),
		MinVersion: fmt.Sprintf("%d.%d", m.GetProtoVersionMin()>>16, m.GetProtoVersionMin()&0xFFFF),
		PublicKey:  key,
		Signature:  sig,
	}, nil
}

// TMProposeSet extension
type ProposeSet struct {
	*TMProposeSet
	PublicKey  *btcec.PublicKey
	Signature  *crypto.Signature
	NodePublic crypto.Hash
}

func (m *TMProposeSet) Extend() (ExtendedMessage, error) {
	key, err := crypto.ParsePublicKey(m.NodePubKey)
	if err != nil {
		return nil, err
	}
	public, err := crypto.NewRipplePublicNode(key.SerializeCompressed())
	if err != nil {
		return nil, err
	}
	sig, err := crypto.ParseSignature(m.Signature)
	if err != nil {
		return nil, err
	}
	return &ProposeSet{
		TMProposeSet: m,
		PublicKey:    key,
		Signature:    sig,
		NodePublic:   public,
	}, nil
}

//
func (m *TMProofWork) Extend() (ExtendedMessage, error)          { return m, nil }
func (m *TMCluster) Extend() (ExtendedMessage, error)            { return m, nil }
func (m *TMPeers) Extend() (ExtendedMessage, error)              { return m, nil }
func (m *TMEndpoints) Extend() (ExtendedMessage, error)          { return m, nil }
func (m *TMTransaction) Extend() (ExtendedMessage, error)        { return m, nil }
func (m *TMLedgerData) Extend() (ExtendedMessage, error)         { return m, nil }
func (m *TMStatusChange) Extend() (ExtendedMessage, error)       { return m, nil }
func (m *TMHaveTransactionSet) Extend() (ExtendedMessage, error) { return m, nil }

// Commands
func (m *TMGetPeers) Extend() (ExtendedMessage, error)        { return m, nil }
func (m *TMGetLedger) Extend() (ExtendedMessage, error)       { return m, nil }
func (m *TMGetObjectByHash) Extend() (ExtendedMessage, error) { return m, nil }

//Not implemented in rippled
func (m *TMSearchTransaction) Extend() (ExtendedMessage, error) { return m, nil }
func (m *TMErrorMsg) Extend() (ExtendedMessage, error)          { return m, nil }
func (m *TMGetAccount) Extend() (ExtendedMessage, error)        { return m, nil }
func (m *TMAccount) Extend() (ExtendedMessage, error)           { return m, nil }
func (m *TMGetValidations) Extend() (ExtendedMessage, error)    { return m, nil }
func (m *TMValidation) Extend() (ExtendedMessage, error)        { return m, nil }
func (m *TMGetContacts) Extend() (ExtendedMessage, error)       { return m, nil }
func (m *TMContact) Extend() (ExtendedMessage, error)           { return m, nil }
