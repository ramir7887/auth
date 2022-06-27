package grpcserver

import (
	"context"
	"fmt"
	"github.com/ramir7887/auth-contract/auth"
	"gitlab.com/g6834/team28/auth/internal/usecase"
	"gitlab.com/g6834/team28/auth/pkg/logger"
	"google.golang.org/grpc"

	"net"
	"time"
)

const (
	_defaultShutdownTimeout = 10 * time.Second
	_defaultHost            = ""
	_defaultPort            = "4000"
)

type Server struct {
	host            string
	port            string
	error           chan error
	stop            chan struct{}
	server          *grpc.Server
	shutdownTimeout time.Duration
	logger          logger.Interface
	uc              usecase.Authentication
	auth.UnimplementedAuthServiceServer
}

func NewServer(uc usecase.Authentication, l logger.Interface, opts ...Option) *Server {
	s := &Server{
		uc:              uc,
		logger:          l.WithFields(logger.Fields{"package": "grpcserver"}),
		error:           make(chan error),
		stop:            make(chan struct{}),
		shutdownTimeout: _defaultShutdownTimeout,
		port:            _defaultPort,
		host:            _defaultHost,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	serv := grpc.NewServer()
	auth.RegisterAuthServiceServer(serv, s)

	s.start()

	return s
}

func (s *Server) start() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		s.logger.WithFields(logger.Fields{
			"error": err.Error(),
		}).Fatal("Error create Listener")
	}

	serv := grpc.NewServer()
	auth.RegisterAuthServiceServer(serv, s)
	s.server = serv

	go func() {
		s.error <- serv.Serve(lis)
		close(s.error)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.error
}

// Shutdown -.
func (s *Server) Shutdown() error {
	select {
	case <-s.error:
		return nil
	default:
	}

	close(s.stop)
	time.Sleep(s.shutdownTimeout)

	s.server.Stop()

	return nil
}

// Validate -.
func (s *Server) Validate(ctx context.Context, token *auth.Token) (*auth.Token, error) {
	l := s.logger.WithFields(logger.Fields{
		"method": "Server.Validate",
	})
	l.Info("Start handle Validate")
	defer l.Info("End handle Validate")
	accessToken, refreshToken, err := s.uc.Validate(ctx, token.AccessToken, token.RefreshToken)
	if err != nil {
		l.WithFields(logger.Fields{
			"error": err.Error(),
		}).Error("Error Validate")
		return nil, err
	}

	return &auth.Token{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
