package events

import (
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type EventType string

const (
	EventClientConnected    EventType = "client.connected"
	EventClientDisconnected EventType = "client.disconnected"
	EventCommandIssued      EventType = "command.issued"
	EventCommandCompleted   EventType = "command.completed"
)

type Event struct {
	Type      EventType   `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}

type RabbitMQEvent struct {
	Delivery amqp091.Delivery
	Event    Event
}

type EventHandler func(event RabbitMQEvent) error
