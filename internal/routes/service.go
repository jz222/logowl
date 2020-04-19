package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/internal/controllers"
	"github.com/jz222/loggy/internal/middlewares"
	"github.com/jz222/loggy/internal/store"
)

func serviceRoutes(router *gin.RouterGroup, store store.InterfaceStore) {
	router.Use(middlewares.VerifyUserJwt(store))

	controller := controllers.GetServiceController(store)

	router.POST("/", controller.Create)
	router.PUT("/:id", controller.Edit)
	router.DELETE("/:id", controller.Delete)
}
