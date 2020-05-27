package models

// Keys contains all environment variables.
type Keys struct {
	MONGO_URI             string
	MONGO_DB_NAME         string
	MAILGUN_PRIVATE_KEY   string
	MAILGUN_DOMAIN        string
	PORT                  string
	SECRET                string
	CLIENT_URL            string
	MONTHLY_REQUEST_LIMIT int
	IS_SELFHOSTED         bool
}
