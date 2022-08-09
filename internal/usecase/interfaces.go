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
		Validate(context.Context, string, string) (string, string, error)
	}

	// UserCreate -.
	UserCreate interface {
		Create(context.Context, entity.User) error
	}
)
