package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

var _ storage.Storage = (*Storage)(nil)

var (
	ErrDateBusy       = errors.New("date is already taken")
	ErrEventNotFound  = errors.New("event not found")
	ErrInvalidEvent   = errors.New("invalid event")
	ErrInvalidEventID = errors.New("invalid event ID")
	ErrIDBusy         = errors.New("id is already taken")
)

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
	_, err := s.db.Exec(InsertEvent, event.ID, event.Title, event.StartsAt, event.EndsAt, event.Description, event.UserID, event.NotifyBefore)
	return err
}

func (s *Storage) Get(id storage.EventID) (event storage.Event, err error) {
	err = s.db.QueryRow(SelectEventByID, id).Scan(
		&event.ID,
		&event.Title,
		&event.StartsAt,
		&event.EndsAt,
		&event.Description,
		&event.UserID,
		&event.NotifyBefore,
	)
	return
}

func (s *Storage) Update(id storage.EventID, event storage.Event) error {
	_, err := s.db.Exec(UpdateEvent, event.Title, event.StartsAt, event.EndsAt, event.Description, event.UserID, event.NotifyBefore)
	return err
}

func (s *Storage) Delete(id storage.EventID) error {
	_, err := s.db.Exec(DeleteEvent, id)
	return err
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
