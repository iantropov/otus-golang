package scheduler

import (
	"context"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Logger interface {
	Info(string)
	Infof(string, ...any)
	Error(string)
	Errorf(string, ...any)
}

//go:generate mockery --name Application
type Storage interface {
	ListEventBeforeTime(ctx context.Context, before time.Time) []storage.Event
	ListEventCreatedAfter(ctx context.Context, after time.Time) []storage.Event
}

type Scheduler struct {
	logger  Logger
	storage Storage
}

func New(logger Logger, storage Storage) *Scheduler {
	return &Scheduler{
		logger:  logger,
		storage: storage,
	}
}

// #### Уведомление
// Уведомление - временная сущность, в БД не хранится, складывается в очередь для рассыльщика, содержит поля:
// * ID события;
// * Заголовок события;
// * Дата события;
// * Пользователь, которому отправлять.
