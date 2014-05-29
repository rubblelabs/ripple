package protocol

import (
	"code.google.com/p/goprotobuf/proto"
)

type Message interface {
	proto.Message
	Extend() (ExtendedMessage, error)
}

type ExtendedMessage interface {
	proto.Message
	Log() string
}
