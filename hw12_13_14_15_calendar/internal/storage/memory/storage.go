package memorystorage

import (
	"errors"
	"sync"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu                sync.RWMutex
	eventsByIdMap     map[storage.EventId]storage.Event
	eventsStartsAtMap map[time.Time]storage.EventId
}

var _ storage.Storage = (*Storage)(nil)

var (
	ErrDateBusy       = errors.New("date is already taken")
	ErrEventNotFound  = errors.New("event not found")
	ErrInvalidEvent   = errors.New("invalid event")
	ErrInvalidEventId = errors.New("invalid event id")
	ErrIdBusy         = errors.New("id is already taken")
)

func New() *Storage {
	eventsByIdMap := make(map[storage.EventId]storage.Event)
	eventsStartsAtMap := make(map[time.Time]storage.EventId)
	return &Storage{
		eventsByIdMap:     eventsByIdMap,
		eventsStartsAtMap: eventsStartsAtMap,
	}
}

func (s *Storage) Create(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isValidEvent(event) {
		return ErrInvalidEvent
	}

	if _, exists := s.eventsByIdMap[event.Id]; exists {
		return ErrIdBusy
	}

	if _, exists := s.eventsStartsAtMap[event.StartsAt]; exists {
		return ErrDateBusy
	}

	s.eventsByIdMap[event.Id] = event
	s.eventsStartsAtMap[event.StartsAt] = event.Id

	return nil
}

func (s *Storage) Get(id storage.EventId) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, exists := s.eventsByIdMap[id]
	if !exists {
		return storage.Event{}, ErrEventNotFound
	}

	return event, nil
}

func (s *Storage) Update(id storage.EventId, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isValidEvent(event) {
		return ErrInvalidEvent
	}

	if id != event.Id {
		return ErrInvalidEventId
	}

	existingEvent, exists := s.eventsByIdMap[id]
	if !exists {
		return ErrEventNotFound
	}

	conflictingEventId, exists := s.eventsStartsAtMap[event.StartsAt]
	if exists && conflictingEventId != id {
		return ErrDateBusy
	}

	delete(s.eventsStartsAtMap, existingEvent.StartsAt)

	s.eventsByIdMap[id] = event
	s.eventsStartsAtMap[event.StartsAt] = id

	return nil
}

func (s *Storage) Delete(id storage.EventId) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existingEvent, exists := s.eventsByIdMap[id]
	if !exists {
		return ErrEventNotFound
	}

	delete(s.eventsStartsAtMap, existingEvent.StartsAt)
	delete(s.eventsByIdMap, id)

	return nil
}

func (s *Storage) ListEventForDay(day time.Time) []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.rangeEvents(day, day.AddDate(0, 0, 1))
}

func (s *Storage) ListEventForWeek(weekStart time.Time) []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.rangeEvents(weekStart, weekStart.AddDate(0, 0, 7))
}

func (s *Storage) ListEventForMonth(monthStart time.Time) []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.rangeEvents(monthStart, monthStart.AddDate(0, 1, 0))
}

func (s *Storage) isValidEvent(event storage.Event) bool {
	if event.Id == "" || event.Title == "" || event.Description == "" {
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
	res := make([]storage.Event, 0)

	for _, event := range s.eventsByIdMap {
		if !event.StartsAt.Before(startTime) && event.StartsAt.Before(endTime) {
			res = append(res, event)
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
