package keys

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	"github.com/jz222/loggy/internal/models"
)

const (
	SESSION_TIMEOUT_IN_HOURS = 7
)

var (
	instance models.Keys
	once     sync.Once
)

func loadEnvAsString(key string) string {
	godotenv.Load()

	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}

	if key == "MAILGUN_PRIVATE_KEY" {
		log.Println("[INFO] the environment variable MAILGUN_PRIVATE_KEY was not provided. Sending emails will therefore not be available")
		return ""
	}

	if key == "MAILGUN_DOMAIN" {
		log.Println("[INFO] the environment variable MAILGUN_DOMAIN was not provided. Sending emails will therefore not be available")
		return ""
	}

	log.Fatal("❌ Failed to load environment variable: ", key)

	return ""
}

func loadEnvAsInt(key string) int {
	godotenv.Load()

	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatal("❌ Failed to load environment variable: ", key)
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal("❌ Failed to parse environment variable " + key + ". Please provide a valid number greater than zero")
	}

	return parsed
}

var envVariables = models.Keys{
	MONGO_URI:             loadEnvAsString("MONGO_URI"),
	MONGO_DB_NAME:         loadEnvAsString("MONGO_DB_NAME"),
	MAILGUN_PRIVATE_KEY:   loadEnvAsString("MAILGUN_PRIVATE_KEY"),
	MAILGUN_DOMAIN:        loadEnvAsString("MAILGUN_DOMAIN"),
	PORT:                  loadEnvAsString("PORT"),
	SECRET:                loadEnvAsString("SECRET"),
	CLIENT_URL:            loadEnvAsString("CLIENT_URL"),
	MONTHLY_REQUEST_LIMIT: loadEnvAsInt("MONTHLY_REQUEST_LIMIT"),
	IS_SELFHOSTED:         loadEnvAsString("IS_SELFHOSTED") == "true",
}

// GetKeys returns all environment variables. It can also be
// executed to just load all environment variables.
func GetKeys() models.Keys {
	once.Do(func() {
		instance = envVariables
	})

	return instance
}
