package storage

import (
	"context"
	"time"
)

type Storage interface {
	Create(ctx context.Context, event Event) error
	Update(ctx context.Context, id EventID, event Event) error
	Delete(ctx context.Context, id EventID) error
	Get(ctx context.Context, id EventID) (Event, error)
	ListEventForDay(ctx context.Context, day time.Time) []Event
	ListEventForMonth(ctx context.Context, monthStart time.Time) []Event
	ListEventForWeek(ctx context.Context, weekStart time.Time) []Event
}
