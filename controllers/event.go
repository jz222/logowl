package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/services/event"
	"github.com/jz222/loggy/utils"
)

type eventControllers struct{}

var Event eventControllers

func (e *eventControllers) GetErrors(c *gin.Context) {
	pointer := c.Param("pointer")

	persistedErrors, err := event.GetErrors(pointer)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, persistedErrors)
}
