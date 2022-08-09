//go:build integration

package usecase

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gitlab.com/g6834/team28/auth/config"
	"gitlab.com/g6834/team28/auth/internal/entity"
	"gitlab.com/g6834/team28/auth/internal/repository"
	"gitlab.com/g6834/team28/auth/internal/repository/mongorepository"
	"gitlab.com/g6834/team28/auth/pkg/logger"
	"gitlab.com/g6834/team28/auth/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"testing"
)

var (
	configPath = "../../config/config.yml"
	userPwd    = "qwerty"
	user       = entity.User{Name: "user", Password: []byte("$2a$14$ke.euFhogr5WCqXSagVuceiZc2o6w/Z6hoL9KIVqIPZdEZL4TRF9K")}
	cfg        *config.Config
	repo       repository.UserRepository
	db         *mongodb.Mongo
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
	db = mongodb.New(cfg.Mongo.Dsn, cfg.Mongo.DbName, logger.New("error"))
	repo = mongorepository.New(db)
	clearRepository(cfg.Mongo.DbName, db)

	os.Exit(m.Run())
}

func TestNewUserCreateUseCase(t *testing.T) {
	uc := NewUserCreateUseCase(repo)

	assert.NotNil(t, uc)
}

func TestUserCreateUseCase_Create(t *testing.T) {
	clearRepository(cfg.Mongo.DbName, db)
	defer clearRepository(cfg.Mongo.DbName, db)
	uc := NewUserCreateUseCase(repo)

	err := uc.Create(context.Background(), user)

	assert.NoError(t, err)
}

func TestUserCreateUseCase_Create_Error(t *testing.T) {
	clearRepository(cfg.Mongo.DbName, db)
	defer clearRepository(cfg.Mongo.DbName, db)
	uc := NewUserCreateUseCase(repo)

	err := uc.Create(context.Background(), user)

	assert.NoError(t, err)

	err = uc.Create(context.Background(), user)

	assert.Error(t, err)
}
