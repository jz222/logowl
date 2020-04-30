package models

// Keys contains all environment variables.
type Keys struct {
	MONGO_URI              string
	MONGO_DB_NAME          string
	PORT                   string
	SECRET                 string
	CLIENT_URL             string
	TOTAL_MONTHLY_REQUESTS int
	IS_SELFHOSTED          bool
}
