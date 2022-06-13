package v1

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/g6834/team28/auth/internal/usecase"
	"gitlab.com/g6834/team28/auth/pkg/logger"
	"net/http"
	"time"
)

const keyUserData key = 1

type key uint

type authenticationRoutes struct {
	uc     usecase.Authentication
	logger logger.Interface
}

type userData struct {
	name         string
	token        string
	refreshToken string
}

func newAuthenticationRoutes(r *mux.Router, l logger.Interface, uc usecase.Authentication) {
	ctx := context.Background()
	ar := &authenticationRoutes{
		uc:     uc,
		logger: l,
	}

	r.HandleFunc("/login", ar.login(ctx)).Queries("redirect_uri", "{redirect_uri}").Methods(http.MethodPost).Name("login with redirect v1")
	r.HandleFunc("/login", ar.login(ctx)).Methods(http.MethodPost).Name("login v1")

	sub := r.PathPrefix("/").Subrouter()
	sub.HandleFunc("/logout", ar.logout(ctx)).Queries("redirect_uri", "{redirect_uri}").Methods(http.MethodPost).Name("logout with redirect v1")
	sub.HandleFunc("/logout", ar.logout(ctx)).Methods(http.MethodPost).Name("logout v1")
	sub.HandleFunc("/i", ar.info(ctx)).Methods(http.MethodGet).Name("information v1")
	sub.Use(TokenMiddleware(ctx, l))
}

func (a *authenticationRoutes) login(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Start login handler v1")
		defer a.logger.Info("End login handler v1")

		name, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, refreshToken, err := a.uc.Login(ctx, name, password)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
		}

		cookieAccess := http.Cookie{
			Name:    "accessToken",
			Value:   token,
			Path:    "/",
			Expires: time.Now().Add(1 * time.Minute),
		}
		cookieRefresh := http.Cookie{
			Name:    "refreshToken",
			Value:   refreshToken,
			Path:    "/",
			Expires: time.Now().Add(1 * time.Hour),
		}

		http.SetCookie(w, &cookieAccess)
		http.SetCookie(w, &cookieRefresh)

		vars := mux.Vars(r)
		if redirectUri, ok := vars["redirect_uri"]; ok {
			a.logger.Info("redirect to ", redirectUri)
			http.Redirect(w, r, redirectUri, http.StatusPermanentRedirect)
			return
		}
	}
}

func (a *authenticationRoutes) logout(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Start logout handler v1")
		defer a.logger.Info("End logout handler v1")

		user, ok := r.Context().Value(keyUserData).(*userData)
		if !ok {
			http.Error(w, "error user data", http.StatusInternalServerError)
			return
		}
		if err := a.uc.Logout(ctx, user.name); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		cookieAccess := http.Cookie{
			Name:  "accessToken",
			Value: "",
			Path:  "/",
		}
		cookieRefresh := http.Cookie{
			Name:  "refreshToken",
			Value: "",
			Path:  "/",
		}
		http.SetCookie(w, &cookieAccess)
		http.SetCookie(w, &cookieRefresh)

		vars := mux.Vars(r)
		if redirectUri, ok := vars["redirect_uri"]; ok {
			a.logger.Info("redirect to ", redirectUri)
			http.Redirect(w, r, redirectUri, http.StatusPermanentRedirect)
			return
		}
	}
}

func (a *authenticationRoutes) info(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Start info handler v1")
		defer a.logger.Info("End info handler v1")
		user, ok := r.Context().Value(keyUserData).(*userData)
		if !ok {
			http.Error(w, "error user data", http.StatusInternalServerError)
			return
		}

		u, err := a.uc.Info(ctx, user.name)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
