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

type userControllers struct {
	UserService services.InterfaceUser
}

func (u *userControllers) Get(c *gin.Context) {
	userData, ok := c.Get("user")
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "could not parse user data")
		return
	}

	userDetails, err := u.UserService.FetchAllInformation(bson.M{"_id": userData.(models.User).ID})
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

	persistedUser, err := u.UserService.Invite(newUser)
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

	if userData.(models.User).Role != "admin" {
		utils.RespondWithError(c, http.StatusForbidden, "you need to be admin to delete users")
		return
	}

	userID := c.Param("id")

	parsedUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "the provided user ID is invalid")
		return
	}

	deleteCount, err := u.UserService.Delete(bson.M{"_id": parsedUserID, "organizationId": userData.(models.User).OrganizationID})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if deleteCount == 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "the user with the ID "+userID+" does not exist")
		return
	}

	utils.RespondWithSuccess(c)
}

func (u *userControllers) DeleteUserAccount(c *gin.Context) {
	userData, ok := c.Get("user")
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "could not parse user data")
		return
	}

	if userData.(models.User).IsOrganizationOwner {
		utils.RespondWithError(c, http.StatusForbidden, "you can not delete your account as organization owner")
		return
	}

	deleteCount, err := u.UserService.Delete(bson.M{"_id": userData.(models.User).ID})
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

func GetUserController(store store.InterfaceStore) userControllers {
	userService := services.GetUserService(store)

	return userControllers{
		UserService: &userService,
	}
}

func GetUserControllerMock() userControllers {
	userService := services.GetUserServiceMock()

	return userControllers{
		UserService: &userService,
	}
}
