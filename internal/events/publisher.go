package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/n0tB0b17/isner/internal/config"
	"github.com/rabbitmq/amqp091-go"
)

type EventPublisher struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	config  *config.Config
}

func NewEventPublisher(config *config.Config) (*EventPublisher, error) {
	conn, err := amqp091.Dial(config.RabbitMQURL)
	if err != nil {
		fmt.Printf("[NewEventPublisher]::> unable to dial rabbitMQ server: %v \n", err)
		return nil, ErrConnectionFailed
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		fmt.Printf("[NewEventPublisher]::> unable to access channel from connection: %v \n", err)
		return nil, ErrChannelCreation
	}

	err = channel.ExchangeDeclare(
		config.EventExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		conn.Close()
		channel.Close()
		fmt.Printf("[NewEventPublisher]::> unable to delete new channel-exhange: %v \n", err)
		return nil, ErrExchangeDeclare
	}

	return &EventPublisher{
		conn:    conn,
		channel: channel,
		config:  config,
	}, nil
}

func (ep *EventPublisher) Publish(eventType EventType, payload interface{}) error {
	eve := Event{
		Type:      eventType,
		Payload:   payload,
		Timestamp: time.Now(),
	}

	docs, err := json.Marshal(eve)
	if err != nil {
		fmt.Printf("[events::Publish]::> unable to encode event to json: %v \n", err)
		return err
	}

	err = ep.channel.Publish(
		ep.config.EventExchange,
		string(eventType),
		false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        docs,
		},
	)

	if err != nil {
		fmt.Printf("[events::Publish]::> unable to publish event to channel: %v \n", err)
		return ErrPublishFailed
	}

	return err
}

func (ep *EventPublisher) Close() {
	if ep.channel != nil {
		ep.channel.Close()
	}

	if ep.conn != nil {
		ep.conn.Close()
	}
}
