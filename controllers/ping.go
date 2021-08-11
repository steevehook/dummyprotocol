package controllers

import (
	"fmt"
	"github.com/steevehook/vprotocol/server"
	"github.com/steevehook/vprotocol/transport"
)

func (router Router) ping() func(msg transport.Message) (server.Response, error) {
	return func(msg transport.Message) (server.Response, error) {
		fmt.Println("pinged")
		res := server.Response{
			Body: "pong",
		}
		fmt.Println("time to pong")
		return res, nil
	}
}
