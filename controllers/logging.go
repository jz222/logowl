package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services/logging"
	"github.com/jz222/loggy/utils"
)

type loggingControllers struct{}

// Logging contains all controllers related to logging.
var Logging loggingControllers

func (l *loggingControllers) RegisterError(c *gin.Context) {
	errorEvent := models.Error{
		Badges:    map[string]string{},
		Host:      c.Request.Host,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Count:     1,
		Timestamp: time.Now().Unix(),
	}

	err := json.NewDecoder(c.Request.Body).Decode(&errorEvent)
	if err != nil {
		log.Println(err.Error())
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	go logging.SaveError(errorEvent)

	utils.RespondWithSuccess(c)
}
