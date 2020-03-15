package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/keys"
	"github.com/jz222/loggy/utils"
)

func VerifyUserJwt(c *gin.Context) {
	authenticationHeader := c.GetHeader("Authorization")
	splitHeader := strings.Split(authenticationHeader, " ")

	if len(splitHeader) != 2 {
		utils.RespondWithError(c, http.StatusBadRequest, "authorization header malformed")
		c.Abort()
		return
	}

	token, err := jwt.Parse(splitHeader[1], func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("")
		}

		return []byte(keys.GetKeys().SECRET), nil
	})
	if err != nil || !token.Valid {
		utils.RespondWithError(c, http.StatusUnauthorized, "incorrect JWT")
		c.Abort()
	}

	c.Next()
}
