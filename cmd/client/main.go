package main

import (
	"fmt"
	"log"

	"github.com/steevehook/vprotocol/client"
)

func main() {
	c := client.VClient{}

	err := c.Connect("localhost:8080")
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}

	msg, err := c.Ping()
	if err != nil {
		log.Fatalf("could not ping server: %v", err)
	}

	fmt.Println("response from server:", msg)
}
