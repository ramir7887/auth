//go:build integration

package usecase

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAuthenticationUseCase(t *testing.T) {
	uc := NewAuthenticationUseCase(repo)

	assert.NotNil(t, uc)
}

func TestAuthenticationUseCase_Login(t *testing.T) {
	clearRepository(cfg.Mongo.DbName, db)
	defer clearRepository(cfg.Mongo.DbName, db)
	ucCreate := NewUserCreateUseCase(repo)
	err := ucCreate.Create(context.Background(), user)

	assert.NoError(t, err)

	uc := NewAuthenticationUseCase(repo)

	accessToken, refreshToken, err := uc.Login(context.Background(), user.Name, userPwd)

	assert.NoError(t, err)
	assert.NotEqual(t, "", accessToken)
	assert.NotEqual(t, "", refreshToken)
}

func TestAuthenticationUseCase_Login_Error(t *testing.T) {
	clearRepository(cfg.Mongo.DbName, db)
	defer clearRepository(cfg.Mongo.DbName, db)
	ucCreate := NewUserCreateUseCase(repo)
	err := ucCreate.Create(context.Background(), user)

	assert.NoError(t, err)

	uc := NewAuthenticationUseCase(repo)

	_, _, err = uc.Login(context.Background(), user.Name, "*******")

	assert.Error(t, err)
}

func TestAuthenticationUseCase_Logout(t *testing.T) {
	clearRepository(cfg.Mongo.DbName, db)
	defer clearRepository(cfg.Mongo.DbName, db)
	ucCreate := NewUserCreateUseCase(repo)
	err := ucCreate.Create(context.Background(), user)

	assert.NoError(t, err)

	uc := NewAuthenticationUseCase(repo)

	err = uc.Logout(context.Background(), user.Name)

	assert.NoError(t, err)
}

func TestAuthenticationUseCase_Info(t *testing.T) {
	clearRepository(cfg.Mongo.DbName, db)
	defer clearRepository(cfg.Mongo.DbName, db)
	ucCreate := NewUserCreateUseCase(repo)
	err := ucCreate.Create(context.Background(), user)

	assert.NoError(t, err)

	uc := NewAuthenticationUseCase(repo)

	u, err := uc.Info(context.Background(), user.Name)

	assert.NoError(t, err)
	assert.Equal(t, user.Name, u.Name)
}

func TestAuthenticationUseCase_Info_Error(t *testing.T) {
	clearRepository(cfg.Mongo.DbName, db)
	defer clearRepository(cfg.Mongo.DbName, db)
	ucCreate := NewUserCreateUseCase(repo)
	err := ucCreate.Create(context.Background(), user)

	assert.NoError(t, err)

	uc := NewAuthenticationUseCase(repo)

	_, err = uc.Info(context.Background(), "testtesttest")

	assert.Error(t, err)
}

func TestAuthenticationUseCase_Validate(t *testing.T) {
	clearRepository(cfg.Mongo.DbName, db)
	defer clearRepository(cfg.Mongo.DbName, db)
	ucCreate := NewUserCreateUseCase(repo)
	err := ucCreate.Create(context.Background(), user)

	assert.NoError(t, err)

	uc := NewAuthenticationUseCase(repo)
	accessToken, refreshToken, err := uc.Login(context.Background(), user.Name, userPwd)

	newAccessToken, newRefreshToken, err := uc.Validate(context.Background(), accessToken, refreshToken)

	assert.NoError(t, err)
	assert.Equal(t, accessToken, newAccessToken)
	assert.Equal(t, refreshToken, newRefreshToken)
}

func TestAuthenticationUseCase_Validate_Error(t *testing.T) {
	clearRepository(cfg.Mongo.DbName, db)
	defer clearRepository(cfg.Mongo.DbName, db)

	uc := NewAuthenticationUseCase(repo)
	_, _, err := uc.Validate(context.Background(), "dfdfdf", "dfdfdfdfde")

	assert.Error(t, err)
}
