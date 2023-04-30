package internalhttp

import (
	"context"
	"errors"
	"testing"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/server/mocks"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestGetEvent(t *testing.T) {
	ctx := context.Background()
	existingEvent := storage.Event{}
	errEventNotFound := errors.New("event not found")
	eventID := storage.EventID("123")

	tests := []struct {
		title      string
		event      storage.Event
		err        error
		GetAppMock func(t *testing.T) *mocks.Application
	}{
		{
			title: "positive case",
			event: existingEvent,
			GetAppMock: func(t *testing.T) *mocks.Application {
				t.Helper()

				appMock := mocks.NewApplication(t)
				appMock.On("GetEvent", ctx, eventID).Return(existingEvent, nil)
				return appMock
			},
		},
		{
			title: "event not found",
			err:   errEventNotFound,
			GetAppMock: func(t *testing.T) *mocks.Application {
				t.Helper()

				appMock := mocks.NewApplication(t)
				appMock.On("GetEvent", ctx, eventID).Return(storage.Event{}, errEventNotFound)
				return appMock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			logger := mocks.NewLogger(t)
			appMock := tt.GetAppMock(t)
			server := NewServer("", "", logger, appMock)

			res, err := server.getEvent(ctx, eventIDRequest{string(eventID)})
			if tt.err != nil {
				require.Equal(t, tt.err, err)
			} else {
				require.Equal(t, outEvent(tt.event), res.Event)
				require.Equal(t, nil, err)
			}

			appMock.AssertExpectations(t)
		})
	}
}
