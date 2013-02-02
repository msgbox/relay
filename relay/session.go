package handler

import (
	"bufio"
	"net"
)

type Session struct {
	rwc net.Conn
	br  *bufio.Reader
	bw  *bufio.Writer
}

func NewSession(rwc net.Conn) (s *Session, err error) {
	s = &Session{
		rwc: rwc,
		br:  bufio.NewReader(rwc),
		bw:  bufio.NewWriter(rwc),
	}

	return
}
