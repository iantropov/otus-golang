package rabbit

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	logger Logger
	conn   *Connection
}

func NewProducer(logger Logger, conn *Connection) *Producer {
	return &Producer{
		logger: logger,
		conn:   conn,
	}
}

func (p *Producer) Produce(ctx context.Context, message []byte) error {
	return p.conn.ch.PublishWithContext(ctx,
		"",                // exchange
		p.conn.queue.Name, // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
}
