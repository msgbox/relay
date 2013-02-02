// Incoming Relay Handler
//
// Handles an incoming message from another relay.
// It simply accepts the message and writes it into a queue
// for further processing.

package handler

import (
	"bufio"
	"log"
	"net"
)

type Outgoing struct {
	Sess *Session
}

// Read from the Session buffer and send to a handler function
func (out *Outgoing) Read() {
	defer out.Sess.rwc.Close()

	for {
		line, err := readln(out.Sess.br)
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}

		out.handleMessage(sl)
	}
}

func readln(r *bufio.Reader) ([]byte, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return ln, err
}

// Connect to an AMQP Exchange and write the message
func (out *Outgoing) handleMessage(body []byte) {

	// Find Receiving Node
	node := "127.0.0.1:7834"

	// Dial Net Connection
	conn, err := net.Dial("tcp", node)
	if err != nil {
		// Handle error
	}
	defer conn.Close()

	// Transfer Payload
	status, err := conn.Write(body)
	if err != nil {
		// Handle Error
	}

	if status < 0 {
		// Figure out what status is returning
	}

}
