// Outgoing Relay Handler
//
// Handles an outgoing message.
// It serves as a relay that a worker can dial and have the
// relay handle delivering the message.

package handler

import (
	"fmt"
	"github.com/msgbox/message"
	"io/ioutil"
	"net"
	"strings"
)

// Read from the Session buffer and send to a handler function
func ReadOutgoing(conn net.Conn) {
	// close connection on exit
	defer conn.Close()

	result, err := ioutil.ReadAll(conn)
	checkError(err)

	err = sendMessage(result)
	checkError(err)
}

func sendMessage(result []byte) error {
	// Lookup Receiving Node
	node := lookupAddress(result)

	conn, err := net.DialTCP("tcp", nil, node)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(result)
	if err != nil {
		return err
	}

	return nil
}

// Get it working with with IP addresses for now
func lookupAddress(data []byte) *net.TCPAddr {
	port := ":7834" // Hardcoded for now
	msg := messages.Parse(data)
	receiver := msg.GetReceiver()
	addr := strings.Split(receiver, "@")
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr[1]+port)
	checkError(err)
	return tcpAddr
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s", err.Error())
	}
}
