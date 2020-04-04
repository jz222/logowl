package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/keys"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services"
	"github.com/jz222/loggy/services/auth"
	"github.com/jz222/loggy/services/organization"
	"github.com/jz222/loggy/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type authControllers struct {
	UserService services.InterfaceUser
}

func (a *authControllers) Setup(c *gin.Context) {
	var setup models.Setup

	err := json.NewDecoder(c.Request.Body).Decode(&setup)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if keys.GetKeys().IS_SELFHOSTED {
		exists, err := organization.CheckPresence(bson.M{})
		if err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}

		if exists {
			utils.RespondWithError(c, http.StatusForbidden, "there can only be one organization in self-hosted mode")
			return
		}
	}

	userExists, err := a.UserService.CheckPresence(bson.M{"email": setup.User.Email})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userExists {
		utils.RespondWithError(c, http.StatusForbidden, "could not create user")
		return
	}

	organizationID, err := organization.Create(setup.Organization)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	setup.User.OrganizationID = organizationID
	setup.User.IsOrganizationOwner = true
	setup.User.Role = "admin"

	_, err = a.UserService.Create(setup.User)
	if err != nil {
		fmt.Println(err.Error())
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
	}

	utils.RespondWithSuccess(c)
}

func (a *authControllers) SignUp(c *gin.Context) {
	var credentials models.Credentials

	err := json.NewDecoder(c.Request.Body).Decode(&credentials)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if credentials.Email == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "email address was not provided")
		return
	}

	if credentials.Password == "" || len(credentials.Password) < 12 {
		utils.RespondWithError(c, http.StatusBadRequest, "password was not provided or is invalid")
		return
	}

	filter := bson.M{"email": credentials.Email, "inviteCode": credentials.InviteCode, "isVerified": false}
	update := bson.M{"password": credentials.Password, "isVerified": true}

	err = a.UserService.Update(filter, update)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	userData, err := a.UserService.FetchAllInformation(bson.M{"email": credentials.Email})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jwt, expirationTime, err := auth.CreateJWT(userData.ID.Hex())
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	userData.Password = ""

	response := models.SignInResponse{
		User:           userData,
		JWT:            jwt,
		ExpirationTime: expirationTime,
	}

	utils.RespondWithJSON(c, response)
}

func (a *authControllers) SignIn(c *gin.Context) {
	var credentials models.Credentials

	err := json.NewDecoder(c.Request.Body).Decode(&credentials)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	persistedUser, err := a.UserService.FetchAllInformation(bson.M{"email": credentials.Email})
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

func GetAuthControllers(db *mongo.Database) authControllers {
	userService := services.GetUserService(db)

	return authControllers{
		UserService: &userService,
	}
}
