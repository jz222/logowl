package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/controllers"
	"github.com/jz222/loggy/middlewares"
	"github.com/jz222/loggy/store"
)

func serviceRoutes(router *gin.RouterGroup, store store.InterfaceStore) {
	router.Use(middlewares.VerifyUserJwt(store))

	controller := controllers.GetServiceController(store)

	router.POST("/", controller.Create)
	router.DELETE("/:id", controller.Delete)
}
