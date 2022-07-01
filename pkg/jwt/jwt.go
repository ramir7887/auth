package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	jwtKey            = []byte("qwerty")
	ParseError        = errors.New("couldn't parse claims")
	TokenExpiredError = errors.New("token expired")
)

type Claim struct {
	jwt.StandardClaims
	ID       string `json:"id"`
	Username string `json:"username"`
}

func GenerateAllJwt(id, username string) (string, string, error) {
	accessToken, err := GenerateAccessJwt(id, username)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := GenerateRefreshJwt(id, username)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func GenerateAccessJwt(id, username string) (string, error) {
	return GenerateJwt(id, username, 1*time.Minute)
}

func GenerateRefreshJwt(id, username string) (string, error) {
	return GenerateJwt(id, username, 1*time.Hour)
}

func GenerateJwt(id, username string, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	claims := &Claim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		ID:       id,
		Username: username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func Parse(signedToken string) (*Claim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&Claim{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claim)
	if !ok {
		return nil, ParseError
	}
	return claims, nil
}

func ValidateToken(signedToken string) (*Claim, error) {
	claims, err := Parse(signedToken)
	if err != nil {
		return nil, err
	}
	if Expired(claims) {
		return nil, TokenExpiredError
	}
	return claims, nil
}

func Expired(claim *Claim) bool {
	return claim.ExpiresAt < time.Now().Local().Unix()
}
