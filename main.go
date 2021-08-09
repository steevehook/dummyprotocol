package main

import (
	"log"
	"time"

	"github.com/steevehook/dummyprotocol/proto"
)

func main() {
	go func() {
		err := proto.Listen(8080)
		if err != nil {
			log.Fatalf("could not start server: %v", err)
		}
	}()

	c := proto.NewClient()
	err := c.Connect("127.0.0.1", 8080)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	//c.Ping()
	//c.Disconnect()
	time.Sleep(time.Second * 20)
}
