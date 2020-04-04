package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/controllers"
	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/middlewares"
)

func eventRoutes(router *gin.RouterGroup) {
	router.Use(middlewares.VerifyUserJwt)

	eventController := controllers.GetEventController(mongodb.GetClient())

	router.GET(":service/error/:id", eventController.GetError)
	router.GET(":service/errors/", eventController.GetErrors)
	router.GET(":service/errors/:pointer", eventController.GetErrors)
	router.PUT(":service/error/:id", eventController.UpdateError)
	router.DELETE(":service/error/:id", eventController.DeleteError)
}
