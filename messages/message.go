package messages

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/rand"
	"encoding/hex"
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

// Generate a UUID to use as a unique Message ID
// Need to work on the UUID library to ensure UUID's
// are created correctly
func (m *Item) generateUUID() string {
	uuid, err := genUUID()
	if err != nil {
		fmt.Errorf("Error Generating UUID: %s", err)
	}

	return *&uuid
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

// This should take an incoming message in Protocol Buffer
// format and push it to an AMQP exchange
func Receive(data []byte, connection *amqp.Connection) error {
	err := queue.Publish("incoming", data, connection)
	if err != nil {
		return fmt.Errorf("Publishing Error: %s", err)
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
		Id:        proto.String(i.generateUUID()),
		Payload:   proto.String(i.encrypt()),
	}

	p, err := proto.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshaling error: ", err)
	}

	return p, nil
}

// UUID v4 Generator
// http://www.ashishbanerjee.com/home/go/go-generate-uuid
func genUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// TODO: verify the two lines implement RFC 4122 correctly
	uuid[8] = 0x80 // variant bits see page 5
	uuid[4] = 0x40 // version 4 Pseudo Random, see page 7

	return hex.EncodeToString(uuid), nil
}
