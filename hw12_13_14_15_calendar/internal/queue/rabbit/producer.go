package rabbit

import (
	"context"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	logger *logger.Logger
	conn   *Connection
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
