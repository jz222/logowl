package main

import (
	"github.com/jz222/logowl/internal/server"
)

func main() {
	server := server.CreateInstance()

	server.Start()
}
