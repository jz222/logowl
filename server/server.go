package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/keys"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/routes"
	"github.com/jz222/loggy/store"
)

type instance struct {
	Keys   models.Keys
	Server *gin.Engine
}

// Start runs the server.
func (s *instance) Start() {
	storeInstance := store.GetStore()
	storeInstance.Connect()
	defer storeInstance.Disconnect()

	s.Keys = keys.GetKeys()
	s.Server = gin.Default()

	routes.InitRoutes(s.Server, storeInstance)

	port := fmt.Sprintf(":%s", s.Keys.PORT)

	s.Server.Run(port)
}

// CreateInstance creates a new server instance.
func CreateInstance() *instance {
	return &instance{}
}
