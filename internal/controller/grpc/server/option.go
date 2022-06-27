package grpcserver

import (
	"time"
)

// Option -.
type Option func(*Server)

// Port -.
func Port(port string) Option {
	return func(s *Server) {
		s.port = port
	}
}

// Host -.
func Host(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

// Timeout -.
func Timeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
