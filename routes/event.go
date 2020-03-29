package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/controllers"
	"github.com/jz222/loggy/middlewares"
)

func eventRoutes(router *gin.RouterGroup) {
	router.Use(middlewares.VerifyUserJwt)

	router.GET(":service/error/:id", controllers.Event.GetError)
	router.GET(":service/errors/", controllers.Event.GetErrors)
	router.GET(":service/errors/:pointer", controllers.Event.GetErrors)

	router.PUT(":service/error/:id", controllers.Event.UpdateError)

	router.DELETE(":service/error/:id", controllers.Event.DeleteError)
}
