package keys

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/jz222/loggy/internal/models"
)

var (
	instance models.Keys
	once     sync.Once
)

func loadEnv(key string) string {
	if err := godotenv.Load(); err != nil {
		log.Println("Could not load .env file. Proceeding to read environment variable")
	}

	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return ""
}

var envVariables = models.Keys{
	MONGO_URI:     loadEnv("MONGO_URI"),
	MONGO_DB_NAME: loadEnv("MONGO_DB_NAME"),
	PORT:          loadEnv("PORT"),
	SECRET:        loadEnv("SECRET"),
	IS_SELFHOSTED: loadEnv("IS_SELFHOSTED") == "true",
}

// GetKeys returns all environment variables. It can also be
// executed to just load all environment variables.
func GetKeys() models.Keys {
	once.Do(func() {
		instance = envVariables
	})

	return instance
}
