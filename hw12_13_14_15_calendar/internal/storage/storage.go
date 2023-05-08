package storage

import (
	"context"
	"fmt"
	"time"
)

type InternalError struct {
	err error
}

func (e InternalError) Error() string {
	return fmt.Sprintf("internal storage error: %v", e.err)
}

type Storage interface {
	Close(ctx context.Context) error
	Create(ctx context.Context, event Event) error
	Update(ctx context.Context, id EventID, event Event) error
	Delete(ctx context.Context, id EventID) error
	Get(ctx context.Context, id EventID) (Event, error)
	ListEventForDay(ctx context.Context, day time.Time) []Event
	ListEventForMonth(ctx context.Context, monthStart time.Time) []Event
	ListEventForWeek(ctx context.Context, weekStart time.Time) []Event
	ListEventBeforeTime(ctx context.Context, before time.Time) []Event
	ListEventsCreatedAfter(ctx context.Context, after time.Time) []Event
}

// ### Описание методов
// * Создать (событие);
// * Обновить (ID события, событие);
// * Удалить (ID события);
// * СписокСобытийНаДень (дата);
// * СписокСобытийНаНеделю (дата начала недели);
// * СписокСобытийНaМесяц (дата начала месяца).
