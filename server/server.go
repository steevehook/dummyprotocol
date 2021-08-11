package server

import (
	"crypto/ecdsa"
	"net"
	"time"

	"go.uber.org/zap"

	"github.com/steevehook/vprotocol/crypto"
	"github.com/steevehook/vprotocol/logging"
	"github.com/steevehook/vprotocol/transport"
)

type router interface {
	Switch(transport.Message) (Response, error)
}

type Response struct {
	Exited bool
	Body   interface{}
}

type Settings struct {
	Addr     string
	Router   router
	Deadline time.Duration
}

func ListenAndServe(settings Settings) (*VServer, error) {
	li, err := net.Listen("tcp", settings.Addr)
	if err != nil {
		return nil, err
	}
	srv := &VServer{
		listener: li,
		quit:     make(chan struct{}),
		exited:   make(chan struct{}),
		router:   settings.Router,
		deadline: settings.Deadline,
	}
	go srv.serve()
	return srv, nil
}

type VServer struct {
	listener net.Listener
	quit     chan struct{}
	exited   chan struct{}
	router   router
	deadline time.Duration
}

// Stop is responsible for cleanup process before application server shutdown
func (srv *VServer) Stop() error {
	logger := logging.Logger
	close(srv.quit)
	<-srv.exited
	logger.Info("server was shut down")
	return nil
}

func (srv *VServer) serve() {
	logger := logging.Logger
	logger.Info(
		"server is up and running on address",
		zap.String("addr", srv.listener.Addr().String()),
	)
	for {
		select {
		case <-srv.quit:
			// avoid accepting new connections
			logger.Info("shutting down the server")
			err := srv.listener.Close()
			if err != nil {
				logger.Error("could not close listener", zap.Error(err))
			}

			close(srv.exited)
			return
		default:
			tcpListener := srv.listener.(*net.TCPListener)
			err := tcpListener.SetDeadline(time.Now().Add(srv.deadline))
			if err != nil {
				logger.Error("failed to set listener deadline", zap.Error(err))
			}

			conn, err := tcpListener.Accept()
			if oppErr, ok := err.(*net.OpError); ok && oppErr.Timeout() {
				continue
			}
			if err != nil {
				logger.Error("failed to accept connection", zap.Error(err))
				return
			}

			// we could also manage connections here
			// and pass connection specific info
			go srv.handle(conn)
		}
	}
}

func (srv *VServer) handle(conn net.Conn) {
	logger := logging.Logger
	defer func() {
		err := conn.Close()
		if err != nil {
			logger.Error("could not close connection", zap.Error(err))
		}
		logger.Debug("connection closed")
	}()

	privateKey, err := crypto.NewECDHKey()
	if err != nil {
		logger.Error("could not create ECDH private key", zap.Error(err))
		return
	}

	var clientPublicKey *ecdsa.PublicKey
	err = crypto.DecodeECDHPublicKey(conn, &clientPublicKey)
	if err != nil {
		logger.Error("could not decode public key", zap.Error(err))
		return
	}

	err = crypto.EncodeECDHPublicKey(conn, privateKey.PublicKey)
	if err != nil {
		logger.Error("could not encode server public key", zap.Error(err))
	}
	secret := crypto.ECDHSecret(clientPublicKey, privateKey)

	scanner := transport.NewVScanner(conn)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			logger.Error("empty request line")
			continue
		}

		msg, err := transport.Decode(scanner.Bytes(), secret)
		if err != nil {
			logger.Error("could not decode data", zap.Error(err))
			continue
		}

		res, err := srv.router.Switch(msg)
		if err != nil {
			logger.Error("switch error", zap.Error(err))
			continue
		}
		if res.Exited {
			break
		}
		if res.Body == nil {
			continue
		}

		err = transport.Encode(conn, secret, msg.Operation, res.Body)
		if err != nil {
			logger.Error("could not encode response body", zap.Error(err))
		}
	}
}
