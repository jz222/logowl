package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services/user"
	"github.com/jz222/loggy/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type userControllers struct{}

var User userControllers

func (u *userControllers) GetUser(c *gin.Context) {
	userData, ok := c.Get("user")
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "user ID is not available")
		return
	}

	userDetails, err := user.FetchAllInformation(bson.M{"_id": userData.(models.User).ID})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	userDetails.Password = ""

	utils.RespondWithJSON(c, userDetails)
}
