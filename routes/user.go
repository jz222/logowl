package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/controllers"
	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/middlewares"
)

func userRoutes(router *gin.RouterGroup) {
	router.Use(middlewares.VerifyUserJwt(mongodb.GetClient()))

	controller := controllers.GetUserController(mongodb.GetClient())

	router.GET("/", controller.Get)
	router.POST("/invite", controller.Invite)
	router.DELETE("/", controller.DeleteUserAccount)
	router.DELETE("/:id", controller.Delete)
}
