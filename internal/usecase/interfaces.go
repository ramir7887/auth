// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"gitlab.com/g6834/team28/auth/internal/entity"
)

type (
	// Authentication -.
	Authentication interface {
		Login(context.Context, string, string) (string, string, error)
		Logout(context.Context, string) error
		Info(context.Context, string) (entity.User, error)
	}

	// UserRepository -.
	UserRepository interface {
		UserByName(string) (entity.User, error)
	}
)
