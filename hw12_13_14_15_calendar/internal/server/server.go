package server

import (
	"context"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

//go:generate mockery --name Logger
type Logger interface {
	Info(string)
	Infof(string, ...any)
	Error(string)
	Errorf(string, ...any)
}

//go:generate mockery --name Application
type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, id storage.EventID, event storage.Event) error
	DeleteEvent(ctx context.Context, id storage.EventID) error
	GetEvent(ctx context.Context, id storage.EventID) (storage.Event, error)
	ListEventForDay(ctx context.Context, day time.Time) []storage.Event
	ListEventForMonth(ctx context.Context, monthStart time.Time) []storage.Event
	ListEventForWeek(ctx context.Context, weekStart time.Time) []storage.Event
}

type Server interface {
	Start() error
	Stop(ctx context.Context) error
}
