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

func TestStorageCreate(t *testing.T) {
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
			title: "empty event",
			event: storage.Event{},
			err:   ErrInvalidEvent,
		},
		{
			title: "invalid event (outdated)",
			event: buildEventWith(map[string]any{
				"StartsAt": date(2013, 6, 2),
				"EndsAt":   date(2013, 6, 2),
			}),
			err: ErrInvalidEvent,
		},
		{
			title: "invalid event (starts_at > ends_at)",
			event: buildEventWith(map[string]any{
				"StartsAt": date(2033, 6, 1),
				"EndsAt":   date(2032, 6, 2),
			}),
			err: ErrInvalidEvent,
		},
		{
			title: "invalid event (without id)",
			event: buildEventWith(map[string]any{
				"Id": storage.EventId(""),
			}),
			err: ErrInvalidEvent,
		},
		{
			title: "invalid event (without title)",
			event: buildEventWith(map[string]any{
				"Title": "",
			}),
			err: ErrInvalidEvent,
		},
		{
			title: "invalid event (without description)",
			event: buildEventWith(map[string]any{
				"Description": "",
			}),
			err: ErrInvalidEvent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			memStorage := New()
			err := memStorage.Create(tt.event)
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
