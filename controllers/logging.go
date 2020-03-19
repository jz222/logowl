package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services/events"
	"github.com/jz222/loggy/utils"
)

type eventControllers struct{}

var Logging eventControllers

func (l *eventControllers) RegisterError(c *gin.Context) {
	var errorLog models.Error

	err := json.NewDecoder(c.Request.Body).Decode(&errorLog)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	events.SaveError(errorLog)

	utils.RespondWithSuccess(c)
}

func (l *eventControllers) LoadAllErrors(c *gin.Context) {
	persistedErrors := events.GetErrors()

	utils.RespondWithJSON(c, persistedErrors)
}
