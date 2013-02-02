package relay

import (
	"github.com/msgbox/relay/relay"
	"net"
)

func ListenIncoming(addr string) net.Listener {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		// handle error
	}

	for {
		rw, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}

		sess, err := handler.NewSession(rw)
		if err != nil {
			continue
		}

		handler := handler.Incoming{Sess: sess}

		go handler.Read()
	}

	return ln
}

func ListenOutgoing(addr string) net.Listener {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		// handle error
	}

	for {
		rw, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}

		sess, err := handler.NewSession(rw)
		if err != nil {
			continue
		}

		handler := handler.Outgoing{Sess: sess}

		go handler.Read()
	}

	return ln
}
