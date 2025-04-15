package events

import (
	"encoding/json"
	"fmt"

	"github.com/n0tB0b17/isner/internal/config"
	"github.com/rabbitmq/amqp091-go"
)

type EventConsumer struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	config  *config.Config
	handler map[EventType]EventHandler
}

func NewEventConsumer(config *config.Config) (*EventConsumer, error) {
	conn, err := amqp091.Dial(config.RabbitMQURL)
	if err != nil {
		fmt.Printf("[NewEventConsumer]::> error while creating new consumer: %v \n", err)
		return nil, ErrConnectionFailed
	}

	channel, err := conn.Channel()
	if err != nil {
		fmt.Printf("[NewEventConsumer]::> error while getting channel from from connection: %v \n", err)
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
		fmt.Printf("[NewEventConsumer]::> error while declaring new exchange: %v \n", err)
		return nil, ErrExchangeDeclare
	}

	return &EventConsumer{
		conn:    conn,
		channel: channel,
		config:  config,
		handler: make(map[EventType]EventHandler),
	}, nil
}

func (ec *EventConsumer) StartConsuming(queueName string, routingKey []string) error {
	queue, err := ec.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		fmt.Printf("[events::StartConsuming]::> error while declaring queue: %v \n", err)
		return ErrQueueDeclare
	}

	// binding queue to exchange for routing
	for _, key := range routingKey {
		if err := ec.channel.QueueBind(queue.Name, key, ec.config.EventExchange, false, nil); err != nil {
			fmt.Printf("[event::StartConsuming]::> error while binding queue: %v \n", err)
			return ErrQueueBind
		}
	}

	// start consuming message
	msgs, err := ec.channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Printf("[events::StartConsuming]::> error while consuming: %v \n", err)
		return ErrConsumeFailed
	}

	go func() {
		for msg := range msgs {
			var eve Event
			if err := json.Unmarshal(msg.Body, &eve); err != nil {
				fmt.Printf("[events::StartConsuming]::> error while consuming message inside of channel: %v \n", err)
				msg.Nack(false, true)
				continue
			}

			if handler, exist := ec.handler[eve.Type]; exist {
				rabbitEvent := RabbitMQEvent{
					Delivery: msg,
					Event:    eve,
				}

				if err := handler(rabbitEvent); err != nil {
					fmt.Printf("[events::StartConsuming::handler]::> error while triggering handler with event: %v \n", err)
					msg.Nack(false, true)
					continue
				}

			} else {
				fmt.Printf("[events::StartConsuming]::> no handler registered for eventType: %s", eve.Type)
			}

			msg.Ack(false)
		}
	}()

	fmt.Printf("Starting consuming for queue: %s \n", queueName)
	return nil
}

func (ec *EventConsumer) RegisterEvent(et EventType, eh EventHandler) {
	ec.handler[et] = eh
}

func (ec *EventConsumer) Close() {
	if ec.conn != nil {
		ec.conn.Close()
	}

	if ec.channel != nil {
		ec.channel.Close()
	}
}
