package relay

import (
	"fmt"
	"github.com/msgbox/queue"
	"github.com/msgbox/relay/relay"
	"github.com/streadway/amqp"
	"net"
)

type session struct {
	queue *amqp.Connection
}

var s *session

func init() {
	// Create an AMQP Connection
	queueConn, err := queue.Connect()
	checkError(err)

	s = &session{
		queue: queueConn,
	}
}

func ListenIncoming(addr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handler.ReadIncoming(conn, s.queue)
	}
}

func ListenOutgoing(addr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handler.ReadOutgoing(conn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s", err.Error())
	}
}
