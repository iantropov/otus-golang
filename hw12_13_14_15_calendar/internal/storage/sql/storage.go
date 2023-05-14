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

func (s *Storage) Create(ctx context.Context, event storage.Event) error {
	_, err := s.db.ExecContext(
		ctx,
		InsertEvent,
		event.ID,
		event.Title,
		event.StartsAt,
		event.EndsAt,
		event.CreatedAt,
		event.Description,
		event.UserID,
		event.NotifyBefore.Abs(),
	)
	return err
}

func (s *Storage) Get(ctx context.Context, id storage.EventID) (event storage.Event, err error) {
	err = s.db.QueryRowContext(ctx, SelectEventByID, id).Scan(
		&event.ID,
		&event.Title,
		&event.StartsAt,
		&event.EndsAt,
		&event.CreatedAt,
		&event.Description,
		&event.UserID,
		&event.NotifyBefore,
	)
	return
}

func (s *Storage) Update(ctx context.Context, id storage.EventID, event storage.Event) error {
	_, err := s.db.ExecContext(
		ctx,
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

func (s *Storage) Delete(ctx context.Context, id storage.EventID) error {
	_, err := s.db.ExecContext(ctx, DeleteEvent, id)
	return err
}

func (s *Storage) ListEventForDay(ctx context.Context, day time.Time) []storage.Event {
	return s.rangeEventsIn(ctx, day, day.AddDate(0, 0, 1))
}

func (s *Storage) ListEventForWeek(ctx context.Context, weekStart time.Time) []storage.Event {
	return s.rangeEventsIn(ctx, weekStart, weekStart.AddDate(0, 0, 7))
}

func (s *Storage) ListEventForMonth(ctx context.Context, monthStart time.Time) []storage.Event {
	return s.rangeEventsIn(ctx, monthStart, monthStart.AddDate(0, 1, 0))
}

func (s *Storage) ListEventBeforeTime(ctx context.Context, before time.Time) []storage.Event {
	rows, err := s.db.QueryContext(ctx, SelectEventsBeforeTime, before)
	if err != nil {
		s.logger.Error("failed to query: " + err.Error())
		return nil
	}
	defer rows.Close()
	return s.rangeEvents(rows)
}

func (s *Storage) ListEventCreatedAfter(ctx context.Context, after time.Time) []storage.Event {
	rows, err := s.db.QueryContext(ctx, SelectEventsCreatedAfter, after.UTC())
	if err != nil {
		s.logger.Error("failed to query: " + err.Error())
		return nil
	}
	defer rows.Close()
	return s.rangeEvents(rows)
}

func (s *Storage) rangeEventsIn(ctx context.Context, startTime time.Time, endTime time.Time) []storage.Event {
	rows, err := s.db.QueryContext(ctx, SelectEventsForPeriod, startTime, endTime)
	if err != nil {
		s.logger.Error("failed to query: " + err.Error())
		return nil
	}
	defer rows.Close()
	return s.rangeEvents(rows)
}

func (s *Storage) rangeEvents(rows *sql.Rows) []storage.Event {
	res := make([]storage.Event, 0)

	for rows.Next() {
		var event storage.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.StartsAt,
			&event.EndsAt,
			&event.CreatedAt,
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

	err := rows.Err()
	if err != nil {
		s.logger.Error("failed to iterate: " + err.Error())
		return nil
	}

	return res
}
