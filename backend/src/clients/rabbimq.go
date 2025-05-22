package clients

import (
	"backend/src/config"
	"fmt"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

var (
	rabbitMQConn *amqp091.Connection
	rabbitMQOnce sync.Once
)

// NewRabbitMQClient ensures thread-safe initialization of the RabbitMQ connection.
func NewRabbitMQClient() (*amqp091.Connection, error) {
	var err error
	rabbitMQOnce.Do(func() {
		rabbitMQURL := fmt.Sprintf("amqp://%s:%s@%s/",
			config.GetEnv("RABBITMQ_USER", "cdr-intelligence"),
			config.GetEnv("RABBITMQ_PASSWORD", "cdr-intelligence"),
			config.GetEnv("RABBITMQ_HOST", "localhost:5672"),
		)

		rabbitMQConn, err = amqp091.Dial(rabbitMQURL)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	return rabbitMQConn, nil
}
