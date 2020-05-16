package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jz222/loggy/internal/keys"
	"github.com/jz222/loggy/internal/models"
	"github.com/jz222/loggy/internal/store"
	"github.com/jz222/loggy/internal/utils"
)

type InterfaceAuth interface {
	CreateJWT(string) (string, int64, error)
	ResetPassword(user models.User) (string, error)
}

type Auth struct {
	Store   store.InterfaceStore
	Request InterfaceRequest
}

func (a *Auth) CreateJWT(id string) (string, int64, error) {
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

func (a *Auth) ResetPassword(user models.User) (string, error) {
	resetToken, err := utils.GenerateRandomString(50)
	if err != nil {
		return "", errors.New("an error occured while creating a password reset token")
	}

	passwordResetToken := models.PasswordResetToken{
		Email:     user.Email,
		Token:     resetToken,
		Used:      false,
		ExpiresAt: time.Now().Unix() + 60,
		CreatedAt: time.Now(),
	}

	_, err = a.Store.PasswordResetTokens().InsertOne(passwordResetToken)
	if err != nil {
		return "", err
	}

	data := map[string]interface{}{
		"FirstName": user.FirstName,
		"URL":       fmt.Sprintf("%s/auth/newpassword?token=%s", keys.GetKeys().CLIENT_URL, resetToken),
	}

	err = a.Request.SendEmail(user.Email, "resetPassword", data)
	if err != nil {
		return "", err
	}

	return passwordResetToken.Token, nil
}

func GetAuthService(store store.InterfaceStore) Auth {
	return Auth{store, &Request{}}
}
