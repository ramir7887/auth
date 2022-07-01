package v2

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/g6834/team28/auth/internal/controller/http/responder"
	"gitlab.com/g6834/team28/auth/internal/entity"
	"gitlab.com/g6834/team28/auth/internal/usecase"
	"gitlab.com/g6834/team28/auth/pkg/logger"
	"gitlab.com/g6834/team28/auth/pkg/password"
	"net/http"
)

type userCreateRoutes struct {
	uc     usecase.UserCreate
	logger logger.Interface
}

func newUserCreateRoutes(r *mux.Router, l logger.Interface, uc usecase.UserCreate) {
	ctx := context.Background()
	ur := &userCreateRoutes{
		uc: uc,
		logger: l.WithFields(logger.Fields{
			"package": "v2",
		}),
	}

	r.HandleFunc("", ur.create(ctx)).Methods(http.MethodPost).Name("user create")
}

func (ur *userCreateRoutes) create(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := ur.logger.WithFields(logger.Fields{
			"method": "userCreateRoutes.create",
		})
		l.Info("Start create handler v2")
		defer l.Info("End create handler v2")

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
		user := entity.User{
			Name:     req.Login,
			Password: []byte(req.Password),
		}
		user.Password = password.HashPassword(user.Password)
		err := ur.uc.Create(ctx, user)
		if err != nil {
			err = responder.JsonRespond(w, http.StatusConflict, responseError{Error: http.StatusText(http.StatusConflict)})
			if err != nil {
				l.WithFields(logger.Fields{
					"error": err.Error(),
				}).Error("Error respond")
			}
			return
		}

		err = responder.JsonRespond(w, http.StatusCreated, map[string]string{"message": http.StatusText(http.StatusCreated)})
		if err != nil {
			l.WithFields(logger.Fields{
				"error": err.Error(),
			}).Error("Error respond")
		}
	}
}
