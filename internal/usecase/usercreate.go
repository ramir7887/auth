package usecase

import (
	"context"
	"gitlab.com/g6834/team28/auth/internal/entity"
	"gitlab.com/g6834/team28/auth/internal/repository"
)

type UserCreateUseCase struct {
	repo repository.UserRepository
}

func NewUserCreateUseCase(r repository.UserRepository) *UserCreateUseCase {
	return &UserCreateUseCase{
		repo: r,
	}
}

func (uc *UserCreateUseCase) Create(ctx context.Context, user entity.User) error {
	err := uc.repo.Create(user)
	if err != nil {
		return err
	}
	return nil
}
