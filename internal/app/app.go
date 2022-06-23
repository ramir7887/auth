// Package app configures and runs application.
package app

import (
	"fmt"
	"gitlab.com/g6834/team28/auth/internal/controller/http/profile"
	v2 "gitlab.com/g6834/team28/auth/internal/controller/http/v2"
	"gitlab.com/g6834/team28/auth/internal/entity"
	"gitlab.com/g6834/team28/auth/internal/repository/inmemory"
	"gitlab.com/g6834/team28/auth/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"gitlab.com/g6834/team28/auth/config"
	"gitlab.com/g6834/team28/auth/internal/usecase"
	"gitlab.com/g6834/team28/auth/pkg/httpserver"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	var err error
	l := logger.New(cfg.Log.Level)

	// Repository
	// 1. Create repository
	users := make([]entity.User, 0, len(cfg.Users))
	for _, u := range cfg.Users {
		user := entity.User{
			Name:     u.Name,
			Password: u.Password,
		}
		users = append(users, user)
	}
	repository := inmemory.New(users)

	// Use case
	// 1. Create UseCase
	authenticationUseCase := usecase.New(repository)

	handler := mux.NewRouter()
	// HTTP Server
	prof := profile.New(cfg.Profile.Enabled, cfg.Profile.Login, cfg.Profile.Password, l)
	prof.NewRouter(handler.PathPrefix("/debug/pprof/").Subrouter())

	// 1. Create Router for Postman tests
	v2.NewRouter(handler, l, authenticationUseCase)

	// 2. Create Router for Requirements
	//handler := mux.NewRouter()
	//v1.NewRouter(handler, l, authenticationUseCase)

	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	l.WithFields(logger.Fields{
		"package": "app",
		"method":  "Run",
	}).Infof("Http server started at :%s", cfg.HTTP.Port)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: "+s.String(), nil)
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err).Error(), nil)
		// after: add signal from grpc server
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error(), nil)
	}

	// after: add Shutdown from grpc server
}
