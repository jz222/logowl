package main

import (
	"github.com/jz222/loggy/internal/server"
)

func main() {
	server := server.CreateInstance()

	server.Start()
}
