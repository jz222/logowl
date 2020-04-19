export MONGO_URI=(grep MONGO_URI .env | cut -d '=' -f2)
export MONGO_DB_NAME=(grep MONGO_DB_NAME .env | cut -d '=' -f2)
export PORT=(grep PORT .env | cut -d '=' -f2)
export SECRET=(grep SECRET .env | cut -d '=' -f2)
export CLIENT_URL=(grep CLIENT_URL .env | cut -d '=' -f2)
export IS_SELFHOSTED=(grep IS_SELFHOSTED .env | cut -d '=' -f2)
go test ./internal/unittests/... -v