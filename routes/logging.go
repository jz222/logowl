package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/controllers"
	"github.com/jz222/loggy/store"
)

func loggingRoutes(router *gin.RouterGroup, store store.InterfaceStore) {
	controller := controllers.GetLoggingController(store)

	router.POST("/error", controller.RegisterError)
}
