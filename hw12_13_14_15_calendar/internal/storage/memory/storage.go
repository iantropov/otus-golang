package memorystorage

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu           sync.RWMutex
	sortedEvents []*storage.Event
	eventsMap    map[storage.EventId]storage.Event
}

var _ storage.Storage = (*Storage)(nil)

var (
	ErrDateBusy      = errors.New("date is already taken")
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidEvent  = errors.New("invalid event")
	ErrIdBudy        = errors.New("id is already taken")
)

func New() *Storage {
	eventsMap := make(map[storage.EventId]int)
	return &Storage{
		eventsMap: eventsMap,
	}
}

func (s *Storage) Create(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isValidEvent(event) {
		return ErrInvalidEvent
	}

	if _, exists := s.eventsMap[event.Id]; exists {
		return ErrIdBudy
	}

	if s.findEventIndexByDate(event.StartsAt) != -1 {
		return ErrDateBusy
	}

	s.events = append(s.events, event)
	s.sortEvents()

	return nil
}

func (s *Storage) Update(id storage.EventId, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	updatingEventIdx := s.findEventIndexById(id)
	if updatingEventIdx == -1 {
		return ErrEventNotFound
	}

	conflictingEventIdx := s.findEventIndexByDate(event.StartsAt)
	if conflictingEventIdx != -1 {
		return ErrDateBusy
	}

	s.events[updatingEventIdx] = event
	s.sortEvents()

	return nil
}

func (s *Storage) Delete(id storage.EventId) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.eventsMap[event.Id]; exists {
		return ErrIdBudy
	}

	eventIdx := s.findEventIndexById(id)
	if eventIdx == -1 {
		return ErrEventNotFound
	}

	s.events = append(s.events[:eventIdx], s.events[eventIdx+1:]...)

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

func (s *Storage) findEventIndexByDate(date time.Time) int {
	for i := range s.events {
		if s.events[i].StartsAt == date {
			return i
		}
	}
	return -1
}

func (s *Storage) sortEvents() {
	sort.Slice(s.events, func(i, j int) bool {
		return s.events[i].StartsAt.Before(s.events[j].StartsAt)
	})
}

func (s *Storage) findEventIndexById(id storage.EventId) int {
	for i := range s.events {
		if s.events[i].Id == id {
			return i
		}
	}
	return -1
}

func (s *Storage) rangeEvents(startTime time.Time, endTime time.Time) []storage.Event {
	startIdx, endIdx := -1, -1

	for i := range s.events {
		if startIdx == -1 && !s.events[i].StartsAt.Before(startTime) {
			startIdx = i
		}
		if startIdx > -1 && s.events[i].StartsAt.After(endTime) {
			endIdx = i - 1
			break
		}
	}

	if startIdx == -1 {
		return nil
	}

	if endIdx == -1 {
		endIdx = len(s.events)
	}

	eventsCopy := make([]storage.Event, endIdx-startIdx+1)
	copy(eventsCopy, s.events[startIdx:endIdx])
	return eventsCopy
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
