package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.httpfromtcp/internal/response"
)

type stateofserver int

const (
	stateinit stateofserver = iota
	statedone
)

type Server struct {
	Stateofserver stateofserver
	listene       net.Listener
	closed        atomic.Bool
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listene.Close()
}

func (s *Server) Listen() {
	for {
		conn, err := s.listene.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			fmt.Printf("error in accepting %q", err)
			continue
		}
		go s.handel(conn)
	}
}

func (s *Server) handel(conn net.Conn) {
	defer conn.Close()
	body := "hello world!\n"
	if err := response.WriteStatusLine(conn, response.ServerOk); err != nil {
		return
	}
	h := response.GetDefaultHeaders(len(body))
	if err := response.WriteHeaders(conn, h); err != nil {
		return
	}

	conn.Write([]byte(body))
}

func Serv(port int) (*Server, error) {
	var serv Server
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	serv.listene = ln
	go serv.Listen()
	return &serv, nil
}
