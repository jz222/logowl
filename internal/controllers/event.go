package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/internal/models"
	"github.com/jz222/loggy/internal/services"
	"github.com/jz222/loggy/internal/store"
	"github.com/jz222/loggy/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventControllers struct {
	EventService   services.InterfaceEvent
	ServiceService services.InterfaceService
}

func (e *EventControllers) GetError(c *gin.Context) {
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

	persistedService, err := e.ServiceService.FindOne(bson.M{"_id": parsedServiceID, "organizationId": userData.(models.User).OrganizationID})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	errorEvent, err := e.EventService.GetError(bson.M{"_id": parsedErrorID, "ticket": persistedService.Ticket}, userData.(models.User).ID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, errorEvent)
}

func (e *EventControllers) GetErrors(c *gin.Context) {
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

	requestedService, err := e.ServiceService.FindOne(bson.M{"_id": parsedServiceID, "organizationId": organizationID})
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

func (e *EventControllers) DeleteError(c *gin.Context) {
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

	service, err := e.ServiceService.FindOne(bson.M{"_id": parsedServiceID, "organizationId": user.(models.User).OrganizationID})
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

func (e *EventControllers) UpdateError(c *gin.Context) {
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

	service, err := e.ServiceService.FindOne(bson.M{"_id": parsedServiceID, "organizationId": user.(models.User).OrganizationID})
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

func (e *EventControllers) GetAnalytics(c *gin.Context) {
	serviceID := c.Param("service")
	mode := c.Query("mode")

	if mode == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "the query parameter mode was not provided")
		return
	}

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

	service, err := e.ServiceService.FindOne(bson.M{"_id": parsedServiceID, "organizationId": user.(models.User).OrganizationID})
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	analyticInsights, err := e.EventService.GetAnalytics(service.Ticket, mode)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, analyticInsights)
}

func GetEventController(db store.InterfaceStore) EventControllers {
	eventService := services.GetEventService(db)
	serviceService := services.GetServiceService(db)

	return EventControllers{
		EventService:   &eventService,
		ServiceService: &serviceService,
	}
}
