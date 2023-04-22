package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Logger interface {
	Info(msg string)
	Infof(f string, args ...any)
	Error(msg string)
	Errorf(f string, args ...any)
}

type Storage struct {
	db     *sql.DB
	logger Logger
}

var _ storage.Storage = (*Storage)(nil)

var (
	ErrDateBusy       = errors.New("date is already taken")
	ErrEventNotFound  = errors.New("event not found")
	ErrInvalidEvent   = errors.New("invalid event")
	ErrInvalidEventID = errors.New("invalid event ID")
	ErrIDBusy         = errors.New("id is already taken")
)

func New(logger Logger, db *sql.DB) *Storage {
	return &Storage{
		db:     db,
		logger: logger,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *Storage) Close(_ context.Context) error {
	s.db.Close()
	return nil
}

func (s *Storage) Create(event storage.Event) error {
	_, err := s.db.Exec(
		InsertEvent,
		event.ID,
		event.Title,
		event.StartsAt,
		event.EndsAt,
		event.Description,
		event.UserID,
		event.NotifyBefore.Abs(),
	)
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
	_, err := s.db.Exec(
		UpdateEvent,
		event.Title,
		event.StartsAt,
		event.EndsAt,
		event.Description,
		event.UserID,
		event.NotifyBefore,
		id,
	)
	return err
}

func (s *Storage) Delete(id storage.EventID) error {
	_, err := s.db.Exec(DeleteEvent, id)
	return err
}

func (s *Storage) ListEventForDay(day time.Time) []storage.Event {
	return s.rangeEvents(day, day.AddDate(0, 0, 1))
}

func (s *Storage) ListEventForWeek(weekStart time.Time) []storage.Event {
	return s.rangeEvents(weekStart, weekStart.AddDate(0, 0, 7))
}

func (s *Storage) ListEventForMonth(monthStart time.Time) []storage.Event {
	return s.rangeEvents(monthStart, monthStart.AddDate(0, 1, 0))
}

func (s *Storage) rangeEvents(startTime time.Time, endTime time.Time) []storage.Event {
	res := make([]storage.Event, 0)

	rows, err := s.db.Query(SelectEventsForPeriod, startTime, endTime)
	if err != nil {
		s.logger.Error("failed to query: " + err.Error())
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var event storage.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.StartsAt,
			&event.EndsAt,
			&event.Description,
			&event.UserID,
			&event.NotifyBefore,
		)
		if err != nil {
			s.logger.Error("failed to scan: " + err.Error())
			return nil
		}
		res = append(res, event)
	}

	err = rows.Err()
	if err != nil {
		s.logger.Error("failed to iterate: " + err.Error())
		return nil
	}

	return res
}
