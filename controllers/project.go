package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services/project"
	"github.com/jz222/loggy/utils"
)

type projectControllers struct{}

var Project projectControllers

func (p *projectControllers) Create(c *gin.Context) {
	var newProject models.Project

	err := json.NewDecoder(c.Request.Body).Decode(&newProject)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	userData, ok := c.Get("user")
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "could not parse user data")
		return
	}

	newProject.OrganizationID = userData.(models.User).OrganizationID

	createdProject, err := project.Create(newProject)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, createdProject)
}
