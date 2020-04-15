package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/internal/controllers"
	"github.com/jz222/loggy/internal/middlewares"
	"github.com/jz222/loggy/internal/store"
)

func organizationRoutes(router *gin.RouterGroup, store store.InterfaceStore) {
	router.Use(middlewares.VerifyUserJwt(store))

	controller := controllers.GetOrganizationController(store)

	router.DELETE("/", controller.Delete)
}
