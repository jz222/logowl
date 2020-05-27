package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jz222/logowl/internal/keys"
	"github.com/jz222/logowl/internal/models"
	"github.com/jz222/logowl/internal/routes"
	"github.com/jz222/logowl/internal/store"
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

	srv := &http.Server{
		Addr:    port,
		Handler: s.Server,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server with error: ", err.Error())
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("failed to shut down server with error: ", err.Error())
	}

	log.Println("✅ Server shut down gracefully")
}

// CreateInstance creates a new server instance.
func CreateInstance() *instance {
	return &instance{}
}
