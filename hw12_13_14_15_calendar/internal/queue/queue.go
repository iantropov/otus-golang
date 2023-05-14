package queue

import (
	"context"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/pkg/logger"
)

type Connection interface {
	Close()
	NewProducer(logg *logger.Logger) Producer
	NewConsumer(logg *logger.Logger) Consumer
}

type Producer interface {
	Produce(ctx context.Context, message []byte) error
}

type Consumer interface {
	Consume(ctx context.Context) (<-chan []byte, error)
}
