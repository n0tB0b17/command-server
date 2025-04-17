package main

import (
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

	connections.NewWebSocketHandler(conn, publisher)
}
