package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/keys"
	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/routes"
)

type instance struct {
	Keys   models.Keys
	Server *gin.Engine
}

// Start runs the server.
func (s *instance) Start() {
	db := mongodb.GetClient()
	defer db.Client().Disconnect(context.TODO())

	s.Keys = keys.GetKeys()
	s.Server = gin.Default()

	routes.InitRoutes(s.Server)

	port := fmt.Sprintf(":%s", s.Keys.PORT)

	s.Server.Run(port)
}

// CreateInstance creates a new server instance.
func CreateInstance() *instance {
	return &instance{}
}
