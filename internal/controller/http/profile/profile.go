package profile

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/g6834/team28/auth/internal/controller/http/responder"
	"gitlab.com/g6834/team28/auth/pkg/logger"
	"net/http"
	"net/http/pprof"
)

type Profile struct {
	enabled  bool
	login    string
	password string
	logger   logger.Interface
}

func New(enabled bool, login, password string, logger logger.Interface) *Profile {
	return &Profile{
		enabled:  enabled,
		login:    login,
		password: password,
		logger:   logger,
	}
}

func (p *Profile) NewRouter(r *mux.Router) {
	sub := r.PathPrefix("/").Subrouter()
	sub.HandleFunc("/", pprof.Index).Name("pprof-*")
	sub.HandleFunc("/cmdline", pprof.Cmdline).Name("pprof-cmdline")
	sub.HandleFunc("/profile", pprof.Profile).Name("pprof-profile")
	sub.HandleFunc("/symbol", pprof.Symbol).Name("pprof-symbol")
	sub.HandleFunc("/trace", pprof.Trace).Name("pprof-trace")
	sub.Use(p.CheckEnabledMiddleware, p.AuthMiddleware)

	subEnable := r.PathPrefix("/").Subrouter()
	subEnable.HandleFunc("/enable", p.EnableHttp).Methods(http.MethodPost).Name("pprof-enable")
	subEnable.Use(p.AuthMiddleware)
}

func (p *Profile) EnableHttp(w http.ResponseWriter, r *http.Request) {
	type RequestData struct {
		Enable bool `json:"enable"`
	}

	var req RequestData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if err := responder.JsonRespond(w, http.StatusBadRequest, map[string]string{"error": http.StatusText(http.StatusBadRequest)}); err != nil {
			p.logger.Error("Error respond: ", err.Error())
		}
		return
	}
	p.enabled = req.Enable
	if err := responder.JsonRespond(w, http.StatusOK, req); err != nil {
		p.logger.Error("Error respond: ", err.Error())
	}
}

func (p *Profile) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, password, ok := r.BasicAuth()
		if !ok {
			if err := responder.JsonRespond(w, http.StatusUnauthorized, map[string]string{"error": http.StatusText(http.StatusUnauthorized)}); err != nil {
				p.logger.Error("Error respond: ", err.Error())
			}
			return
		}
		if login != p.login || password != p.password {
			if err := responder.JsonRespond(w, http.StatusUnauthorized, map[string]string{"error": http.StatusText(http.StatusUnauthorized)}); err != nil {
				p.logger.Error("Error respond: ", err.Error())
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (p *Profile) CheckEnabledMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !p.enabled {
			if err := responder.JsonRespond(w, http.StatusNotFound, map[string]string{"error": http.StatusText(http.StatusNotFound)}); err != nil {
				p.logger.Error("Error respond: ", err.Error())
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}
