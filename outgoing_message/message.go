package outgoing_message

import (
	"code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"github.com/msgbox/relay/structs"
	"log"
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
func Send(data []byte) {
	pb := createProtocolBuffer(data)

	// Send pb to AMQP Exchange
	log.Println(string(pb))
}

// Takes a byte slice, usually of JSON data, and converts
// it into a Protocol Buffer.
//
// @param {String} data
// @return {[]byte}
// @api private
func createProtocolBuffer(data []byte) []byte {

	// Marshall the JSON input into an Item struct
	var i Item
	err := json.Unmarshal(data, &i)
	if err != nil {
		log.Fatal("unmarshaling error:", err)
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
		log.Fatal("marshaling error: ", err)
	}

	return p
}
