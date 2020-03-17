package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services/auth"
	"github.com/jz222/loggy/services/organization"
	"github.com/jz222/loggy/services/user"
	"github.com/jz222/loggy/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type authControllers struct{}

var Auth authControllers

func (a *authControllers) Setup(c *gin.Context) {
	var setup models.Setup

	err := json.NewDecoder(c.Request.Body).Decode(&setup)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	organizationID, err := organization.Create(setup.Organization)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	setup.User.OrganizationID = organizationID

	_, err = user.Create(setup.User)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
	}
}

func (a *authControllers) SignIn(c *gin.Context) {
	var credentials models.Credentials

	err := json.NewDecoder(c.Request.Body).Decode(&credentials)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	persistedUser, err := user.FetchAllInformation(bson.M{"email": credentials.Email})
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "the provided email and password don't match")
		return
	}

	passwordIsValid := persistedUser.VerifyPassword(credentials.Password)
	if !passwordIsValid {
		utils.RespondWithError(c, http.StatusUnauthorized, "the provided email and password don't match")
		return
	}

	jwt, expirationTime, err := auth.CreateJWT(persistedUser.ID.Hex())
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	persistedUser.Password = ""

	response := models.SignInResponse{
		User:           persistedUser,
		JWT:            jwt,
		ExpirationTime: expirationTime,
	}

	utils.RespondWithJSON(c, response)
}
