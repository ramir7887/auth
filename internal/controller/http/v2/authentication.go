package v2

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

func newAuthenticationRoutes(r *mux.Router, l logger.Interface, uc usecase.Authentication) {
	ctx := context.Background()
	ar := &authenticationRoutes{
		uc:     uc,
		logger: l,
	}

	r.HandleFunc("/login", ar.login(ctx)).Queries("redirect_uri", "{redirect_uri}").Methods(http.MethodPost).Name("login with redirect v2")
	r.HandleFunc("/login", ar.login(ctx)).Methods(http.MethodPost).Name("login v2")
	r.HandleFunc("/logout", ar.logout(ctx)).Queries("redirect_uri", "{redirect_uri}").Methods(http.MethodPost).Name("logout with redirect v2")
	r.HandleFunc("/logout", ar.logout(ctx)).Methods(http.MethodPost).Name("logout v2")

	// SubRoutes for middleware TokenMiddleware
	sub := r.PathPrefix("/").Subrouter()
	sub.HandleFunc("/validate", ar.info(ctx)).Methods(http.MethodPost).Name("information v2")
	sub.Use(TokenMiddleware(ctx, l))
}

func (a *authenticationRoutes) login(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Start login handler v2")
		defer a.logger.Info("End login handler v2")

		var req requestLogin

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			err = jsonRespond(w, http.StatusInternalServerError, responseError{Error: http.StatusText(http.StatusInternalServerError)})
			if err != nil {
				a.logger.Error("Error respond: ", err.Error())
			}
			return
		}

		token, refreshToken, err := a.uc.Login(ctx, req.Login, req.Password)
		if err != nil {
			err = jsonRespond(w, http.StatusForbidden, responseError{Error: http.StatusText(http.StatusForbidden)})
			if err != nil {
				a.logger.Error("Error respond: ", err.Error())
			}
			return
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
		res := responseLogin{
			Name:         req.Login,
			AccessToken:  token,
			RefreshToken: refreshToken,
		}

		http.SetCookie(w, &cookieAccess)
		http.SetCookie(w, &cookieRefresh)
		err = jsonRespond(w, http.StatusOK, res)
		if err != nil {
			a.logger.Error("Error respond: ", err.Error())
		}
		return

		// Пока не понятно надо или нет
		//vars := mux.Vars(r)
		//if redirectUri, ok := vars["redirect_uri"]; ok {
		//	a.logger.Info("redirect to", redirectUri)
		//	http.Redirect(w, r, redirectUri, http.StatusPermanentRedirect)
		//	return
		//}
	}
}

func (a *authenticationRoutes) logout(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Start logout handler v2")
		defer a.logger.Info("End logout handler v2")

		cookieAccess := http.Cookie{
			Name:     "accessToken",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		}
		cookieRefresh := http.Cookie{
			Name:     "refreshToken",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		}
		http.SetCookie(w, &cookieAccess)
		http.SetCookie(w, &cookieRefresh)

		// Пока не понятно надо или нет
		//vars := mux.Vars(r)
		//if redirectUri, ok := vars["redirect_uri"]; ok {
		//	a.logger.Info("redirect to", redirectUri)
		//	http.Redirect(w, r, redirectUri, http.StatusPermanentRedirect)
		//	return
		//}
	}
}

func (a *authenticationRoutes) info(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Start info handler v2")
		defer a.logger.Info("End info handler v2")
		user, ok := r.Context().Value(keyUserData).(*userData)
		if !ok {
			err := jsonRespond(w, http.StatusInternalServerError, responseError{Error: http.StatusText(http.StatusInternalServerError)})
			if err != nil {
				a.logger.Error("Error respond: ", err.Error())
			}
			return
		}

		u, err := a.uc.Info(ctx, user.name)
		if err != nil {
			err = jsonRespond(w, http.StatusNotFound, responseError{Error: http.StatusText(http.StatusNotFound)})
			if err != nil {
				a.logger.Error("Error respond: ", err.Error())
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(u); err != nil {
			err = jsonRespond(w, http.StatusInternalServerError, responseError{Error: http.StatusText(http.StatusInternalServerError)})
			if err != nil {
				a.logger.Error("Error respond: ", err.Error())
			}
			return
		}
	}
}
