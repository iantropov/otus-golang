package rabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

func NewConnection(conn *amqp.Connection, ch *amqp.Channel, queue amqp.Queue) *Connection {
	return &Connection{
		conn:  conn,
		ch:    ch,
		queue: queue,
	}
}

func (r *Connection) Close() {
	r.ch.Close()
	r.conn.Close()
}
