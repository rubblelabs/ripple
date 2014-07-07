package protocol

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/rubblelabs/ripple/data"
	"strings"
)

func b2h(b []byte) string {
	return hex.EncodeToString(b)
}

func short(s string) string {
	if len(s) > 10 {
		return s[:4] + ".." + s[len(s)-4:]
	}
	return s
}

func (m *Ping) Log() string {
	netTime, pingTime := data.NewRippleTime(uint32(m.GetNetTime())), data.NewRippleTime(uint32(m.GetPingTime()))
	return fmt.Sprintf("Ping: %s Seq: %d Time: %s/%s ", m.GetType(), m.GetSeq(), netTime, pingTime)
}

func (m *TMProposeSet) Log() string {
	buf := new(bytes.Buffer)
	closeTime := data.NewRippleTime(m.GetCloseTime())
	fmt.Fprintf(buf, "Proposal: %d Closed: %s Hash: %s ", m.GetProposeSeq(), closeTime.Short(), short(b2h(m.GetCurrentTxHash())))
	fmt.Fprintf(buf, "Sig: %s ", short(b2h(m.GetSignature())))
	for _, tx := range m.GetAddedTransactions() {
		fmt.Fprintf(buf, "\nAdded: %s", short(b2h(tx)))
	}
	for _, tx := range m.GetRemovedTransactions() {
		fmt.Fprintf(buf, "\nRemoved: %s", short(b2h(tx)))
	}
	return buf.String()
}

func (m *Hello) Log() string {
	return fmt.Sprintf("Hello: %s (V:%s Min:%s) Port: %d Proof:%v", m.GetFullVersion(), m.Version, m.MinVersion, m.GetIpv4Port(), m.GetProofOfWork())
}

func (m *TMValidation) Log() string {
	return fmt.Sprintf("Validation: Checked: %t Validation: %s", m.GetCheckedSignature(), short(b2h(m.Validation)))
}

func (m *TMTransaction) Log() string {
	received := data.NewRippleTime(uint32(m.GetReceiveTimestamp()))
	return fmt.Sprintf("Transaction: Status: %s Checked: %t Received: %s Raw:%s", m.GetStatus(), m.GetCheckedSignature(), received.String(), short(b2h(m.RawTransaction)))
}

func (m *TMStatusChange) Log() string {
	t := data.NewRippleTime(uint32(m.GetNetworkTime()))
	return fmt.Sprintf("Status: %s Event: %s Seq: %d (%d-%d) Time: %s", m.GetNewStatus(), m.GetNewEvent(), m.GetLedgerSeq(), m.GetFirstSeq(), m.GetLastSeq(), t.Short())
}

func (m *TMHaveTransactionSet) Log() string {
	return fmt.Sprintf("HaveTransactionSet: %s %s", m.GetStatus(), short(b2h(m.GetHash())))
}

func (m *TMGetPeers) Log() string {
	return "GetPeers:"
}

func (m *TMEndpoints) Log() string {
	var endpoints []string
	for _, ep := range m.Endpoints {
		endpoints = append(endpoints, ep.GetIpv4().Host())
	}
	return fmt.Sprintf("Endpoints: %s", strings.Join(endpoints, ","))
}

func (m *TMGetLedger) Log() string {
	var nodes []string
	for _, n := range m.GetNodeIDs() {
		nodes = append(nodes, short(b2h(n)))
	}
	return fmt.Sprintf("GetLedger: %d Hash:%s Type:%s Query:%s Cookie:%d Nodes: %s", m.GetLedgerSeq(), short(b2h(m.GetLedgerHash())), m.GetLtype().String(), m.GetQueryType().String(), m.GetRequestCookie(), strings.Join(nodes, ","))
}

func (m *TMLedgerData) Log() string {
	if m.Error != nil {
		return fmt.Sprintf("LedgerData: %s", m.GetError().String())
	}
	var nodes []string
	for _, node := range m.GetNodes() {
		nodes = append(nodes, short(b2h(node.GetNodedata())))
	}
	return fmt.Sprintf("LedgerData: %d Hash: %s Cookie: %d Nodes: %s", m.GetLedgerSeq(), short(b2h(m.GetLedgerHash())), m.GetRequestCookie(), strings.Join(nodes, ","))
}

func (m *TMGetObjectByHash) Log() string {
	var objects []string
	for _, obj := range m.GetObjects() {
		objects = append(objects, fmt.Sprintf("%s:%s", short(b2h(obj.GetHash())), short(b2h(obj.GetData()))))
	}
	return fmt.Sprintf("GetObjectByHash: %d %s", m.GetSeq(), strings.Join(objects, ","))
}

func (m *TMProofWork) Log() string { return m.String() }
func (m *TMCluster) Log() string   { return m.String() }
func (m *TMPeers) Log() string     { return m.String() }
