package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/logowl/internal/controllers"
	"github.com/jz222/logowl/internal/middlewares"
	"github.com/jz222/logowl/internal/store"
)

func serviceRoutes(router *gin.RouterGroup, store store.InterfaceStore) {
	router.Use(middlewares.VerifyUserJwt(store))

	controller := controllers.GetServiceController(store)

	router.POST("/", controller.Create)
	router.PUT("/:id", controller.Edit)
	router.DELETE("/:id", controller.Delete)
}
