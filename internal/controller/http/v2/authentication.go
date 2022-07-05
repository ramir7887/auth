package v2

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/g6834/team28/auth/internal/controller/http/responder"
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

func newAuthenticationRoutes(r *mux.Router, l logger.Interface, uc usecase.Authentication) {
	ctx := context.Background()
	ar := &authenticationRoutes{
		uc: uc,
		logger: l.WithFields(logger.Fields{
			"package": "v2",
		}),
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

// login godoc
// @tags auth
// @Summary authentication user
// @Description authentication user by email and password
// @Accept json
// @Produce json
// @Param loginData body requestLogin true "Login Data"
// @Success 200 {object} responseLogin
// @response default {object} responseLogin
// @Header 200 {string} SetCookie "set accessToken and refreshToken"
// @Failure 500 {object} responseError
// @Failure 403 {object} responseError
// @Router /login [post]
func (a *authenticationRoutes) login(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := a.logger.WithFields(logger.Fields{
			"method": "authenticationRoutes.login",
		})
		l.Info("Start login handler v2")
		defer l.Info("End login handler v2")

		var req requestLogin

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			err = responder.JsonRespond(w, http.StatusInternalServerError, responseError{Error: http.StatusText(http.StatusInternalServerError)})
			if err != nil {
				l.WithFields(logger.Fields{
					"error": err.Error(),
				}).Error("Error respond")
			}
			return
		}

		token, refreshToken, err := a.uc.Login(ctx, req.Login, req.Password)
		if err != nil {
			err = responder.JsonRespond(w, http.StatusForbidden, responseError{Error: http.StatusText(http.StatusForbidden)})
			if err != nil {
				l.WithFields(logger.Fields{
					"error": err.Error(),
				}).Error("Error respond")
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
		err = responder.JsonRespond(w, http.StatusOK, res)
		if err != nil {
			l.WithFields(logger.Fields{
				"error": err.Error(),
			}).Error("Error respond")
		}

		// Пока не понятно надо или нет
		//vars := mux.Vars(r)
		//if redirectUri, ok := vars["redirect_uri"]; ok {
		//	a.logger.Info("redirect to", redirectUri)
		//	http.Redirect(w, r, redirectUri, http.StatusPermanentRedirect)
		//	return
		//}
	}
}

// logout godoc
// @tags auth
// @Summary logout user
// @Description logout user
// @Success 200
// @response default
// @Header 200 {string} SetCookie "set empty cookie"
// @Router /logout [post]
func (a *authenticationRoutes) logout(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := a.logger.WithFields(logger.Fields{
			"method": "authenticationRoutes.logout",
		})
		l.Info("Start logout handler v2")
		defer l.Info("End logout handler v2")

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

// info godoc
// @tags auth
// @Summary info by user
// @Description info user by token
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} entity.User
// @Failure 403 {object} responseError
// @Failure 404 {object} responseError
// @Failure 500 {object} responseError
// @Router /validate [post]
func (a *authenticationRoutes) info(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := a.logger.WithFields(logger.Fields{
			"method": "authenticationRoutes.logout",
		})
		l.Info("Start info handler v2")
		defer l.Info("End info handler v2")
		user, ok := r.Context().Value(keyUserData).(*userData)
		if !ok {
			err := responder.JsonRespond(w, http.StatusInternalServerError, responseError{Error: http.StatusText(http.StatusInternalServerError)})
			if err != nil {
				l.WithFields(logger.Fields{
					"error": err.Error(),
				}).Error("Error respond")
			}
			return
		}

		u, err := a.uc.Info(ctx, user.name)
		if err != nil {
			err = responder.JsonRespond(w, http.StatusNotFound, responseError{Error: http.StatusText(http.StatusNotFound)})
			if err != nil {
				l.WithFields(logger.Fields{
					"error": err.Error(),
				}).Error("Error respond")
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(u); err != nil {
			err = responder.JsonRespond(w, http.StatusInternalServerError, responseError{Error: http.StatusText(http.StatusInternalServerError)})
			if err != nil {
				l.WithFields(logger.Fields{
					"error": err.Error(),
				}).Error("Error respond")
			}
			return
		}
	}
}
