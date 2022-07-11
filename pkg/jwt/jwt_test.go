//go:build !integration

package jwt

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	id       string = "sdfg8939jfsk"
	username string = "tester"
)

func TestGenerateJwt(t *testing.T) {
	token, err := GenerateJwt(id, username, 5*time.Minute)

	assert.NoError(t, err, "not error")
	assert.NotEqual(t, "", token, "must be not empty string")
}

func TestGenerateAccessJwt(t *testing.T) {
	token, err := GenerateAccessJwt(id, username)

	assert.NoError(t, err, "not error")
	assert.NotEqual(t, "", token, "must be not empty string")
}

func TestGenerateRefreshJwt(t *testing.T) {
	token, err := GenerateRefreshJwt(id, username)

	assert.NoError(t, err, "not error")
	assert.NotEqual(t, "", token, "must be not empty string")
}

func TestGenerateAllJwt(t *testing.T) {
	token, refreshToken, err := GenerateAllJwt(id, username)

	assert.NoError(t, err, "not error")
	assert.NotEqual(t, "", token, "must be not empty string")
	assert.NotEqual(t, "", refreshToken, "must be not empty string")
}

func TestParse(t *testing.T) {
	token, err := GenerateJwt(id, username, 5*time.Minute)
	if err != nil {
		t.Fatalf("error GenerateJwt: %s", err.Error())
	}
	claim, err := Parse(token)

	assert.NoError(t, err, "not error")
	assert.Equal(t, id, claim.ID)
	assert.Equal(t, username, claim.Username)
}

func TestExpired_Expired(t *testing.T) {
	token, err := GenerateJwt(id, username, 5*time.Minute)
	if err != nil {
		t.Fatalf("error GenerateJwt: %s", err.Error())
	}
	claim, err := Parse(token)
	if err != nil {
		t.Fatalf("error Parse: %s", err.Error())
	}

	expired := Expired(claim)

	assert.Equal(t, false, expired)
}

func TestExpired_NotExpired(t *testing.T) {
	token, err := GenerateJwt(id, username, 1*time.Second)
	if err != nil {
		t.Fatalf("error GenerateJwt: %s", err.Error())
	}
	claim, err := Parse(token)
	if err != nil {
		t.Fatalf("error Parse: %s", err.Error())
	}
	time.Sleep(2 * time.Second)

	expired := Expired(claim)

	assert.Equal(t, true, expired)
}

func TestValidateToken_Valid(t *testing.T) {
	token, err := GenerateJwt(id, username, 5*time.Minute)
	if err != nil {
		t.Fatalf("error GenerateJwt: %s", err.Error())
	}

	claim, err := ValidateToken(token)
	assert.NoError(t, err, "not error")
	assert.Equal(t, id, claim.ID)
	assert.Equal(t, username, claim.Username)
}

func TestValidateToken_NotValid(t *testing.T) {
	token, err := GenerateJwt(id, username, 1*time.Second)
	if err != nil {
		t.Fatalf("error GenerateJwt: %s", err.Error())
	}
	time.Sleep(2 * time.Second)

	_, err = ValidateToken(token)
	assert.Error(t, err, "error TokenExpiredError")
}
