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
	Username string `json:"username"`
}

func GenerateJwt(username string, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	claims := &Claim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		Username: username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(signedToken string) (*Claim, error) {
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
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, TokenExpiredError
	}
	return claims, nil
}
