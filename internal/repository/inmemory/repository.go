package inmemory

import (
	"errors"
	"gitlab.com/g6834/team28/auth/internal/entity"
)

type UserRepository struct {
	users []entity.User
}

func New(users []entity.User) *UserRepository {
	return &UserRepository{
		users: users,
	}
}

func (r *UserRepository) UserByName(name string) (entity.User, error) {
	for _, user := range r.users {
		if user.Name == name {
			return user, nil
		}
	}

	return entity.User{}, errors.New("user not found")
}

func (r *UserRepository) Create(user entity.User) error {
	if _, err := r.UserByName(user.Name); err == nil {
		return errors.New("name already use")
	}
	r.users = append(r.users, user)
	return nil
}
