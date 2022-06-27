package usecase

import (
	"context"
	"errors"
	"gitlab.com/g6834/team28/auth/internal/entity"
	"gitlab.com/g6834/team28/auth/internal/repository"
	"gitlab.com/g6834/team28/auth/pkg/jwt"
	"gitlab.com/g6834/team28/auth/pkg/password"
	"time"
)

var (
	InvalidNameOrPassword error = errors.New("invalid username or password")
)

type AuthenticationUseCase struct {
	repo repository.UserRepository
}

func New(r repository.UserRepository) *AuthenticationUseCase {
	return &AuthenticationUseCase{
		repo: r,
	}
}

func (uc *AuthenticationUseCase) Login(ctx context.Context, name, pass string) (string, string, error) {
	u, err := uc.repo.UserByName(name)
	if err != nil {
		return "", "", InvalidNameOrPassword
	}
	if !password.ComparePassword(pass, u.Password) {
		return "", "", InvalidNameOrPassword
	}

	//JWT create Access and Refresh
	accessToken, err := jwt.GenerateJwt(u.Name, 1*time.Minute)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwt.GenerateJwt(u.Name, 1*time.Hour)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, err
}

func (uc *AuthenticationUseCase) Logout(ctx context.Context, name string) error {
	if _, err := uc.repo.UserByName(name); err != nil {
		return err
	}
	return nil
}

func (uc *AuthenticationUseCase) Info(ctx context.Context, name string) (entity.User, error) {
	u, err := uc.repo.UserByName(name)
	if err != nil {
		return entity.User{}, err
	}
	return u, nil
}

func (uc *AuthenticationUseCase) Validate(ctx context.Context, accessToken, refreshToken string) (string, string, error) {
	var newAccessToken, newRefreshToken string

	claimAccess, err := jwt.Parse(accessToken)
	if err != nil {
		return "", "", err
	}

	if !jwt.Expired(claimAccess) {
		return accessToken, refreshToken, nil
	}

	claimRefresh, err := jwt.Parse(refreshToken)
	if err != nil {
		return "", "", err
	}

	if !jwt.Expired(claimRefresh) {
		newAccessToken, newRefreshToken, err = jwt.GenerateAllJwt(claimAccess.Username)
		if err != nil {
			return "", "", err
		}
	}
	return newAccessToken, newRefreshToken, nil
}
