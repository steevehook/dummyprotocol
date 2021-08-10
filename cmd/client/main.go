package main

import (
	"log"

	"github.com/steevehook/vprotocol/client"
)

func main() {
	c := client.New()
	err := c.Connect("localhost:8080")
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
}
