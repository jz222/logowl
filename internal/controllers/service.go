package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/internal/models"
	"github.com/jz222/loggy/internal/services"
	"github.com/jz222/loggy/internal/store"
	"github.com/jz222/loggy/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceController struct {
	ServiceService services.InterfaceService
}

func (s *ServiceController) Create(c *gin.Context) {
	var newService models.Service

	err := json.NewDecoder(c.Request.Body).Decode(&newService)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	userData, ok := c.Get("user")
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "could not parse user data")
		return
	}

	newService.OrganizationID = userData.(models.User).OrganizationID

	createdService, err := s.ServiceService.Create(newService)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, createdService)
}

func (s *ServiceController) Delete(c *gin.Context) {
	id := c.Param("id")

	serviceID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	userData, ok := c.Get("user")
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "could not parse user data")
		return
	}

	filter := bson.M{"_id": serviceID, "organizationId": userData.(models.User).OrganizationID}

	count, err := s.ServiceService.Delete(filter)
	if err != nil {
		utils.RespondWithError(c, http.StatusForbidden, err.Error())
		return
	}

	if count == 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "the service with the id "+id+" does not exist")
		return
	}

	utils.RespondWithSuccess(c)
}

func (s *ServiceController) Edit(c *gin.Context) {
	id := c.Param("id")

	var serviceUpdate map[string]interface{}

	err := json.NewDecoder(c.Request.Body).Decode(&serviceUpdate)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	serviceId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	userData, ok := c.Get("user")
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "could not parse user data")
		return
	}

	filter := bson.M{"_id": serviceId, "organizationId": userData.(models.User).OrganizationID}
	update := bson.M{}

	slackWebHookURL, ok := serviceUpdate["slackWebhookURL"].(string)
	if ok {
		update["slackWebhookURL"] = slackWebHookURL
	}

	_, err = s.ServiceService.FindOneAndUpdate(filter, update)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c)
}

func GetServiceController(store store.InterfaceStore) ServiceController {
	serviceService := services.GetServiceService(store)

	return ServiceController{
		ServiceService: &serviceService,
	}
}
