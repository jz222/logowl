package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/controllers"
	"github.com/jz222/loggy/middlewares"
	"github.com/jz222/loggy/store"
)

func eventRoutes(router *gin.RouterGroup, store store.InterfaceStore) {
	router.Use(middlewares.VerifyUserJwt(store))

	controller := controllers.GetEventController(store)

	router.GET(":service/error/:id", controller.GetError)
	router.GET(":service/errors/", controller.GetErrors)
	router.GET(":service/errors/:pointer", controller.GetErrors)
	router.PUT(":service/error/:id", controller.UpdateError)
	router.DELETE(":service/error/:id", controller.DeleteError)
}
