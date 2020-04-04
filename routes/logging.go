package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/controllers"
	"github.com/jz222/loggy/libs/mongodb"
)

func loggingRoutes(router *gin.RouterGroup) {
	controller := controllers.GetLoggingController(mongodb.GetClient())

	router.POST("/error", controller.RegisterError)
}
