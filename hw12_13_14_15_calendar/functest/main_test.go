package functest

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
)

type PostgresEvent struct {
	ID           string        `db:"id"`
	Title        string        `db:"title"`
	StartsAt     time.Time     `db:"starts_at"`
	EndsAt       time.Time     `db:"ends_at"`
	Description  string        `db:"description"`
	UserID       string        `db:"user_id"`
	NotifyBefore time.Duration `db:"notify_before"`
}

type HttpEvent struct {
	ID                  string    `json:"id"`
	Title               string    `json:"title"`
	StartsAt            time.Time `json:"startsAt"`
	EndsAt              time.Time `json:"endsAt"`
	Description         string    `json:"description"`
	UserID              string    `json:"userId"`
	NotifyBeforeSeconds int       `json:"notifyBeforeSeconds"`
}

type HttpEventResponse struct {
	Event HttpEvent `json:"event"`
}

// func TestMain(m *testing.M) {
// 	// call flag.Parse() here if TestMain uses flags
// 	os.Exit(m.Run())
// }

var amqpDSN = os.Getenv("TESTS_AMQP_DSN")
var postgresDSN = os.Getenv("TESTS_POSTGRES_DSN")
var apiBaseUrl = os.Getenv("API_BASE_URL")

func init() {
	if amqpDSN == "" {
		amqpDSN = "amqp://guest:guest@localhost:5672/"
	}
	if postgresDSN == "" {
		postgresDSN = "host=localhost port=5432 user=calendar password=password dbname=calendar sslmode=disable"
	}
	if apiBaseUrl == "" {
		apiBaseUrl = "http://localhost:8888"
	}
}

func TestTest(t *testing.T) {
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

		fmt.Printf("%s\n", body)
		var httpEventResponse HttpEventResponse
		err := json.Unmarshal(body, &httpEventResponse)
		panicOnErr(err)

		require.Equal(t, convertToHttpEvent(postgresEvent), httpEventResponse.Event)
	})
}

func convertToHttpEvent(event PostgresEvent) HttpEvent {
	return HttpEvent{
		ID:                  event.ID,
		Title:               event.Title,
		StartsAt:            event.StartsAt,
		EndsAt:              event.EndsAt,
		Description:         event.Description,
		UserID:              event.UserID,
		NotifyBeforeSeconds: int(event.NotifyBefore.Seconds()),
	}
}

func makeHttpRequest(method, url, body string) (int, []byte) {
	req, err := http.NewRequest(method, apiBaseUrl+url, strings.NewReader(body))
	panicOnErr(err)

	res, err := http.DefaultClient.Do(req)
	panicOnErr(err)

	resBody, err := ioutil.ReadAll(res.Body)
	panicOnErr(err)

	return res.StatusCode, resBody
}

func testAMQPMessage(t *testing.T, queue, message string) {
	conn, err := amqp.Dial(amqpDSN)
	panicOnErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	panicOnErr(err)
	defer ch.Close()

	// _, err = ch.QueueDeclare(queue, true, true, true, false, nil)
	// panicOnErr(err)

	events, err := ch.Consume(queue, "", true, true, false, false, nil)
	panicOnErr(err)

	event := <-events
	require.Equal(t, message, event.Body)
}

func testPostgresEvent(t *testing.T, expectedEvent PostgresEvent) {
	db, err := sql.Open("postgres", postgresDSN)
	panicOnErr(err)
	defer db.Close()

	var event PostgresEvent
	selectQuery := `SELECT id, title, starts_at, ends_at, description, user_id, notify_before FROM events WHERE id=$1`
	err = db.QueryRow(selectQuery, expectedEvent.ID).Scan(
		&event.ID,
		&event.Title,
		&event.StartsAt,
		&event.EndsAt,
		&event.Description,
		&event.UserID,
		&event.NotifyBefore,
	)
	panicOnErr(err)

	require.Equal(t, expectedEvent, event)
}

func insertPostresEvent(event PostgresEvent) {
	db, err := sql.Open("postgres", postgresDSN)
	panicOnErr(err)
	defer db.Close()

	deleteQuery := `DELETE FROM events WHERE id=$1`
	_, err = db.Exec(deleteQuery, event.ID)
	panicOnErr(err)

	insertQuery :=
		`INSERT INTO events(id, title, starts_at, ends_at, created_at, description, user_id, notify_before) ` +
			`VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = db.Exec(
		insertQuery,
		&event.ID,
		&event.Title,
		&event.StartsAt,
		&event.EndsAt,
		time.Now(),
		&event.Description,
		&event.UserID,
		&event.NotifyBefore,
	)
	panicOnErr(err)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
