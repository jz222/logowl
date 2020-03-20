package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/keys"
	"github.com/jz222/loggy/services/user"
	"github.com/jz222/loggy/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		utils.RespondWithError(c, http.StatusUnauthorized, "incorrect JWT")
		c.Abort()
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims["id"].(string))
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	userData, err := user.FindOne(bson.M{"_id": userID})
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set("user", userData)

	c.Next()
}
