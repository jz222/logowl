package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/controllers"
	"github.com/jz222/loggy/libs/mongodb"
)

func authRoutes(router *gin.RouterGroup) {
	controller := controllers.GetAuthControllers(mongodb.GetClient())

	router.POST("/setup", controller.Setup)
	router.POST("/signup", controller.SignUp)
	router.POST("/signin", controller.SignIn)
}
