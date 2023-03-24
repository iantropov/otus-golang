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
		StartsAt:     time.Now().Add(time.Minute),
		EndsAt:       time.Now().Add(time.Hour),
		Title:        gofakeit.FarmAnimal(),
		Description:  gofakeit.Adjective(),
		Id:           storage.EventId(gofakeit.UUID()),
		UserId:       gofakeit.Email(),
		NotifyBefore: time.Second * 2,
	}
	err := memStorage.Create()

	require.Equal(t, nil, err)

}
