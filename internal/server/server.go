package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

func Serve(port int) (*Server, error) {
	portStr := ":" + strconv.Itoa(port)
	l, err := net.Listen("tcp", portStr)
	if err != nil {
		return nil, err
	}

	var open atomic.Bool

	s := &Server{
		Port:     portStr,
		Listener: l,
		Open:     open,
	}
	s.Open.Store(true)
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.Open.Store(false)
	return s.Listener.Close()
}

func (s *Server) listen() {
	for s.Open.Load() {
		conn, err := s.Listener.Accept()
		if err != nil && s.Open.Load() {
			log.Fatal(err)
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!"))
	conn.Close()
}
