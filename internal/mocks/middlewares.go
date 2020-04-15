package mocks

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/internal/models"
)

func VerifyUserJWT(c *gin.Context) {
	userData := models.User{
		Email: "test@example.com",
	}

	c.Set("user", userData)

	c.Next()
}
