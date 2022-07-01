package repository

import "gitlab.com/g6834/team28/auth/internal/entity"

type (
	UserRepository interface {
		UserByName(string) (entity.User, error)
		Create(user entity.User) error
	}
)
