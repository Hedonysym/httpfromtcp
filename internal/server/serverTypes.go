package server

import (
	"net"
	"sync/atomic"
)

type Server struct {
	Port     string
	Listener net.Listener
	Open     atomic.Bool
}
