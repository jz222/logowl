package server

import (
	"fmt"
	"net/http"
	"time"

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

	server := &http.Server{
		Addr:              port,
		Handler:           router,
		ReadTimeout:       2 * time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 4 * time.Second,
		MaxHeaderBytes:    1 << 40,
	}

	server.ListenAndServe()
}
