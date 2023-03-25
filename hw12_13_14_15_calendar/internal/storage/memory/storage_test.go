package memorystorage

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorageBusinessLogic(t *testing.T) {
	memStorage := New()

	event := buildEvent()

	err := memStorage.Create(event)
	require.Equal(t, nil, err)

	resEvent, err := memStorage.Get(event.Id)
	require.Equal(t, event, resEvent)
	require.Equal(t, nil, err)

	event.UserId = gofakeit.Email()
	err = memStorage.Update(event.Id, event)
	require.Equal(t, nil, err)

	resEvent, err = memStorage.Get(event.Id)
	require.Equal(t, event, resEvent)
	require.Equal(t, nil, err)

	resEvents := memStorage.ListEventForDay(date(2033, 6, 1))
	require.Equal(t, []storage.Event{event}, resEvents)

	err = memStorage.Delete(event.Id)
	require.Equal(t, nil, err)

	_, err = memStorage.Get(event.Id)
	require.Equal(t, ErrEventNotFound, err)
}

func TestStorageValidateEvent(t *testing.T) {
	tests := []struct {
		title   string
		event   storage.Event
		isValid bool
	}{
		{
			title:   "valid event",
			event:   buildEvent(),
			isValid: true,
		},
		{
			title:   "empty event",
			event:   storage.Event{},
			isValid: false,
		},
		{
			title: "invalid event (outdated)",
			event: buildEventWith(map[string]any{
				"StartsAt": date(2013, 6, 2),
				"EndsAt":   date(2013, 6, 2),
			}),
			isValid: false,
		},
		{
			title: "invalid event (starts_at > ends_at)",
			event: buildEventWith(map[string]any{
				"StartsAt": date(2033, 6, 1),
				"EndsAt":   date(2032, 6, 2),
			}),
			isValid: false,
		},
		{
			title: "invalid event (without id)",
			event: buildEventWith(map[string]any{
				"Id": storage.EventId(""),
			}),
			isValid: false,
		},
		{
			title: "invalid event (without title)",
			event: buildEventWith(map[string]any{
				"Title": "",
			}),
			isValid: false,
		},
		{
			title: "invalid event (without description)",
			event: buildEventWith(map[string]any{
				"Description": "",
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			memStorage := New()
			isValid := memStorage.isValidEvent(tt.event)
			require.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestStorageCreate(t *testing.T) {
	existingEvent := buildEventWith(map[string]any{
		"StartsAt": date(2035, 6, 2),
		"EndsAt":   date(2035, 6, 3),
	})

	tests := []struct {
		title string
		event storage.Event
		err   error
	}{
		{
			title: "valid event",
			event: buildEvent(),
			err:   nil,
		},
		{
			title: "invalid event",
			event: storage.Event{},
			err:   ErrInvalidEvent,
		},
		{
			title: "duplicated id",
			event: buildEventWith(map[string]any{
				"Id": existingEvent.Id,
			}),
			err: ErrIdBusy,
		},
		{
			title: "duplicated StartsAt",
			event: buildEventWith(map[string]any{
				"StartsAt": existingEvent.StartsAt,
				"EndsAt":   date(2035, 7, 2),
			}),
			err: ErrDateBusy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			memStorage := New()

			err := memStorage.Create(existingEvent)
			require.Equal(t, nil, err)

			err = memStorage.Create(tt.event)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestStorageUpdate(t *testing.T) {
	existingEvent := buildEventWith(map[string]any{
		"StartsAt": date(2035, 6, 2),
		"EndsAt":   date(2035, 6, 3),
	})

	existingEvent2 := buildEventWith(map[string]any{
		"StartsAt": date(2045, 6, 2),
		"EndsAt":   date(2045, 6, 3),
	})

	tests := []struct {
		title   string
		eventId storage.EventId
		event   storage.Event
		err     error
	}{
		{
			title:   "positive case",
			eventId: existingEvent.Id,
			event: buildEventWith(map[string]any{
				"Id": existingEvent.Id,
			}),
			err: nil,
		},
		{
			title:   "invalid event",
			eventId: existingEvent.Id,
			event:   storage.Event{},
			err:     ErrInvalidEvent,
		},
		{
			title:   "invalid event id",
			eventId: existingEvent.Id,
			event:   buildEvent(),
			err:     ErrInvalidEventId,
		},
		{
			title:   "event not found",
			eventId: "123",
			event: buildEventWith(map[string]any{
				"Id": storage.EventId("123"),
			}),
			err: ErrEventNotFound,
		},
		{
			title:   "duplicated StartsAt",
			eventId: existingEvent.Id,
			event: buildEventWith(map[string]any{
				"Id":       existingEvent.Id,
				"StartsAt": existingEvent2.StartsAt,
				"EndsAt":   existingEvent2.StartsAt.Add(time.Hour),
			}),
			err: ErrDateBusy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			memStorage := New()

			err := memStorage.Create(existingEvent)
			require.Equal(t, nil, err)

			err = memStorage.Create(existingEvent2)
			require.Equal(t, nil, err)

			err = memStorage.Update(tt.eventId, tt.event)
			require.Equal(t, tt.err, err)
		})
	}
}

func buildEvent() storage.Event {
	return storage.Event{
		StartsAt:     date(2033, 6, 1),
		EndsAt:       date(2033, 6, 2),
		Title:        gofakeit.FarmAnimal(),
		Description:  gofakeit.Adjective(),
		Id:           storage.EventId(gofakeit.UUID()),
		UserId:       gofakeit.Email(),
		NotifyBefore: time.Second * 2,
	}
}

func buildEventWith(attrs map[string]any) storage.Event {
	event := buildEvent()

	for key := range attrs {
		switch key {
		case "Id":
			event.Id = attrs[key].(storage.EventId)
		case "StartsAt":
			event.StartsAt = attrs[key].(time.Time)
		case "EndsAt":
			event.EndsAt = attrs[key].(time.Time)
		case "Title":
			event.Title = attrs[key].(string)
		case "Description":
			event.Description = attrs[key].(string)
		default:
			panic("invalid key")
		}
	}

	return event
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}
