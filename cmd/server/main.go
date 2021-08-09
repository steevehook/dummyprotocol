package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/steevehook/vprotocol/config"
	"github.com/steevehook/vprotocol/controllers"
	"github.com/steevehook/vprotocol/logging"
	"github.com/steevehook/vprotocol/server"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "config file path")
	flag.Parse()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	cfg, err := config.NewManager(*configPath)
	if err != nil {
		log.Fatalf("could not create config manager: %v", err)
	}

	loggerSettings := logging.Settings{
		Level:  cfg.GetLoggerLevel(),
		Output: cfg.GetLoggerOutput(),
	}
	err = logging.Init(loggerSettings)
	if err != nil {
		log.Fatal("could not initialize logger: ", err)
	}

	router := controllers.NewRouter(cfg)
	serverSettings := server.Settings{
		Addr:     cfg.GetServerAddr(),
		Router:   router,
		Deadline: cfg.GetServerDeadline(),
	}
	srv, err := server.ListenAndServe(serverSettings)
	if err != nil {
		log.Fatalf("could not create server listener: %v", err)
	}

	select {
	case <-stop:
		err := srv.Stop()
		if err != nil {
			log.Fatalf("could not stop the server: %v", err)
		}
	}
}
