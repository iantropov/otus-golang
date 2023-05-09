package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/notifications"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

//go:generate mockery --name Logger
type Logger interface {
	Info(string)
	Infof(string, ...any)
	Error(string)
	Errorf(string, ...any)
}

//go:generate mockery --name Storage
type Storage interface {
	ListEventBeforeTime(ctx context.Context, before time.Time) []storage.Event
	ListEventCreatedAfter(ctx context.Context, after time.Time) []storage.Event
	Delete(ctx context.Context, id storage.EventID) error
}

//go:generate mockery --name Storage
type Producer interface {
	Produce(ctx context.Context, message []byte) error
}

type Scheduler struct {
	logger          Logger
	storage         Storage
	producer        Producer
	period          time.Duration
	lastScheduledAt time.Time
}

func New(logger Logger, storage Storage, producer Producer, period time.Duration) *Scheduler {
	return &Scheduler{
		logger:          logger,
		storage:         storage,
		producer:        producer,
		period:          period,
		lastScheduledAt: time.Now(),
	}
}

func (s *Scheduler) Schedule(ctx context.Context) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(s.period):
				s.cleanupEvents(ctx)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(s.period):
				s.scheduleEvents(ctx)
			}
		}
	}()

	s.logger.Info("Scheduling events...")

	wg.Wait()
}

func (s *Scheduler) cleanupEvents(ctx context.Context) {
	events := s.storage.ListEventBeforeTime(ctx, time.Now().AddDate(-1, 0, 0))
	for i := range events {
		err := s.storage.Delete(ctx, events[i].ID)
		if err != nil {
			s.logger.Errorf("Failed to delete event %s: %\n", events[i].ID, err)
			continue
		}
		s.logger.Infof("Successfully deleted event %s\n", events[i].ID)
	}
}

func (s *Scheduler) scheduleEvents(ctx context.Context) {
	events := s.storage.ListEventCreatedAfter(ctx, s.lastScheduledAt)
	for i := range events {
		if events[i].StartsAt.Add(-events[i].NotifyBefore).After(time.Now()) {
			err := s.produceEvent(ctx, events[i])
			if err != nil {
				s.logger.Errorf("Failed to schedule event %s: %\n", events[i].ID, err)
				continue
			}
			s.logger.Infof("Successfully scheduled event %s\n", events[i].ID)
		}
		s.lastScheduledAt = events[i].CreatedAt
	}
}

func (s *Scheduler) produceEvent(ctx context.Context, event storage.Event) error {
	bytes, err := json.Marshal(notifications.Notification{
		ID:       event.ID,
		Title:    event.Title,
		StartsAt: event.StartsAt,
		UserID:   event.UserID,
	})
	if err != nil {
		return fmt.Errorf("failed to marsal notification: %w", err)
	}

	err = s.producer.Produce(ctx, bytes)
	if err != nil {
		return fmt.Errorf("failed to produce notification: %w", err)
	}

	return nil
}
