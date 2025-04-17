package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/n0tB0b17/isner/internal/apis"
	"github.com/n0tB0b17/isner/internal/command"
	"github.com/n0tB0b17/isner/internal/config"
	"github.com/n0tB0b17/isner/internal/connections"
	"github.com/n0tB0b17/isner/internal/events"
	"github.com/n0tB0b17/isner/internal/storage"
)

func main() {
	cfg := config.GetConfig()
	store := storage.NewMemoryStorage()
	publisher, err := events.NewEventPublisher(cfg)
	if err != nil {
		return
	}
	defer publisher.Close()

	consumer, err := events.NewEventConsumer(cfg)
	if err != nil {
		return
	}
	defer consumer.Close()

	conn := connections.NewConnectionManager()
	registry := command.NewCommandRegistry()
	registry.InitializeDefaultCmd()

	_ = command.NewCommandExecutor(conn, publisher, registry)
	processor := command.NewCommandProcessor(publisher, store)
	processor.RegisterEventHandler(consumer)

	wsHandler := connections.NewWebSocketHandler(conn, publisher)

	s := apis.NewAPIServer(wsHandler, store)
	if err := s.Start(); err != nil {
		fmt.Printf("error starting api server: %v \n", err)
	}

	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, syscall.SIGINT, syscall.SIGTERM)
	<-osChan

	if err := s.Shutdown(); err != nil {
		fmt.Printf("error while shutting down api server: %v \n", err)
	}

	fmt.Println("server shutdown gracefully...")
}
