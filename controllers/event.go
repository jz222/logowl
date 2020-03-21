package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services/event"
	"github.com/jz222/loggy/services/service"
	"github.com/jz222/loggy/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type eventControllers struct{}

var Event eventControllers

func (e *eventControllers) GetErrors(c *gin.Context) {
	serviceID := c.Param("service")
	pointer := c.Param("pointer")

	parsedPage, err := strconv.ParseInt(pointer, 10, 64)
	if err != nil {
		parsedPage = 0
	}

	user, ok := c.Get("user")
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "could not read user data")
		return
	}

	organizationID := user.(models.User).OrganizationID

	parsedServiceID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "the provided service ID is invalid")
		return
	}

	requestedService, err := service.FindOne(bson.M{"_id": parsedServiceID, "organizationId": organizationID})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	persistedErrors, err := event.GetErrors(requestedService.Ticket, parsedPage)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, persistedErrors)
}
