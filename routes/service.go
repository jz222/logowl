package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/controllers"
	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/middlewares"
)

func serviceRoutes(router *gin.RouterGroup) {
	router.Use(middlewares.VerifyUserJwt(mongodb.GetClient()))

	controller := controllers.GetServiceController(mongodb.GetClient())

	router.POST("/", controller.Create)
	router.DELETE("/:id", controller.Delete)
}
