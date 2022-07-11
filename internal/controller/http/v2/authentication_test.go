//go:build integration

package v2

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gitlab.com/g6834/team28/auth/config"
	"gitlab.com/g6834/team28/auth/internal/entity"
	"gitlab.com/g6834/team28/auth/internal/repository/mongorepository"
	"gitlab.com/g6834/team28/auth/internal/usecase"
	"gitlab.com/g6834/team28/auth/pkg/httpserver"
	"gitlab.com/g6834/team28/auth/pkg/jwt"
	"gitlab.com/g6834/team28/auth/pkg/logger"
	"gitlab.com/g6834/team28/auth/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	configPath = "../../../../config/config.yml"
	user       = entity.User{Name: "user", Password: []byte("secret")}
	cfg        *config.Config
	db         *mongodb.Mongo
	repo       *mongorepository.UserRepository
	baseUrl    string
)

func clearRepository(dbName string, db *mongodb.Mongo) {
	collection := db.Database(dbName).Collection("users")
	_, err := collection.DeleteMany(context.Background(), bson.D{primitive.E{Key: "name", Value: user.Name}})
	if err != nil {
		logger.New("warn").Fatalf("error clearRepository: %s", err.Error())
	}
}

func TestMain(m *testing.M) {
	var err error
	l := logger.New("warn")
	cfg, err = config.NewConfig(configPath)
	if err != nil {
		l.Fatalf("error NewConfig: %s", err.Error())
	}
	baseUrl = "http://localhost:" + cfg.HTTP.Port
	db = mongodb.New(cfg.Mongo.Dsn, cfg.Mongo.DbName, l)
	repo = mongorepository.New(db)
	authenticationUseCase := usecase.NewAuthenticationUseCase(repo)
	userCreateUsecase := usecase.NewUserCreateUseCase(repo)
	handler := mux.NewRouter()
	NewRouter(handler, cfg.HTTP, l, authenticationUseCase, userCreateUsecase)

	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	defer httpServer.Shutdown()

	os.Exit(m.Run())
}

func TestServerHealtzHandler(t *testing.T) {
	c := http.Client{}

	r, err := c.Get(fmt.Sprintf("%s/healthz", baseUrl))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
}

func TestServerLoginHandler(t *testing.T) {
	c := http.Client{}

	testCases := []struct {
		name         string
		data         []byte
		expectedCode int
	}{
		{
			name:         "authenticated 200",
			data:         []byte(`{"login": "test123", "password": "qwerty"}`),
			expectedCode: http.StatusOK,
		},
		{
			name:         "not authenticated 403 incorrect login or password",
			data:         []byte(`{"login": "test123", "password": "weqweqweqw"}`),
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "not authenticated 500 incorrect body",
			data:         []byte(`{"login": "test123", "password": "weqweqweqw`),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBuffer(tc.data)
			r, err := c.Post(fmt.Sprintf("%s/login", baseUrl), "application/json", buf)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, r.StatusCode)
		})
	}
}

func TestServerLogoutHandler(t *testing.T) {
	c := http.Client{}

	r, err := c.Post(fmt.Sprintf("%s/logout", baseUrl), "", nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.NotEmpty(t, r.Header.Get("Set-Cookie"))
}

func TestServerInfoHandler(t *testing.T) {
	c := http.Client{}

	testCases := []struct {
		name         string
		login        string
		password     string
		setToken     bool
		expectedCode int
	}{
		{
			name:         "info 200",
			login:        "test123",
			password:     "qwerty",
			setToken:     true,
			expectedCode: http.StatusOK,
		},
		{
			name:         "info user not found",
			login:        "test13435523",
			password:     "qwerty",
			setToken:     true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "info 403 forbidden (middleware)",
			login:        "test123",
			password:     "qwerty",
			setToken:     false,
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				accessToken string
				body        = make([]byte, 512)
				err         error
			)

			accessToken, err = jwt.GenerateJwt("1", tc.login, 1*time.Minute)
			if err != nil {
				t.Fatalf("error GenerateJwt: %s", err.Error())
			}
			buf := bytes.NewBuffer(body)

			r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/validate", baseUrl), buf)
			if err != nil {
				t.Fatalf("error NewRequest: %s", err.Error())
			}
			if tc.setToken {
				r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
			}
			r.Header.Set("Content-Type", "application/json")

			res, err := c.Do(r)
			if err != nil {
				t.Fatalf("error Do: %s", err.Error())
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, res.StatusCode)
		})
	}

}
