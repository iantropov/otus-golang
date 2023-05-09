package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu                sync.RWMutex
	eventsByIDMap     map[storage.EventID]*storage.Event
	eventsStartsAtMap map[time.Time]storage.EventID
}

var _ storage.Storage = (*Storage)(nil)

var (
	ErrDateBusy       = errors.New("date is already taken")
	ErrEventNotFound  = errors.New("event not found")
	ErrInvalidEvent   = errors.New("invalid event")
	ErrInvalidEventID = errors.New("invalid event ID")
	ErrIDBusy         = errors.New("id is already taken")
)

func New() *Storage {
	eventsByIDMap := make(map[storage.EventID]*storage.Event)
	eventsStartsAtMap := make(map[time.Time]storage.EventID)
	return &Storage{
		eventsByIDMap:     eventsByIDMap,
		eventsStartsAtMap: eventsStartsAtMap,
	}
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}

func (s *Storage) Create(_ context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isValidEvent(event) {
		return ErrInvalidEvent
	}

	if _, exists := s.eventsByIDMap[event.ID]; exists {
		return ErrIDBusy
	}

	if _, exists := s.eventsStartsAtMap[event.StartsAt]; exists {
		return ErrDateBusy
	}

	s.eventsByIDMap[event.ID] = &event
	s.eventsStartsAtMap[event.StartsAt] = event.ID

	return nil
}

func (s *Storage) Get(_ context.Context, id storage.EventID) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, exists := s.eventsByIDMap[id]
	if !exists {
		return storage.Event{}, ErrEventNotFound
	}

	return *event, nil
}

func (s *Storage) Update(_ context.Context, id storage.EventID, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isValidEvent(event) {
		return ErrInvalidEvent
	}

	if id != event.ID {
		return ErrInvalidEventID
	}

	existingEvent, exists := s.eventsByIDMap[id]
	if !exists {
		return ErrEventNotFound
	}

	conflictingEventID, exists := s.eventsStartsAtMap[event.StartsAt]
	if exists && conflictingEventID != id {
		return ErrDateBusy
	}

	delete(s.eventsStartsAtMap, existingEvent.StartsAt)

	s.eventsByIDMap[id] = &event
	s.eventsStartsAtMap[event.StartsAt] = id

	return nil
}

func (s *Storage) Delete(_ context.Context, id storage.EventID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existingEvent, exists := s.eventsByIDMap[id]
	if !exists {
		return ErrEventNotFound
	}

	delete(s.eventsStartsAtMap, existingEvent.StartsAt)
	delete(s.eventsByIDMap, id)

	return nil
}

func (s *Storage) ListEventForDay(_ context.Context, day time.Time) []storage.Event {
	return s.rangeEventsIn(day, day.AddDate(0, 0, 1))
}

func (s *Storage) ListEventForWeek(_ context.Context, weekStart time.Time) []storage.Event {
	return s.rangeEventsIn(weekStart, weekStart.AddDate(0, 0, 7))
}

func (s *Storage) ListEventForMonth(_ context.Context, monthStart time.Time) []storage.Event {
	return s.rangeEventsIn(monthStart, monthStart.AddDate(0, 1, 0))
}

func (s *Storage) ListEventBeforeTime(_ context.Context, before time.Time) []storage.Event {
	return s.rangeEventsWith(func(e *storage.Event) bool {
		return e.EndsAt.Before(before)
	})
}

func (s *Storage) ListEventCreatedAfter(_ context.Context, after time.Time) []storage.Event {
	return s.rangeEventsWith(func(e *storage.Event) bool {
		return e.CreatedAt.After(after)
	})
}

func (s *Storage) isValidEvent(event storage.Event) bool {
	if event.ID == "" || event.Title == "" || event.Description == "" {
		return false
	}

	if !event.StartsAt.Before(event.EndsAt) {
		return false
	}

	if !event.StartsAt.After(time.Now()) {
		return false
	}

	return true
}

func (s *Storage) rangeEventsIn(startTime time.Time, endTime time.Time) []storage.Event {
	return s.rangeEventsWith(func(e *storage.Event) bool {
		return !e.StartsAt.Before(startTime) && e.StartsAt.Before(endTime)
	})
}

func (s *Storage) rangeEventsWith(filter func(*storage.Event) bool) []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]storage.Event, 0)

	for _, event := range s.eventsByIDMap {
		if filter(event) {
			res = append(res, *event)
		}
	}

	return res
}
