package mongorepository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gitlab.com/g6834/team28/auth/config"
	"gitlab.com/g6834/team28/auth/internal/entity"
	"gitlab.com/g6834/team28/auth/pkg/logger"
	"gitlab.com/g6834/team28/auth/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"testing"
)

var (
	configPath = "../../../config/config.yml"
	cfg        *config.Config
	db         *mongodb.Mongo
	repo       *UserRepository
	user       = entity.User{Name: "user", Password: []byte("secret")}
)

func clearRepository(dbName string, repo *UserRepository) {
	collection := repo.client.Database(dbName).Collection("users")
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
	repo = New(db)
	clearRepository(cfg.Mongo.DbName, repo)
	defer clearRepository(cfg.Mongo.DbName, repo)
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	repo := New(db)

	assert.NotNil(t, repo)
}

func TestUserRepository_Create(t *testing.T) {
	repo := New(db)
	clearRepository(cfg.Mongo.DbName, repo)
	defer clearRepository(cfg.Mongo.DbName, repo)

	err := repo.Create(user)

	assert.NoError(t, err)

	u, err := repo.UserByName(user.Name)

	assert.NoError(t, err)
	assert.Equal(t, user.Name, u.Name)
	assert.Equal(t, user.Password, u.Password)
}

func TestUserRepository_Create_Error(t *testing.T) {
	repo := New(db)
	clearRepository(cfg.Mongo.DbName, repo)

	err := repo.Create(user)

	assert.NoError(t, err)

	err = repo.Create(user)

	assert.Error(t, err)
}

func TestUserRepository_UserByName(t *testing.T) {
	repo := New(db)
	clearRepository(cfg.Mongo.DbName, repo)

	err := repo.Create(user)

	assert.NoError(t, err)

	u, err := repo.UserByName(user.Name)

	assert.NoError(t, err)
	assert.Equal(t, user.Name, u.Name)
	assert.Equal(t, user.Password, u.Password)
}

func TestUserRepository_UserByName_NotFound(t *testing.T) {
	repo := New(db)

	_, err := repo.UserByName("testtest")

	assert.Error(t, err)
}
