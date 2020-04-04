package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services"
	"github.com/jz222/loggy/store"
	"github.com/jz222/loggy/utils"
)

type organizationControllers struct {
	OrganizationService services.InterfaceOrganization
}

func (o *organizationControllers) Delete(c *gin.Context) {
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

func GetOrganizationController(store store.InterfaceStore) organizationControllers {
	organizationService := services.GetOrganizationService(store)

	return organizationControllers{
		OrganizationService: &organizationService,
	}
}
