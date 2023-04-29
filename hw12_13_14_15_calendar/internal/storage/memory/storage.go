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
	return s.rangeEvents(day, day.AddDate(0, 0, 1))
}

func (s *Storage) ListEventForWeek(_ context.Context, weekStart time.Time) []storage.Event {
	return s.rangeEvents(weekStart, weekStart.AddDate(0, 0, 7))
}

func (s *Storage) ListEventForMonth(_ context.Context, monthStart time.Time) []storage.Event {
	return s.rangeEvents(monthStart, monthStart.AddDate(0, 1, 0))
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

func (s *Storage) rangeEvents(startTime time.Time, endTime time.Time) []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]storage.Event, 0)

	for _, event := range s.eventsByIDMap {
		if !event.StartsAt.Before(startTime) && event.StartsAt.Before(endTime) {
			res = append(res, *event)
		}
	}

	return res
}

// TODO

// Событие - основная сущность, содержит в себе поля:
// * ID - уникальный идентификатор события (можно воспользоваться UUID);
// * Заголовок - короткий текст;
// * Дата и время события;
// * Длительность события (или дата и время окончания);
// * Описание события - длинный текст, опционально;
// * ID пользователя, владельца события;
// * За сколько времени высылать уведомление, опционально.

// #### Уведомление
// Уведомление - временная сущность, в БД не хранится, складывается в очередь для рассыльщика, содержит поля:
// * ID события;
// * Заголовок события;
// * Дата события;
// * Пользователь, которому отправлять.

// ### Описание методов
// * Создать (событие);
// * Обновить (ID события, событие);
// * Удалить (ID события);
// * СписокСобытийНаДень (дата);
// * СписокСобытийНаНеделю (дата начала недели);
// * СписокСобытийНaМесяц (дата начала месяца).
