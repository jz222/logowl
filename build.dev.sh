env CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/logowl ./cmd/logowl/
docker build -t logowl .
rm -rf build/