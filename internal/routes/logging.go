package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/internal/controllers"
	"github.com/jz222/loggy/internal/store"
)

func loggingRoutes(router *gin.RouterGroup, store store.InterfaceStore) {
	controller := controllers.GetLoggingController(store)

	router.POST("/error", controller.RegisterError)
	router.POST("/analytics", controller.RegisterAnalyticEvent)
}
