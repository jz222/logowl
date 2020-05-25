package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/logowl/internal/controllers"
	"github.com/jz222/logowl/internal/store"
)

func authRoutes(router *gin.RouterGroup, store store.InterfaceStore) {
	controller := controllers.GetAuthControllers(store)

	router.POST("/setup", controller.Setup)
	router.POST("/signup", controller.SignUp)
	router.POST("/signin", controller.SignIn)
	router.POST("/resetpassword", controller.ResetPassword)
	router.POST("/setnewpassword", controller.SetNewPassword)
}
