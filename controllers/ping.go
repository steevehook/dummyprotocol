package controllers

import (
	"github.com/steevehook/vprotocol/logging"
	"github.com/steevehook/vprotocol/server"
	"github.com/steevehook/vprotocol/transport"
)

func (router Router) ping() func(msg transport.Message) (server.Response, error) {
	return func(msg transport.Message) (server.Response, error) {
		logger := logging.Logger
		logger.Debug("ping")
		res := server.Response{
			Body: "pong",
		}
		return res, nil
	}
}
