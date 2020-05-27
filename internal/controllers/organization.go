package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jz222/logowl/internal/models"
	"github.com/jz222/logowl/internal/services"
	"github.com/jz222/logowl/internal/store"
	"github.com/jz222/logowl/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type OrganizationControllers struct {
	OrganizationService services.InterfaceOrganization
}

func (o *OrganizationControllers) Delete(c *gin.Context) {
	userData, ok := c.Get("user")
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "could not parse user data")
		return
	}

	currentUser := userData.(models.User)
	if !currentUser.IsAdmin() {
		utils.RespondWithError(c, http.StatusForbidden, "you need to be admin to perform this action")
		return
	}

	err := o.OrganizationService.Delete(userData.(models.User).OrganizationID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c)
}

func (o *OrganizationControllers) Update(c *gin.Context) {
	var organizationUpdate map[string]interface{}

	err := json.NewDecoder(c.Request.Body).Decode(&organizationUpdate)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	userData, ok := c.Get("user")
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "could not parse user data")
		return
	}

	currentUser := userData.(models.User)
	if !currentUser.IsOwnerOfOrganization() {
		utils.RespondWithError(c, http.StatusForbidden, "you need to be the owner of the organization to perform this action")
		return
	}

	filter := bson.M{"_id": currentUser.OrganizationID}
	update := bson.M{}

	isSetUp, ok := organizationUpdate["isSetUp"].(bool)
	if ok {
		update["isSetUp"] = isSetUp
	}

	_, err = o.OrganizationService.FindOneAndUpdate(filter, update)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c)
}

func GetOrganizationController(store store.InterfaceStore) OrganizationControllers {
	organizationService := services.GetOrganizationService(store)

	return OrganizationControllers{
		OrganizationService: &organizationService,
	}
}
