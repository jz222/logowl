package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/internal/models"
	"github.com/jz222/loggy/internal/services"
	"github.com/jz222/loggy/internal/store"
	"github.com/jz222/loggy/internal/utils"
)

type loggingControllers struct {
	LoggingService services.InterfaceLogging
}

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

	go l.LoggingService.SaveError(errorEvent)

	utils.RespondWithSuccess(c)
}

func GetLoggingController(store store.InterfaceStore) loggingControllers {
	loggingService := services.GetLoggingService(store)

	return loggingControllers{
		LoggingService: &loggingService,
	}
}
