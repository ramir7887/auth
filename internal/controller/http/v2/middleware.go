package v2

import (
	"context"
	"fmt"
	"gitlab.com/g6834/team28/auth/pkg/jwt"
	"gitlab.com/g6834/team28/auth/pkg/logger"
	"net/http"
	"strings"
	"time"
)

func TokenMiddleware(ctx context.Context, l logger.Interface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Info("Start TokenMiddleware")

			token := r.Header.Get("Authorization")
			if token == "" {
				cookie, err := r.Cookie("accessToken")
				if err != nil || cookie.Value == "" {
					jsonRespond(w, http.StatusForbidden, map[string]string{"error": http.StatusText(http.StatusForbidden)})
					return
				}
				token = cookie.Value
			}

			token = strings.Replace(token, "Bearer ", "", 1)
			claim, err := jwt.ValidateToken(token)
			if err != nil {
				jsonRespond(w, http.StatusInternalServerError, map[string]string{"error": http.StatusText(http.StatusInternalServerError)})
				return
			}
			user := &userData{
				name:  claim.Username,
				token: token,
			}
			r = r.WithContext(context.WithValue(r.Context(), keyUserData, user))
			next.ServeHTTP(w, r)
		})
	}
}

func LoggingMiddleware(l logger.Interface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			l.Info(fmt.Sprintf("Request: %s %s - %s - %s", r.Method, r.URL.String(), r.RemoteAddr, time.Now().Sub(start).String()))
		})
	}
}
