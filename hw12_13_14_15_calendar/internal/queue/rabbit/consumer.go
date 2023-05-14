package rabbit

import (
	"context"
	"fmt"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/queue"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	logger *logger.Logger
	conn   *Connection
}

func NewConsumer(logger *logger.Logger, conn *Connection) queue.Consumer {
	return &Consumer{
		logger: logger,
		conn:   conn,
	}
}

func (c *Consumer) Consume(ctx context.Context) (<-chan []byte, error) {
	msgs, err := c.conn.ch.Consume(
		c.conn.queue.Name, // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to regisger a consumer: %w", err)
	}

	out := make(chan []byte)
	go func() {
		defer close(out)

		var (
			message amqp.Delivery
			ok      bool
		)
		for {
			select {
			case <-ctx.Done():
				return
			case message, ok = <-msgs:
			}

			if !ok {
				return
			}

			select {
			case <-ctx.Done():
				return
			case out <- message.Body:
			}
		}
	}()

	return out, nil
}
