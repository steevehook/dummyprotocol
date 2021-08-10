package server

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"time"

	"go.uber.org/zap"

	"github.com/steevehook/vprotocol/logging"
)

type router interface {
	Switch(io.Writer, io.Reader) (bool, error)
}

type Settings struct {
	Addr     string
	Router   router
	Deadline time.Duration
}

type Server struct {
	listener net.Listener
	quit     chan struct{}
	exited   chan struct{}
	router   router
	deadline time.Duration
}

func ListenAndServe(settings Settings) (*Server, error) {
	li, err := net.Listen("tcp", settings.Addr)
	if err != nil {
		return nil, err
	}
	srv := &Server{
		listener: li,
		quit:     make(chan struct{}),
		exited:   make(chan struct{}),
		router:   settings.Router,
		deadline: settings.Deadline,
	}
	go srv.serve()
	return srv, nil
}

// Stop is responsible for cleanup process before application server shutdown
func (srv *Server) Stop() error {
	logger := logging.Logger
	close(srv.quit)
	<-srv.exited
	logger.Info("server was shut down")
	return nil
}

func (srv *Server) serve() {
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

			go srv.handle(conn)
		}
	}
}

func (srv *Server) handle(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	scanner := bufio.NewScanner(conn)
	logger := logging.Logger
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			logger.Error("empty request line")
			continue
		}

		exited, err := srv.router.Switch(conn, bytes.NewReader(scanner.Bytes()))
		if err != nil {
			logger.Error("switch error", zap.Error(err))
		}
		if exited {
			break
		}
	}
}
