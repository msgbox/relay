// Outgoing Relay Handler
//
// Handles an outgoing message.
// It serves as a relay that a worker can dial and have the
// relay handle delivering the message.

package handler

import (
	"fmt"
	"io/ioutil"
	"net"
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
	node := "127.0.0.1:7834"

	tcpAddr, err := net.ResolveTCPAddr("tcp4", node)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
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

func checkError(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s", err.Error())
	}
}
