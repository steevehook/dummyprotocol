package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/steevehook/vprotocol/client"
)

func main() {
	// SINGLE CLIENT USES THE SAME CONNECTION CONCURRENTLY
	c := client.VClient{}
	err := c.Connect("localhost:8080")
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(10)
	// test for write broken pipe / add some rate limiting
	// test with many concurrent persistent connections
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			msg, err := c.Ping()
			if err != nil {
				log.Fatalf("could not ping server: %v", err)
			}
			fmt.Println("ping:", i, msg)
		}(i + 1)
	}

	wg.Wait()
	err = c.Disconnect()
	if err != nil {
		log.Fatalf("could not disconnect from server: %v", err)
	}

	err = c.Connect("localhost:8080")
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}

	// RECONNECT TO THE SAME CLIENT
	msg, err := c.Ping()
	if err != nil {
		log.Fatalf("could not ping server: %v", err)
	}
	fmt.Println("reconnected ping:", msg)

	err = c.Disconnect()
	if err != nil {
		log.Fatalf("could not disconnect from server: %v", err)
	}

	// MULTIPLE CLIENTS
	var clients []client.VClient
	for i := 0; i < 10; i++ {
		clients = append(clients, client.VClient{})
	}

	for _, c := range clients {
		wg.Add(1)
		go func(c client.VClient) {
			defer wg.Done()

			err := c.Connect("localhost:8080")
			if err != nil {
				log.Fatalf("could not connect to server: %v", err)
			}

			msg, err := c.Ping()
			if err != nil {
				log.Fatalf("could not ping server: %v", err)
			}
			fmt.Println("concurrent ping:", msg)

			err = c.Disconnect()
			if err != nil {
				log.Fatalf("could not disconnect from server: %v", err)
			}
		}(c)
	}
	wg.Wait()
}
