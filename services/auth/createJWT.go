package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jz222/loggy/keys"
)

func CreateJWT(id string) (string, error) {
	timestamp := time.Now().Unix()
	expiresAt := timestamp + 1000

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"iat": timestamp,
		"exp": expiresAt,
	})

	signedToken, err := token.SignedString([]byte(keys.GetKeys().SECRET))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
