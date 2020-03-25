package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/controllers"
	"github.com/jz222/loggy/middlewares"
)

func serviceRoutes(router *gin.RouterGroup) {
	router.Use(middlewares.VerifyUserJwt)

	router.POST("/", controllers.Service.Create)
	router.DELETE("/:id", controllers.Service.Delete)
}
