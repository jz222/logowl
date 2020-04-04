package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services"
	"github.com/jz222/loggy/services/service"
	"github.com/jz222/loggy/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type eventControllers struct {
	EventService services.InterfaceEvent
}

func (e *eventControllers) GetError(c *gin.Context) {
	errorID := c.Param("id")
	serviceID := c.Param("service")

	parsedErrorID, err := primitive.ObjectIDFromHex(errorID)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "the provided error ID is invalid")
		return
	}

	parsedServiceID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "the provided service ID is invalid")
		return
	}

	userData, ok := c.Get("user")
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "failed to parse user data")
		return
	}

	persistedService, err := service.FindOne(bson.M{"_id": parsedServiceID, "organizationId": userData.(models.User).OrganizationID})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	errorEvent, err := e.EventService.GetError(bson.M{"_id": parsedErrorID, "ticket": persistedService.Ticket})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, errorEvent)
}

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

	persistedErrors, err := e.EventService.GetErrors(requestedService.Ticket, parsedPage)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, persistedErrors)
}

func (e *eventControllers) DeleteError(c *gin.Context) {
	serviceID := c.Param("service")
	errorID := c.Param("id")

	user, ok := c.Get("user")
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "could not parse user")
		return
	}

	parsedServiceID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "the provided service ID is invalid")
		return
	}

	parsedErrorID, err := primitive.ObjectIDFromHex(errorID)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "the provided error ID is invalid")
		return
	}

	service, err := service.FindOne(bson.M{"_id": parsedServiceID, "organizationId": user.(models.User).OrganizationID})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	count, err := e.EventService.DeleteError(bson.M{"_id": parsedErrorID, "ticket": service.Ticket})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if count == 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "error with id "+errorID+" does not exist")
		return
	}

	utils.RespondWithSuccess(c)
}

func (e *eventControllers) UpdateError(c *gin.Context) {
	serviceID := c.Param("service")
	errorID := c.Param("id")

	var update bson.M

	err := json.NewDecoder(c.Request.Body).Decode(&update)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	user, ok := c.Get("user")
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "could not parse user")
		return
	}

	parsedServiceID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "the provided service ID is invalid")
		return
	}

	parsedErrorID, err := primitive.ObjectIDFromHex(errorID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "the provided error ID is invalid")
		return
	}

	service, err := service.FindOne(bson.M{"_id": parsedServiceID, "organizationId": user.(models.User).OrganizationID})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = e.EventService.UpdateError(bson.M{"_id": parsedErrorID, "ticket": service.Ticket}, update)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c)
}

func GetEventController(db *mongo.Database) eventControllers {
	eventService := services.GetEventService(db)

	return eventControllers{
		EventService: &eventService,
	}
}
