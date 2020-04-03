package main

import (
	"github.com/jz222/loggy/server"
)

func main() {
	server := server.CreateInstance()

	server.Start()
}
