package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/logowl/internal/middlewares"
	"github.com/jz222/logowl/internal/store"
)

func createSubrouter(prefix string, router *gin.Engine) *gin.RouterGroup {
	return router.Group(prefix)
}

// InitRoutes attaches all routes to the router.
func InitRoutes(router *gin.Engine, store store.InterfaceStore) {
	router.Use(middlewares.Cors)

	authRoutes(createSubrouter("/auth", router), store)
	eventRoutes(createSubrouter("/event", router), store)
	loggingRoutes(createSubrouter("/logging", router), store)
	serviceRoutes(createSubrouter("/service", router), store)
	userRoutes(createSubrouter("/user", router), store)
	organizationRoutes(createSubrouter("/organization", router), store)
}
