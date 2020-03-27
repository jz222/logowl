package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/keys"
	"github.com/jz222/loggy/routes"
)

var (
	port   = fmt.Sprintf(":%s", keys.GetKeys().PORT)
	router *gin.Engine
)

// Start runs the server.
func Start() {
	router = gin.Default()

	routes.InitRoutes(router)

	router.Run(port)
}
