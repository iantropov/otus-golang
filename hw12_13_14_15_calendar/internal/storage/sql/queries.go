package sqlstorage

import (
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

func (s *Storage) Create(event storage.Event) error {
	return nil
}

func (s *Storage) Get(id storage.EventID) (storage.Event, error) {
	return storage.Event{}, nil
}

func (s *Storage) Update(id storage.EventID, event storage.Event) error {
	return nil
}

func (s *Storage) Delete(id storage.EventID) error {
	return nil
}

func (s *Storage) ListEventForDay(day time.Time) []storage.Event {
	return nil
}

func (s *Storage) ListEventForWeek(weekStart time.Time) []storage.Event {
	return nil
}

func (s *Storage) ListEventForMonth(monthStart time.Time) []storage.Event {
	return nil
}
