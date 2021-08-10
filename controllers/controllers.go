package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io"

	"github.com/steevehook/vprotocol/logging"
	"github.com/steevehook/vprotocol/models"
	"github.com/steevehook/vprotocol/transport"
)

// Operations
const (
	pingOperation       = "ping"
	connectOperation    = "connect"
	disconnectOperation = "disconnect"
	decodeOperation     = "decode_request"
)

type request struct {
	Operation string          `json:"operation"`
	Body      json.RawMessage `json:"body"`
}

// ConfigManager represents the application configuration manager
type ConfigManager interface {
	GetLoggerLevel() string
	GetLoggerOutput() []string
}

type operator func(io.Writer, request) error

// Router represents the Event Bus operation router switch
type Router struct {
	operations map[string]operator
	cfg        ConfigManager
}

// NewRouter creates a new instance of Router switch operation
func NewRouter(cfg ConfigManager) Router {
	router := Router{
		cfg: cfg,
	}
	router.operations = map[string]operator{
		connectOperation: router.connect(),
		pingOperation: func(w io.Writer, _ request) error {
			transport.SendJSON(w, pingOperation, nil)
			return nil
		},
		disconnectOperation: func(w io.Writer, _ request) error {
			transport.SendJSON(w, disconnectOperation, nil)
			return nil
		},
	}
	return router
}

// Switch represents the switch between Event Bus operations
func (router Router) Switch(w io.Writer, r io.Reader) (bool, error) {
	var req request
	err := transport.Decode(r, &req)
	if err != nil {
		transport.SendError(w, decodeOperation, models.InvalidJSONError{})
		return false, models.Error{Message: err.Error()}
	}

	notFoundErr := models.OperationNotFoundError{}
	operation, ok := router.operations[req.Operation]
	if !ok {
		transport.SendError(w, decodeOperation, notFoundErr)
		return false, notFoundErr
	}

	if req.Operation == pingOperation {
		return false, operation(w, req)
	}
	if req.Operation == disconnectOperation {
		return true, operation(w, req)
	}

	// in case we don't have a ping or exit, execute whatever operation was matched
	return false, operation(w, req)
}

func parseReq(r request, body interface{}) error {
	if r.Body == nil {
		return errors.New("body is nil")
	}
	err := transport.Decode(bytes.NewReader(r.Body), body)
	if err != nil {
		logging.Logger.Debug("could not decode request body", zap.Error(err))
		return models.InvalidJSONError{}
	}
	return nil
}
