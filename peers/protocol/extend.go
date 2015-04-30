package protocol

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/rubblelabs/ripple/data"
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
		for _, hash := range node.Children {
			objects = append(objects, &TMIndexedObject{Hash: hash.Bytes()})
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

// Node Version
type NodeVersion uint32

func NewNodeVersion(major, minor uint16) NodeVersion {
	return NodeVersion(major)<<16 | NodeVersion(minor)
}

func (n NodeVersion) String() string {
	return fmt.Sprintf("%d.%d", n>>16, n&0xFFFF)
}

// TMHello extension
type Hello struct {
	*TMHello
	Version    NodeVersion
	MinVersion NodeVersion
}

func (m *TMHello) Extend() (ExtendedMessage, error) {
	return &Hello{
		TMHello:    m,
		Version:    NodeVersion(m.GetProtoVersion()),
		MinVersion: NodeVersion(m.GetProtoVersionMin()),
	}, nil
}

// Untouched
func (m *TMProposeSet) Extend() (ExtendedMessage, error)         { return m, nil }
func (m *TMProofWork) Extend() (ExtendedMessage, error)          { return m, nil }
func (m *TMCluster) Extend() (ExtendedMessage, error)            { return m, nil }
func (m *TMPeers) Extend() (ExtendedMessage, error)              { return m, nil }
func (m *TMEndpoints) Extend() (ExtendedMessage, error)          { return m, nil }
func (m *TMTransaction) Extend() (ExtendedMessage, error)        { return m, nil }
func (m *TMLedgerData) Extend() (ExtendedMessage, error)         { return m, nil }
func (m *TMStatusChange) Extend() (ExtendedMessage, error)       { return m, nil }
func (m *TMHaveTransactionSet) Extend() (ExtendedMessage, error) { return m, nil }
func (m *TMValidation) Extend() (ExtendedMessage, error)         { return m, nil }

// Commands
func (m *TMGetPeers) Extend() (ExtendedMessage, error)        { return m, nil }
func (m *TMGetLedger) Extend() (ExtendedMessage, error)       { return m, nil }
func (m *TMGetObjectByHash) Extend() (ExtendedMessage, error) { return m, nil }
