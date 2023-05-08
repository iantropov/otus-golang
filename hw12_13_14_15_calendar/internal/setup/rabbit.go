package setup

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/rabbit"
)

func SetupRabbit(config config.RabbitConf) (*rabbit.Connection, error) {
	conn, err := amqp.Dial(config.DSN)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("Failed to open a channel: %w", err)
	}

	queue, err := ch.QueueDeclare(
		config.Queue, // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("Failed to declare a queue: %w", err)
	}

	return rabbit.NewConnection(conn, ch, queue), nil
}
