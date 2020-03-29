package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services/user"
	"github.com/jz222/loggy/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type userControllers struct{}

var User userControllers

func (u *userControllers) Get(c *gin.Context) {
	userData, ok := c.Get("user")
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "could not parse user data")
		return
	}

	userDetails, err := user.FetchAllInformation(bson.M{"_id": userData.(models.User).ID})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, userDetails)
}

func (u *userControllers) Invite(c *gin.Context) {
	userData, ok := c.Get("user")
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "could not parse user data")
		return
	}

	if userData.(models.User).Role != "admin" {
		utils.RespondWithError(c, http.StatusForbidden, "you need to be admin to invite new users")
		return
	}

	var newUser models.User

	err := json.NewDecoder(c.Request.Body).Decode(&newUser)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	newUser.OrganizationID = userData.(models.User).OrganizationID

	persistedUser, err := user.Invite(newUser)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, persistedUser)
}

func (u *userControllers) Delete(c *gin.Context) {
	userData, ok := c.Get("user")
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "could not parse user data")
		return
	}

	if userData.(models.User).IsOrganizationOwner {
		utils.RespondWithError(c, http.StatusForbidden, "you can not delete your account as organization owner")
		return
	}

	deleteCount, err := user.Delete(userData.(models.User).ID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if deleteCount == 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "could not delete user")
		return
	}

	utils.RespondWithSuccess(c)
}
