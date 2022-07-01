package v2

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/g6834/team28/auth/internal/usecase"
	"gitlab.com/g6834/team28/auth/pkg/logger"
	"net/http"
)

type userData struct {
	name  string
	token string
}

type requestLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type responseLogin struct {
	Name         string `json:"username"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type responseError struct {
	Error string `json:"error"`
}

// NewRouter -.
func NewRouter(router *mux.Router, l logger.Interface, a usecase.Authentication, u usecase.UserCreate) {
	// Options
	router.Use(LoggingMiddleware(l))
	router.Use(handlers.RecoveryHandler())

	// K8s probe
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }).Methods(http.MethodGet)

	// Prometheus metrics
	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	// Routers
	userRouter := router.PathPrefix("/user").Subrouter()
	newUserCreateRoutes(userRouter, l, u)

	authRouter := mux.NewRouter()
	newAuthenticationRoutes(authRouter, l, a)
	router.PathPrefix("/").Handler(authRouter).Name("auth router v2")
}
