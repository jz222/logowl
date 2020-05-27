package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/logowl/internal/controllers"
	"github.com/jz222/logowl/internal/middlewares"
	"github.com/jz222/logowl/internal/store"
)

func organizationRoutes(router *gin.RouterGroup, store store.InterfaceStore) {
	router.Use(middlewares.VerifyUserJwt(store))

	controller := controllers.GetOrganizationController(store)

	router.DELETE("/", controller.Delete)
	router.PUT("/", controller.Update)
}
