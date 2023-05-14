package rabbit

import (
	"fmt"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/queue"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

func NewConnection(config config.QueueConf) (queue.Connection, error) {
	conn, err := amqp.Dial(config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	queue, err := ch.QueueDeclare(
		config.Name, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &Connection{
		conn:  conn,
		ch:    ch,
		queue: queue,
	}, nil
}

func (conn *Connection) Close() {
	conn.ch.Close()
	conn.conn.Close()
}

func (conn *Connection) NewConsumer(logger *logger.Logger) queue.Consumer {
	return &Consumer{
		logger: logger,
		conn:   conn,
	}
}

func (conn *Connection) NewProducer(logger *logger.Logger) queue.Producer {
	return &Producer{
		logger: logger,
		conn:   conn,
	}
}
