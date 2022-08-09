package inmemory

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/g6834/team28/auth/internal/entity"
	"testing"
)

var (
	user = entity.User{Name: "user", Password: []byte("secret")}
)

func users() []entity.User {
	return []entity.User{
		{Name: "testuser2", Password: []byte("secret2")},
		{Name: "testuser3", Password: []byte("secret3")},
		{Name: "testuser4", Password: []byte("secret4")},
	}
}

func TestNew(t *testing.T) {
	repo := New(users())

	assert.NotNil(t, repo)
	assert.Equal(t, users(), repo.users)
}

func TestUserRepository_Create(t *testing.T) {
	repo := New(users())

	err := repo.Create(user)

	assert.NoError(t, err)

	u, err := repo.UserByName(user.Name)

	assert.NoError(t, err)
	assert.Equal(t, user, u)
}

func TestUserRepository_Create_Error(t *testing.T) {
	repo := New(users())

	err := repo.Create(user)

	assert.NoError(t, err)

	err = repo.Create(user)

	assert.Error(t, err)
}

func TestUserRepository_UserByName(t *testing.T) {
	repo := New(users())

	err := repo.Create(user)

	assert.NoError(t, err)

	u, err := repo.UserByName(user.Name)

	assert.NoError(t, err)
	assert.Equal(t, user, u)
}

func TestUserRepository_UserByName_NotFound(t *testing.T) {
	repo := New(users())

	_, err := repo.UserByName("testtest")

	assert.Error(t, err)
}
