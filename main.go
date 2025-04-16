package main

import (
	"fmt"

	"github.com/n0tB0b17/isner/internal/config"
	"github.com/n0tB0b17/isner/internal/events"
)

func main() {
	cfg := config.GetConfig()
	fmt.Println("starter for command and control server")
	publisher, _ := events.NewEventPublisher(cfg)
	consumer, _ := events.NewEventConsumer(cfg)
	publisher.Publish(events.EventClientConnected, nil)

	consumer.RegisterEvent(events.EventClientConnected, func(event events.RabbitMQEvent) error {
		return nil
	})

	consumer.StartConsuming("client_events", []string{"client.connected"})
}
