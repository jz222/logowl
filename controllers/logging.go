package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services/event"
	"github.com/jz222/loggy/utils"
)

type eventControllers struct{}

var Logging eventControllers

func (l *eventControllers) RegisterError(c *gin.Context) {
	errorEvent := models.Error{
		Host:      c.Request.Host,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Count:     1,
	}

	err := json.NewDecoder(c.Request.Body).Decode(&errorEvent)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	event.SaveError(errorEvent)

	utils.RespondWithSuccess(c)
}

// func (l *eventControllers) LoadAllErrors(c *gin.Context) {
// 	persistedErrors := event.GetErrors()

// 	utils.RespondWithJSON(c, persistedErrors)
// }
