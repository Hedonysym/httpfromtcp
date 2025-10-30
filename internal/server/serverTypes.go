package server

import (
	"net"
	"sync/atomic"

	"github.com/Hedonysym/httpfromtcp/internal/request"
	"github.com/Hedonysym/httpfromtcp/internal/response"
)

type Server struct {
	Port     string
	Listener net.Listener
	Open     atomic.Bool
}

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	StatusCode int
	Err        string
}
