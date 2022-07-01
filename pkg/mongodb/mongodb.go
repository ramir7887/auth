package mongodb

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gitlab.com/g6834/team28/auth/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type Mongo struct {
	DbName string
	logger logger.Interface
	mongo.Client
}

func New(uri string, dbName string, l logger.Interface) *Mongo {
	l = l.WithFields(logger.Fields{
		"package": "mongodb",
		"method":  "New",
	})
	newCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(newCtx, options.Client().ApplyURI(uri))
	if err != nil {
		l.WithFields(logger.Fields{
			"error": err.Error(),
		}).Fatal("Error mongo.Connect")
		return nil
	}

	m := &Mongo{
		DbName: dbName,
		logger: l,
		Client: *client,
	}
	if err := m.Ping(newCtx, readpref.Primary()); err != nil {
		l.WithFields(logger.Fields{
			"error": err.Error(),
		}).Fatal("Error mongo.Client.Ping")
		return nil
	}
	l.Info("Connect to mongodb")

	return m
}

func (m *Mongo) Migrate(path, uri, mode string) error {
	mg, err := migrate.New(path, uri)
	if err != nil {
		return err
	}
	switch mode {
	case "up":
		err = mg.Up()
	case "down":
		err = mg.Down()
	default:
		return errors.New("incorrect mode for migration")
	}
	if err != nil {
		return err
	}
	return nil
}
