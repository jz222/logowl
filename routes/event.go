package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/controllers"
	"github.com/jz222/loggy/middlewares"
)

func eventRoutes(router *gin.RouterGroup) {
	router.Use(middlewares.VerifyUserJwt)

	router.GET(":service/error/all", controllers.Event.GetErrors)
	router.GET(":service/error/all/:pointer", controllers.Event.GetErrors)
	router.DELETE(":service/error/:id", controllers.Event.DeleteError)
}
