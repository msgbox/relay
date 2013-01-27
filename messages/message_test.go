package messages

import (
	"code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"fmt"
	"github.com/msgbox/relay/queue"
	"github.com/msgbox/relay/structs"
	"github.com/streadway/amqp"
	"testing"
)

// Helper for creating an AMQP Connection
func createConnection() *amqp.Connection {
	conn, err := queue.Connect()
	if err != nil {
		fmt.Errorf("Error Connecting: %s", err)
	}

	return conn
}

// Helper for creating fake JSON data
// to use in the tests
func make_json() []byte {
	i := Item{Header{}, Payload{}}

	i.Header.Creator = "sender@example.com"
	i.Header.Receiver = "receiver@example.com"
	i.Payload.Title = "Test Message"
	i.Payload.Body = "Test Message Body"

	data, _ := json.Marshal(i)
	return data
}

// Tests that the Encrypt method concats the payload fields
// THIS WILL INCLUDE ENCRYPTION EVENTUALLY
func Test_Encrypt_1(t *testing.T) {
	i := Item{Payload: Payload{Title: "Title", Body: "Body"}}
	if i.encrypt() != "TitleBody" {
		t.Error("Encryption did not work as expected.")
	}
}

// Ensure Protocol Buffers are Marshaled correctly
func Test_createPB_1(t *testing.T) {
	data := make_json()
	msg, _ := createProtocolBuffer(data)
	newTest := &structs.Message{}
	proto.Unmarshal(msg, newTest)
	if newTest.GetCreator() != "sender@example.com" {
		t.Error("Protocol Buffer Not Correct")
	}
}

// Test the message is sent to an AMQP Exchange
func Test_Send_1(t *testing.T) {
	conn := createConnection()
	defer conn.Close()

	data := make_json()

	err := Send(data, conn)
	if err != nil {
		t.Errorf("Send did not work as expected: %s", err)
	}
}

// Test the message is sent to an AMQP Exchange
func Test_Receive_1(t *testing.T) {
	conn := createConnection()
	defer conn.Close()

	data := make_json()
	msg, _ := createProtocolBuffer(data)

	err := Receive(msg, conn)
	if err != nil {
		t.Errorf("Receive did not work as expected: %s", err)
	}
}
