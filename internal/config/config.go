package config

type Config struct {
	WebSocketPort string
	RabbitMQURL   string
	EventExchange string
	ClientQueue   string
	CommandQueue  string
}

func GetConfig() *Config {
	return &Config{
		WebSocketPort: "6969",
		RabbitMQURL:   "amqp://0.0.0.0:5672",
		EventExchange: "c2_events",
		ClientQueue:   "client_events",
		CommandQueue:  "command_events",
	}
}
