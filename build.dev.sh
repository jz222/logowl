env CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/loggy .
docker build -t loggy .
rm -rf build/