package utils

import (
	"crypto/rand"
	"fmt"
	"log"
)

func GenerateTicket() (string, error) {
	buf := make([]byte, 25)

	_, err := rand.Read(buf)
	if err != nil {
		log.Println("Failed to create ticket")
		return "", err
	}

	ticket := fmt.Sprintf("%X", buf)

	return ticket, nil
}
