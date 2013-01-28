package messages

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/msgbox/relay/structs"
)

// Parses a Message from the queue

func Parse(data []byte) *structs.Message {
	msg := &structs.Message{}
	proto.Unmarshal(data, msg)

	return msg
}
