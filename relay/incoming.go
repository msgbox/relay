// Incoming Relay Handler
//
// Handles an incoming message from another relay.
// It simply accepts the message and writes it into a queue
// for further processing.

package handler

import (
	"github.com/msgbox/queue"
	"log"
)

type Incoming struct {
	Sess *Session
}

// Read from the Session buffer and send to a handler function
func (in *Incoming) Read() {
	defer in.Sess.rwc.Close()

	for {
		line, err := readln(in.Sess.br)
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}

		in.handleMessage(line)
	}
}

// Connect to an AMQP Exchange and write the message
func (in *Incoming) handleMessage(body []byte) {

	// Create an AMQP Connection
	conn, c_err := queue.Connect()
	if c_err != nil {
		// Handle Error
	}
	defer conn.Close()

	err := queue.Publish("incoming", body, conn)
	if err != nil {
		// Handle error
	}
}
