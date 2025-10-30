package server

import (
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/Hedonysym/httpfromtcp/internal/request"
	"github.com/Hedonysym/httpfromtcp/internal/response"
)

func Serve(port int, h Handler) (*Server, error) {
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
	go s.listen(h)
	return s, nil
}

func (s *Server) Close() error {
	s.Open.Store(false)
	return s.Listener.Close()
}

func (s *Server) listen(h Handler) {
	for s.Open.Load() {
		conn, err := s.Listener.Accept()
		if err != nil && s.Open.Load() {
			log.Fatal(err)
		}
		go s.handle(conn, h)
	}
}

func (s *Server) handle(conn net.Conn, h Handler) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatal(err)
	}
	resp := response.NewWriter(conn)
	h(resp, req)
}

func WriteError(w io.Writer, e HandlerError) error {
	head := response.GetDefaultHeaders(len(e.Err))
	err := response.WriteStatus(w, response.StatusCode(e.StatusCode))
	if err != nil {
		return err
	}
	err = response.WriteHeaders(w, head)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(e.Err))
	if err != nil {
		return err
	}
	return nil
}
