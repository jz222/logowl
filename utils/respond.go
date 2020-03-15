package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/models"
)

// RespondWithJSON returns a response in JSON to the client.
func RespondWithJSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// RespondWithError returns an error in JSON to the client.
func RespondWithError(c *gin.Context, code int, message string) {
	err := models.ErrorResponse{
		Ok:      false,
		Code:    code,
		Message: message,
	}

	c.JSON(err.Code, err)
}
