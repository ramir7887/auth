package mongorepository

import (
	"context"
	"gitlab.com/g6834/team28/auth/internal/entity"
	"gitlab.com/g6834/team28/auth/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository struct {
	client *mongodb.Mongo
}

func New(client *mongodb.Mongo) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

func (r *UserRepository) UserByName(name string) (entity.User, error) {
	u := User{}
	collection := r.client.Database(r.client.DbName).Collection("users")
	ctx := context.Background()

	if err := collection.FindOne(ctx, bson.D{primitive.E{Key: "name", Value: name}}).Decode(&u); err != nil {
		return entity.User{}, err
	}

	user := entity.User{
		ID:       u.ID,
		Name:     u.Name,
		Password: []byte(u.Password),
	}
	return user, nil
}

func (r *UserRepository) Create(user entity.User) error {
	collection := r.client.Database(r.client.DbName).Collection("users")
	ctx := context.Background()
	_, err := collection.InsertOne(ctx, bson.D{primitive.E{Key: "name", Value: user.Name}, primitive.E{Key: "password", Value: user.Password}})
	if err != nil {
		return err
	}
	return nil
}
