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

	resEvent, err := memStorage.Get(event.ID)
	require.Equal(t, event, resEvent)
	require.Equal(t, nil, err)

	event.UserID = gofakeit.Email()
	err = memStorage.Update(event.ID, event)
	require.Equal(t, nil, err)

	resEvent, err = memStorage.Get(event.ID)
	require.Equal(t, event, resEvent)
	require.Equal(t, nil, err)

	resEvents := memStorage.ListEventForDay(date(2033, 6, 1))
	require.Equal(t, []storage.Event{event}, resEvents)

	err = memStorage.Delete(event.ID)
	require.Equal(t, nil, err)

	_, err = memStorage.Get(event.ID)
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
				"ID": storage.EventID(""),
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
				"ID": existingEvent.ID,
			}),
			err: ErrIDBusy,
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

func TestStorageGet(t *testing.T) {
	existingEvent := buildEventWith(map[string]any{
		"StartsAt": date(2035, 6, 2),
		"EndsAt":   date(2035, 6, 3),
	})

	tests := []struct {
		title   string
		eventID storage.EventID
		event   storage.Event
		err     error
	}{
		{
			title:   "positive case",
			eventID: existingEvent.ID,
			event:   existingEvent,
		},
		{
			title:   "event not found",
			eventID: "123",
			err:     ErrEventNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			memStorage := New()

			err := memStorage.Create(existingEvent)
			require.Equal(t, nil, err)

			event, err := memStorage.Get(tt.eventID)
			require.Equal(t, tt.event, event)
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
		eventID storage.EventID
		event   storage.Event
		err     error
	}{
		{
			title:   "positive case",
			eventID: existingEvent.ID,
			event: buildEventWith(map[string]any{
				"ID": existingEvent.ID,
			}),
			err: nil,
		},
		{
			title:   "invalid event",
			eventID: existingEvent.ID,
			event:   storage.Event{},
			err:     ErrInvalidEvent,
		},
		{
			title:   "invalid event id",
			eventID: existingEvent.ID,
			event:   buildEvent(),
			err:     ErrInvalidEventID,
		},
		{
			title:   "event not found",
			eventID: "123",
			event: buildEventWith(map[string]any{
				"ID": storage.EventID("123"),
			}),
			err: ErrEventNotFound,
		},
		{
			title:   "duplicated StartsAt",
			eventID: existingEvent.ID,
			event: buildEventWith(map[string]any{
				"ID":       existingEvent.ID,
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

			err = memStorage.Update(tt.eventID, tt.event)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestStorageDelete(t *testing.T) {
	existingEvent := buildEventWith(map[string]any{
		"StartsAt": date(2035, 6, 2),
		"EndsAt":   date(2035, 6, 3),
	})

	tests := []struct {
		title   string
		eventID storage.EventID
		err     error
	}{
		{
			title:   "positive case",
			eventID: existingEvent.ID,
			err:     nil,
		},
		{
			title:   "event not found",
			eventID: "123",
			err:     ErrEventNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			memStorage := New()

			err := memStorage.Create(existingEvent)
			require.Equal(t, nil, err)

			err = memStorage.Delete(tt.eventID)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestStorageListEventForDay(t *testing.T) {
	date := date(2050, 1, 1)
	events := []storage.Event{
		buildEventWith(map[string]any{
			"StartsAt": date.Add(time.Minute),
			"EndsAt":   date.Add(time.Minute * 2),
		}),
		buildEventWith(map[string]any{
			"StartsAt": date.Add(time.Minute * 3),
			"EndsAt":   date.Add(time.Minute * 4),
		}),
		buildEventWith(map[string]any{
			"StartsAt": date.AddDate(0, 0, 1),
			"EndsAt":   date.AddDate(0, 0, 2),
		}),
	}

	memStorage := New()
	for i := range events {
		err := memStorage.Create(events[i])
		require.Equal(t, nil, err)
	}

	dayEvents := memStorage.ListEventForDay(date)
	require.ElementsMatch(t, dayEvents, []storage.Event{events[0], events[1]})
}

func TestStorageListEventForWeek(t *testing.T) {
	date := date(2050, 1, 1)
	events := []storage.Event{
		buildEventWith(map[string]any{
			"StartsAt": date.AddDate(0, 0, 1),
			"EndsAt":   date.AddDate(0, 0, 2),
		}),
		buildEventWith(map[string]any{
			"StartsAt": date.AddDate(0, 0, 3),
			"EndsAt":   date.AddDate(0, 0, 4),
		}),
		buildEventWith(map[string]any{
			"StartsAt": date.AddDate(0, 0, 5),
			"EndsAt":   date.AddDate(0, 0, 15),
		}),
		buildEventWith(map[string]any{
			"StartsAt": date.AddDate(0, 0, 25),
			"EndsAt":   date.AddDate(0, 0, 35),
		}),
	}

	memStorage := New()
	for i := range events {
		err := memStorage.Create(events[i])
		require.Equal(t, nil, err)
	}

	weekEvents := memStorage.ListEventForWeek(date)
	require.ElementsMatch(t, weekEvents, []storage.Event{events[0], events[1], events[2]})
}

func TestStorageListEventFoMonth(t *testing.T) {
	date := date(2050, 1, 1)
	events := []storage.Event{
		buildEventWith(map[string]any{
			"StartsAt": date.AddDate(0, 0, 1),
			"EndsAt":   date.AddDate(0, 0, 2),
		}),
		buildEventWith(map[string]any{
			"StartsAt": date.AddDate(0, 0, 15),
			"EndsAt":   date.AddDate(0, 0, 50),
		}),
		buildEventWith(map[string]any{
			"StartsAt": date.AddDate(0, 0, 65),
			"EndsAt":   date.AddDate(0, 0, 85),
		}),
	}

	memStorage := New()
	for i := range events {
		err := memStorage.Create(events[i])
		require.Equal(t, nil, err)
	}

	weekEvents := memStorage.ListEventForMonth(date)
	require.ElementsMatch(t, weekEvents, []storage.Event{events[0], events[1]})
}

func buildEvent() storage.Event {
	return storage.Event{
		StartsAt:     date(2033, 6, 1),
		EndsAt:       date(2033, 6, 2),
		Title:        gofakeit.FarmAnimal(),
		Description:  gofakeit.Adjective(),
		ID:           storage.EventID(gofakeit.UUID()),
		UserID:       gofakeit.Email(),
		NotifyBefore: time.Second * 2,
	}
}

func buildEventWith(attrs map[string]any) storage.Event {
	event := buildEvent()

	for key := range attrs {
		switch key {
		case "ID":
			event.ID = attrs[key].(storage.EventID)
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
