package protocol

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strings"
)

var messageFactory = map[MessageType]func() Message{
	MessageType_mtHELLO:              func() Message { return &TMHello{} },
	MessageType_mtERROR_MSG:          func() Message { return &TMErrorMsg{} },
	MessageType_mtPING:               func() Message { return &TMPing{} },
	MessageType_mtPROOFOFWORK:        func() Message { return &TMProofWork{} },
	MessageType_mtCLUSTER:            func() Message { return &TMCluster{} },
	MessageType_mtGET_CONTACTS:       func() Message { return &TMGetContacts{} },
	MessageType_mtCONTACT:            func() Message { return &TMContact{} },
	MessageType_mtGET_PEERS:          func() Message { return &TMGetPeers{} },
	MessageType_mtPEERS:              func() Message { return &TMPeers{} },
	MessageType_mtENDPOINTS:          func() Message { return &TMEndpoints{} },
	MessageType_mtSEARCH_TRANSACTION: func() Message { return &TMSearchTransaction{} },
	MessageType_mtGET_ACCOUNT:        func() Message { return &TMGetAccount{} },
	MessageType_mtACCOUNT:            func() Message { return &TMAccount{} },
	MessageType_mtTRANSACTION:        func() Message { return &TMTransaction{} },
	MessageType_mtGET_LEDGER:         func() Message { return &TMGetLedger{} },
	MessageType_mtLEDGER_DATA:        func() Message { return &TMLedgerData{} },
	MessageType_mtPROPOSE_LEDGER:     func() Message { return &TMProposeSet{} },
	MessageType_mtSTATUS_CHANGE:      func() Message { return &TMStatusChange{} },
	MessageType_mtHAVE_SET:           func() Message { return &TMHaveTransactionSet{} },
	MessageType_mtGET_VALIDATIONS:    func() Message { return &TMGetValidations{} },
	MessageType_mtVALIDATION:         func() Message { return &TMValidation{} },
	MessageType_mtGET_OBJECTS:        func() Message { return &TMGetObjectByHash{} },
}

var typeFactory = map[string]MessageType{
	"TMHello":              MessageType_mtHELLO,
	"TMErrorMsg":           MessageType_mtERROR_MSG,
	"TMPing":               MessageType_mtPING,
	"TMProofWork":          MessageType_mtPROOFOFWORK,
	"TMCluster":            MessageType_mtCLUSTER,
	"TMGetContacts":        MessageType_mtGET_CONTACTS,
	"TMContact":            MessageType_mtCONTACT,
	"TMGetPeers":           MessageType_mtGET_PEERS,
	"TMPeers":              MessageType_mtPEERS,
	"TMEndpoints":          MessageType_mtENDPOINTS,
	"TMSearchTransaction":  MessageType_mtSEARCH_TRANSACTION,
	"TMGetAccount":         MessageType_mtGET_ACCOUNT,
	"TMAccount":            MessageType_mtACCOUNT,
	"TMTransaction":        MessageType_mtTRANSACTION,
	"TMGetLedger":          MessageType_mtGET_LEDGER,
	"TMLedgerData":         MessageType_mtLEDGER_DATA,
	"TMProposeSet":         MessageType_mtPROPOSE_LEDGER,
	"TMStatusChange":       MessageType_mtSTATUS_CHANGE,
	"TMHaveTransactionSet": MessageType_mtHAVE_SET,
	"TMGetValidations":     MessageType_mtGET_VALIDATIONS,
	"TMValidation":         MessageType_mtVALIDATION,
	"TMGetObjectByHash":    MessageType_mtGET_OBJECTS,
}

type Header struct {
	Length      int32
	MessageType int16
}

type Decoder struct {
	r      io.Reader
	buffer *proto.Buffer
	buf    bytes.Buffer
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r:      r,
		buffer: proto.NewBuffer(nil),
	}
}

func (dec *Decoder) Decode() (ExtendedMessage, error) {
	var header Header
	if err := binary.Read(dec.r, binary.BigEndian, &header); err != nil {
		return nil, fmt.Errorf("Protocol Decode: Read Header: %s", err.Error())
	}
	factory, ok := messageFactory[MessageType(header.MessageType)]
	if !ok {
		return nil, fmt.Errorf("Protocol Decode: Unknown Message Type: %d", header.MessageType)
	}
	msg := factory()
	dec.buf.Reset()
	if _, err := io.CopyN(&dec.buf, dec.r, int64(header.Length)); err != nil {
		return nil, fmt.Errorf("Protocol Decode: Read Body: %s", err.Error())
	}
	dec.buffer.SetBuf(dec.buf.Bytes())
	if err := dec.buffer.Unmarshal(msg); err != nil {
		return nil, fmt.Errorf("Protocol Decode: Unmarshal: %s", err.Error())
	}
	return msg.Extend()
}

type Encoder struct {
	w      io.Writer
	buffer *proto.Buffer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w:      w,
		buffer: proto.NewBuffer(nil),
	}
}
func (enc *Encoder) Encode(msg proto.Message) error {
	enc.buffer.Reset()
	name := strings.TrimPrefix(reflect.TypeOf(msg).String(), "*protocol.")
	typ, ok := typeFactory[name]
	if !ok {
		return fmt.Errorf("Protocol Encode: Unknown Message Type: %s", name)
	}
	if err := enc.buffer.Marshal(msg); err != nil {
		return fmt.Errorf("Protocol Encode: Marshal: %s", err.Error())
	}
	var header Header
	header.MessageType = int16(typ)
	header.Length = int32(len(enc.buffer.Bytes()))
	if err := binary.Write(enc.w, binary.BigEndian, header); err != nil {
		return fmt.Errorf("Protocol Encode: Write Header: %s", err.Error())
	}
	if _, err := enc.w.Write(enc.buffer.Bytes()); err != nil {
		return fmt.Errorf("Protocol Encode: Write Body: %s", err.Error())
	}
	return nil
}
