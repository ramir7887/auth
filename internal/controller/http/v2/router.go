package v2

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/http-swagger"
	"gitlab.com/g6834/team28/auth/config"
	_ "gitlab.com/g6834/team28/auth/docs"
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

type responseMsg struct {
	Message string `json:"message"`
}

// @title Auth API
// @version 2.0
// @description This is a sample auth server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email ramir7887@yandex.ru

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @query.collection.format multi

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @x-extension-openapi {"example": "value on a json format"}

// NewRouter -.
func NewRouter(router *mux.Router, cfg config.HTTP, l logger.Interface, a usecase.Authentication, u usecase.UserCreate) {
	// Options
	router.Use(LoggingMiddleware(l))
	router.Use(handlers.RecoveryHandler())

	//Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s:%s%s/swagger/doc.json", cfg.ServeAddress, cfg.Port, cfg.BasePath)),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	// K8s probe
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }).Methods(http.MethodGet)

	// Prometheus metrics
	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	// Routers
	userRouter := router.PathPrefix("/user").Subrouter()
	newUserCreateRoutes(userRouter, l, u)

	authRouter := mux.NewRouter()
	newAuthenticationRoutes(authRouter, l, a)
	router.PathPrefix("/").Handler(authRouter).Name("auth router v2")
}
