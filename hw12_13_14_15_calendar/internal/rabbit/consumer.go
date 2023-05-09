package rabbit

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	logger Logger
	conn   *Connection
}

func NewConsumer(logger Logger, conn *Connection) *Consumer {
	return &Consumer{
		logger: logger,
		conn:   conn,
	}
}

func (c *Consumer) Consume(ctx context.Context, message []byte) (<-chan []byte, error) {
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

			select {
			case <-ctx.Done():
				return
			case out <- message.Body:
			}

			if !ok {
				return
			}
		}
	}()

	return out, nil
}
