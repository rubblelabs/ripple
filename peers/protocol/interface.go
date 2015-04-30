package protocol

import (
	"github.com/golang/protobuf/proto"
)

type Message interface {
	proto.Message
	Extend() (ExtendedMessage, error)
}

type ExtendedMessage interface {
	proto.Message
	Log() string
}
