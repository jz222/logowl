package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/controllers"
	"github.com/jz222/loggy/middlewares"
)

func userRoutes(router *gin.RouterGroup) {
	router.Use(middlewares.VerifyUserJwt)

	router.GET("/", controllers.User.Get)
	router.POST("/invite", controllers.User.Invite)
	router.DELETE("/", controllers.User.DeleteUserAccount)
	router.DELETE("/:id", controllers.User.Delete)
}
