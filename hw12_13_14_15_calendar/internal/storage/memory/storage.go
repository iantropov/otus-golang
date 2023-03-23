package memorystorage

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events []storage.Event
}

func New() *Storage {
	fmt.Println("Started in-memory storage!")
	return &Storage{}
}

func (s *Storage) Create(event storage.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = append(s.events, event)
	s.sortEvents()
}

func (s *Storage) Update(id storage.EventId, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	eventIdx := s.findEventIndex(id)
	if eventIdx == -1 {
		return fmt.Errorf("event: %v not found", id)
	}

	s.events[eventIdx] = event
	s.sortEvents()

	return nil
}

func (s *Storage) Delete(id storage.EventId) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	eventIdx := s.findEventIndex(id)
	if eventIdx == -1 {
		return fmt.Errorf("event: %v not found", id)
	}

	s.events = append(s.events[:eventIdx], s.events[eventIdx+1:]...)

	return nil
}

func (s *Storage) ListEventForDay(day time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return selectEventsForRange(startTime, endTime)
}

func (s *Storage) ListEventForWeek(weekStart time.Time) {
}

func (s *Storage) ListEventForMonth(monthStart time.Time) {
}

func (s *Storage) sortEvents() {
	sort.Slice(s.events, func(i, j int) bool {
		return s.events[i].StartsAt.Before(s.events[j].StartsAt)
	})
}

func (s *Storage) findEventIndex(id storage.EventId) int {
	eventIdx := -1
	for i := range s.events {
		if s.events[i].Id == id {
			eventIdx = i
			break
		}
	}
	return eventIdx
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
