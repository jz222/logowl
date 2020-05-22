package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/internal/keys"
	"github.com/jz222/loggy/internal/services"
	"github.com/jz222/loggy/internal/store"
	"github.com/jz222/loggy/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VerifyUserJwt checks if a JWT is present in the "Authorization" header and validates it.
func VerifyUserJwt(store store.InterfaceStore) func(*gin.Context) {
	return func(c *gin.Context) {
		// Parse JWT signature from cookie
		signature, err := c.Cookie("auth-signature")
		if err != nil {
			signature = ""
		}

		// Parse JWT from authorization header
		authorizationHeader := c.GetHeader("Authorization")
		splitHeader := strings.Split(authorizationHeader, " ")

		// Abort if no JWT is present or if the authorization is malformed
		if len(splitHeader) != 2 {
			utils.RespondWithError(c, http.StatusBadRequest, "authorization header malformed")
			c.Abort()
			return
		}

		// If a cookie is present, add the signature from the cookie
		// to the first part of the JWT from the authorization header.
		// Else use the complete JWT from the authorization header.
		JWT := ""

		if signature != "" {
			JWT = splitHeader[1] + "." + signature
		} else {
			JWT = splitHeader[1]
		}

		// Validate the JWT
		token, err := jwt.Parse(JWT, func(token *jwt.Token) (interface{}, error) {
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

		user := services.GetUserService(store)

		userData, err := user.FindOne(bson.M{"_id": userID})
		if err != nil {
			utils.RespondWithError(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set("user", userData)

		c.Next()
	}
}
