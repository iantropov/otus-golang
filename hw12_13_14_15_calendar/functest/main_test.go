package functest

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCalendar(t *testing.T) {
	t.Run("GET event", func(t *testing.T) {
		postgresEvent := PostgresEvent{
			StartsAt:     time.Date(2033, 1, 1, 0, 0, 0, 0, time.UTC),
			EndsAt:       time.Date(2033, 1, 2, 0, 0, 0, 0, time.UTC),
			Title:        "test title",
			Description:  "test description",
			ID:           "123",
			UserID:       "asd@asd.asd",
			NotifyBefore: time.Second * 2,
		}
		insertPostresEvent(postgresEvent)

		status, body := makeHttpRequest(http.MethodGet, "/events/get", `{"id": "123"}`)
		require.Equal(t, http.StatusOK, status)

		var httpEventResponse HttpGetEventResponse
		err := json.Unmarshal(body, &httpEventResponse)
		panicOnErr(err)

		require.Equal(t, convertToHttpEvent(postgresEvent), httpEventResponse.Event)
	})

	t.Run("GET events/listForDay (non-empty)", func(t *testing.T) {
		postgresEvent := PostgresEvent{
			StartsAt:     time.Date(2033, 1, 1, 0, 0, 0, 0, time.UTC),
			EndsAt:       time.Date(2033, 1, 2, 0, 0, 0, 0, time.UTC),
			Title:        "test title",
			Description:  "test description",
			ID:           "123",
			UserID:       "asd@asd.asd",
			NotifyBefore: time.Second * 2,
		}
		insertPostresEvent(postgresEvent)

		status, body := makeHttpRequest(http.MethodGet, "/events/listForDay", `{"at": "2033-01-01T00:00:00Z"}`)
		require.Equal(t, http.StatusOK, status)

		var httpEventResponse HttpListEventResponse
		err := json.Unmarshal(body, &httpEventResponse)
		panicOnErr(err)

		require.Equal(t, 1, len(httpEventResponse.Events))
		require.Equal(t, convertToHttpEvent(postgresEvent), httpEventResponse.Events[0])
	})

	t.Run("GET events/listForDay (empty)", func(t *testing.T) {
		postgresEvent := PostgresEvent{
			StartsAt:     time.Date(2033, 1, 1, 0, 0, 0, 0, time.UTC),
			EndsAt:       time.Date(2033, 1, 2, 0, 0, 0, 0, time.UTC),
			Title:        "test title",
			Description:  "test description",
			ID:           "123",
			UserID:       "asd@asd.asd",
			NotifyBefore: time.Second * 2,
		}
		insertPostresEvent(postgresEvent)

		status, body := makeHttpRequest(http.MethodGet, "/events/listForDay", `{"at": "2038-01-01T00:00:00Z"}`)
		require.Equal(t, http.StatusOK, status)

		var httpEventResponse HttpListEventResponse
		err := json.Unmarshal(body, &httpEventResponse)
		panicOnErr(err)

		require.Equal(t, 0, len(httpEventResponse.Events))
	})

	t.Run("POST events/create (success)", func(t *testing.T) {
		postgresEvent := PostgresEvent{
			StartsAt:     time.Date(2033, 1, 1, 0, 0, 0, 0, time.UTC),
			EndsAt:       time.Date(2033, 1, 2, 0, 0, 0, 0, time.UTC),
			Title:        "test title",
			Description:  "test description",
			ID:           "123",
			UserID:       "asd@asd.asd",
			NotifyBefore: time.Second * 2,
		}
		deletePostgresEvent(postgresEvent.ID)
		httpEvent := convertToHttpEvent(postgresEvent)
		body, err := json.Marshal(HttpCreateEventRequest{
			Event: httpEvent,
		})
		panicOnErr(err)

		status, body := makeHttpRequest(http.MethodPost, "/events/create", string(body))
		require.Equal(t, http.StatusOK, status)

		require.Equal(t, "{}", string(body))

		createdPostgresEvent := getPostgresEvent(postgresEvent.ID)
		require.Equal(t, postgresEvent, createdPostgresEvent)

		notificationMessage := getAMQPMessage("notifications")
		require.Equal(t, `{"id":"123","title":"test title","startsAt":"2033-01-01T00:00:00Z","userId":"asd@asd.asd"}`, notificationMessage)
	})

	t.Run("POST events/create (failure)", func(t *testing.T) {
		postgresEvent := PostgresEvent{
			StartsAt:     time.Date(2033, 1, 1, 0, 0, 0, 0, time.UTC),
			EndsAt:       time.Date(2033, 1, 2, 0, 0, 0, 0, time.UTC),
			Title:        "test title",
			Description:  "test description",
			ID:           "123",
			UserID:       "asd@asd.asd",
			NotifyBefore: time.Second * 2,
		}
		insertPostresEvent(postgresEvent)
		httpEvent := convertToHttpEvent(postgresEvent)
		body, err := json.Marshal(HttpCreateEventRequest{
			Event: httpEvent,
		})
		panicOnErr(err)

		status, body := makeHttpRequest(http.MethodPost, "/events/create", string(body))
		require.Equal(t, http.StatusUnprocessableEntity, status)

		var httpResponse HttpErrorResponse
		err = json.Unmarshal(body, &httpResponse)
		panicOnErr(err)

		require.Equal(t, `running handler`, httpResponse.Message)
		require.Equal(t, `pq: duplicate key value violates unique constraint "events_pkey"`, httpResponse.Error)
	})
}
