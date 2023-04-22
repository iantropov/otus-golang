package sqlstorage

import (
	"context"
	"database/sql"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

var _ storage.Storage = (*Storage)(nil)

func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	return s.db.Ping()
}

func (s *Storage) Close(ctx context.Context) error {
	s.db.Close()
	return nil
}

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
