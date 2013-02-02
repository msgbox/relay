// Incoming Relay Handler
//
// Handles an incoming message from another relay.
// It simply accepts the message and writes it into a queue
// for further processing.

package handler

import (
	"github.com/msgbox/queue"
	"github.com/streadway/amqp"
	"io/ioutil"
	"net"
)

func ReadIncoming(conn net.Conn, queueConn *amqp.Connection) {
	// close connection on exit
	defer conn.Close()

	result, err := ioutil.ReadAll(conn)
	checkError(err)

	err = queue.Publish("incoming", result, queueConn)
	checkError(err)
}
