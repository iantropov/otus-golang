package app

import (
	"context"
	"errors"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage storage.Storage
}

type Logger interface {
	Info(string)
	Infof(string, ...any)
	Error(string)
	Errorf(string, ...any)
}

func New(logger Logger, storage storage.Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	return wrapError(a.storage.Create(ctx, event))
}

func (a *App) UpdateEvent(ctx context.Context, id storage.EventID, event storage.Event) error {
	return wrapError(a.storage.Update(ctx, id, event))
}

func (a *App) GetEvent(ctx context.Context, id storage.EventID) (storage.Event, error) {
	event, err := a.storage.Get(ctx, id)
	return event, wrapError(err)
}

func (a *App) DeleteEvent(ctx context.Context, id storage.EventID) error {
	return wrapError(a.storage.Delete(ctx, id))
}

func (a *App) ListEventForDay(ctx context.Context, at time.Time) []storage.Event {
	return a.storage.ListEventForDay(ctx, at)
}

func (a *App) ListEventForWeek(ctx context.Context, at time.Time) []storage.Event {
	return a.storage.ListEventForWeek(ctx, at)
}

func (a *App) ListEventForMonth(ctx context.Context, at time.Time) []storage.Event {
	return a.storage.ListEventForMonth(ctx, at)
}

func wrapError(err error) error {
	var internalStorageErr storage.InternalError
	if errors.As(err, &internalStorageErr) {
		return InternalError{Err: internalStorageErr}
	}
	return err
}
