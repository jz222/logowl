package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/controllers"
)

func eventRoutes(router *gin.RouterGroup) {
	router.GET("/error/all", controllers.Event.GetErrors)
	router.GET("/error/all/:pointer", controllers.Event.GetErrors)
}
