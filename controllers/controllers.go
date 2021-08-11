package controllers

import (
	"fmt"

	"github.com/steevehook/vprotocol/server"
	"github.com/steevehook/vprotocol/transport"
)

// Operations
const (
	pingOperation       = "ping"
	disconnectOperation = "disconnect"
)

type operator func(message transport.Message) (server.Response, error)

type Router struct {
	operations map[string]operator
}

func NewRouter() Router {
	router := Router{}
	router.operations = map[string]operator{
		pingOperation: router.ping(),
	}
	return router
}

func (router Router) Switch(msg transport.Message) (server.Response, error) {
	if msg.Operation == disconnectOperation {
		res := server.Response{Exited: true}
		return res, nil
	}

	operation, ok := router.operations[msg.Operation]
	if !ok {
		return server.Response{}, fmt.Errorf("operation '%s' not found", msg.Operation)
	}

	res, err := operation(msg)
	if err != nil {
		return server.Response{}, err
	}

	return res, nil
}
