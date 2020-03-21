package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services/service"
	"github.com/jz222/loggy/utils"
)

type serviceController struct{}

var Service serviceController

func (p *serviceController) Create(c *gin.Context) {
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

	createdService, err := service.Create(newService)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, createdService)
}
