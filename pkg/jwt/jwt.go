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

func GenerateAllJwt(username string) (string, string, error) {
	accessToken, err := GenerateAccessJwt(username)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := GenerateRefreshJwt(username)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func GenerateAccessJwt(username string) (string, error) {
	return GenerateJwt(username, 1*time.Minute)
}

func GenerateRefreshJwt(username string) (string, error) {
	return GenerateJwt(username, 1*time.Hour)
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
