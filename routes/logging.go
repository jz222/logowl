package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/controllers"
)

func loggingRoutes(router *gin.RouterGroup) {
	router.POST("/error", controllers.Logging.RegisterError)
	router.GET("/error", controllers.Logging.LoadAllErrors)
}
