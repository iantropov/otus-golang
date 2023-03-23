package storage

import "time"

type Storage interface {
	Create(event Event) error
	Update(id EventId, event Event) error
	Delete(id EventId) error
	ListEventForDay(day time.Time) []Event
	ListEventForMonth(monthStart time.Time) []Event
	ListEventForWeek(weekStart time.Time) []Event
}
