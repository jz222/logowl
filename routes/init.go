package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/middlewares"
)

func createSubrouter(prefix string, router *gin.Engine) *gin.RouterGroup {
	return router.Group(prefix)
}

// InitRoutes attaches all routes to the router.
func InitRoutes(router *gin.Engine) {
	router.Use(middlewares.Cors)

	authRoutes(createSubrouter("/auth", router))
	loggingRoutes(createSubrouter("/logging", router))
	projectRoutes(createSubrouter("/project", router))
}
