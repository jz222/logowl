package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services/logs"
	"github.com/jz222/loggy/utils"
)

type loggingControllers struct{}

var Logging loggingControllers

func (l *loggingControllers) RegisterError(c *gin.Context) {
	var errorLog models.Error

	err := json.NewDecoder(c.Request.Body).Decode(&errorLog)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	logs.SaveError(errorLog)
}
