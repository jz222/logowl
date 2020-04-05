package services

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jz222/loggy/keys"
)

type InterfaceAuth interface {
	CreateJWT(string) (string, int64, error)
}

type auth struct{}

func (a *auth) CreateJWT(id string) (string, int64, error) {
	timestamp := time.Now().Unix()
	expiresAt := timestamp + int64((time.Hour.Seconds() * 7))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"iat": timestamp,
		"exp": expiresAt,
	})

	signedToken, err := token.SignedString([]byte(keys.GetKeys().SECRET))
	if err != nil {
		return "", 0, err
	}

	return signedToken, expiresAt * 1000, nil
}

func GetAuthService() auth {
	return auth{}
}
