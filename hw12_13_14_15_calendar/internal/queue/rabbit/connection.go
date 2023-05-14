package rabbit

import (
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/queue"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

func NewConnection(conn *amqp.Connection, ch *amqp.Channel, queue amqp.Queue) queue.Connection {
	return &Connection{
		conn:  conn,
		ch:    ch,
		queue: queue,
	}
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
