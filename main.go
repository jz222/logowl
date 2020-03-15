package main

import (
	"github.com/jz222/loggy/keys"
	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/server"
)

func init() {
	keys.GetKeys()
	mongodb.InitiateDatabase()
}

func main() {
	server.Start()
}
