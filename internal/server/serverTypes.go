package server

import (
	"io"
	"net"
	"sync/atomic"

	"github.com/Hedonysym/httpfromtcp/internal/request"
)

type Server struct {
	Port     string
	Listener net.Listener
	Open     atomic.Bool
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode int
	Err        string
}
