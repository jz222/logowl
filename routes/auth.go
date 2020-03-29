package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/controllers"
)

func authRoutes(router *gin.RouterGroup) {
	router.POST("/setup", controllers.Auth.Setup)
	router.POST("/signup", controllers.Auth.SignUp)
	router.POST("/signin", controllers.Auth.SignIn)
}
