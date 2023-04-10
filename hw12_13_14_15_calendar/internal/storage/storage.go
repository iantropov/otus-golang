package storage

import "time"

type Storage interface {
	Create(event Event) error
	Update(id EventID, event Event) error
	Delete(id EventID) error
	ListEventForDay(day time.Time) []Event
	ListEventForMonth(monthStart time.Time) []Event
	ListEventForWeek(weekStart time.Time) []Event
}
