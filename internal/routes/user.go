package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/logowl/internal/controllers"
	"github.com/jz222/logowl/internal/middlewares"
	"github.com/jz222/logowl/internal/store"
)

func userRoutes(router *gin.RouterGroup, store store.InterfaceStore) {
	router.Use(middlewares.VerifyUserJwt(store))

	controller := controllers.GetUserController(store)

	router.GET("/", controller.Get)
	router.POST("/invite", controller.Invite)
	router.DELETE("/", controller.DeleteUserAccount)
	router.DELETE("/:id", controller.Delete)
}
