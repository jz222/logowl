package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jz222/logowl/internal/models"
)

// RespondWithJSON returns a response in JSON to the client.
func RespondWithJSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// RespondWithSuccess returns a success message.
func RespondWithSuccess(c *gin.Context) {
	success := models.Response{
		Ok:   true,
		Code: http.StatusOK,
	}

	c.JSON(http.StatusOK, success)
}

// RespondWithError returns an error in JSON to the client.
func RespondWithError(c *gin.Context, code int, message string) {
	err := models.Response{
		Ok:      false,
		Code:    code,
		Message: message,
	}

	c.JSON(err.Code, err)
}
