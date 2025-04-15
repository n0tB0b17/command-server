package events

import "errors"

var (
	ErrConnectionFailed = errors.New("failed to connect to RabbitMQ")
	ErrChannelCreation  = errors.New("failed to create channel")
	ErrExchangeDeclare  = errors.New("failed to declare exchange")
	ErrQueueDeclare     = errors.New("failed to declare queue")
	ErrQueueBind        = errors.New("failed to bind queue")
	ErrPublishFailed    = errors.New("failed to publish message")
	ErrConsumeFailed    = errors.New("failed to consume messages")
)
