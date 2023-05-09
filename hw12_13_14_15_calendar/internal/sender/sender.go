package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/notifications"
)

//go:generate mockery --name Logger
type Logger interface {
	Info(string)
	Infof(string, ...any)
	Error(string)
	Errorf(string, ...any)
}

//go:generate mockery --name Consumer
type Consumer interface {
	Consume(ctx context.Context) (<-chan []byte, error)
}

type Sender struct {
	logger   Logger
	consumer Consumer
}

func New(logger Logger, consumer Consumer) *Sender {
	return &Sender{
		logger:   logger,
		consumer: consumer,
	}
}

func (s *Sender) Start(ctx context.Context) error {
	messages, err := s.consumer.Consume(ctx)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	var notification notifications.Notification
	for msg := range messages {
		err := json.Unmarshal(msg, &notification)
		if err != nil {
			s.logger.Errorf("Failed to unmarshal notification: %v\n", err)
		}
		s.logger.Infof("SEND NOTIFICATION: %+v", notification)
	}

	return nil
}
