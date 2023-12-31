package http

import (
	"net/http"
	"time"
)

// ServerOption is HTTP server option.
type ServerOption func(*Server)

// Network with server network
func Network(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

// Address with server address
func Address(address string) ServerOption {
	return func(s *Server) {
		s.address = address
	}
}

// Timeout with server timeout
func Timeout(t time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = t
	}
}

// Handler with http handler
func Handler(handler http.Handler) ServerOption {
	return func(s *Server) {
		s.handler = handler
	}
}
