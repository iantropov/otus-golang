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

	event := storage.Event{
		StartsAt:     time.Date(2033, 6, 1, 0, 0, 0, 0, time.Local),
		EndsAt:       time.Date(2033, 6, 2, 0, 0, 0, 0, time.Local),
		Title:        gofakeit.FarmAnimal(),
		Description:  gofakeit.Adjective(),
		Id:           storage.EventId(gofakeit.UUID()),
		UserId:       gofakeit.Email(),
		NotifyBefore: time.Second * 2,
	}

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

	resEvents := memStorage.ListEventForDay(time.Date(2033, 6, 1, 0, 0, 0, 0, time.Local))
	require.Equal(t, []storage.Event{event}, resEvents)

	err = memStorage.Delete(event.Id)
	require.Equal(t, nil, err)

	_, err = memStorage.Get(event.Id)
	require.Equal(t, ErrEventNotFound, err)
}
