package outgoing_message

import (
	"code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"fmt"
	"github.com/msgbox/relay/queue"
	"github.com/msgbox/relay/structs"
	"github.com/streadway/amqp"
)

// Interface for encryptors
type Encrypter interface {
	encrypt() []byte
}

type Item struct {
	Header  Header
	Payload Payload
}

type Header struct {
	Creator    string
	Receiver   string
	Created_At int64
	MessageID  string
}

type Payload struct {
	Title string
	Body  string
}

// Eventually this will encrypt the payload
// for now just concat the title and body
// Should return a []byte when encryption is
// added to meet the Encrypter interface
func (m *Item) encrypt() string {
	s := ""
	s += *&m.Payload.Title
	s += *&m.Payload.Body
	return s
}

// This should take a JSON Input and Marshall it into
// a protocol buffer that can be pushed to a AMQP Exchange
func Send(data []byte, connection *amqp.Connection) error {
	pb, err := createProtocolBuffer(data)
	if err != nil {
		return fmt.Errorf("Protocol Buffer Error: %s", err)
	}

	// Send pb to AMQP Exchange
	p_err := queue.Publish("outgoing", pb, connection)
	if p_err != nil {
		return fmt.Errorf("Publishing Error: %s", p_err)
	}

	return nil
}

// Takes a byte slice, usually of JSON data, and converts
// it into a Protocol Buffer.
//
// @param {String} data
// @return {[]byte}, {error}
// @api private
func createProtocolBuffer(data []byte) ([]byte, error) {

	// Marshall the JSON input into an Item struct
	var i Item
	err := json.Unmarshal(data, &i)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling error:", err)
	}

	// Now we want to create a Message from the Item
	msg := &structs.Message{
		Creator:   proto.String(*&i.Header.Creator),
		Receiver:  proto.String(*&i.Header.Receiver),
		CreatedAt: proto.Int64(*&i.Header.Created_At),
		Id:        proto.String(*&i.Header.MessageID),
		Payload:   proto.String(i.encrypt()),
	}

	p, err := proto.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshaling error: ", err)
	}

	return p, nil
}
