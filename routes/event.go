package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/controllers"
	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/middlewares"
)

func eventRoutes(router *gin.RouterGroup) {
	router.Use(middlewares.VerifyUserJwt(mongodb.GetClient()))

	controller := controllers.GetEventController(mongodb.GetClient())

	router.GET(":service/error/:id", controller.GetError)
	router.GET(":service/errors/", controller.GetErrors)
	router.GET(":service/errors/:pointer", controller.GetErrors)
	router.PUT(":service/error/:id", controller.UpdateError)
	router.DELETE(":service/error/:id", controller.DeleteError)
}
