package functest

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
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

type HttpGetEventResponse struct {
	Event HttpEvent `json:"event"`
}

type HttpListEventResponse struct {
	Events []HttpEvent `json:"events"`
}

type HttpCreateEventRequest struct {
	Event HttpEvent `json:"event"`
}

type HttpErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

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

func getAMQPMessage(queue string) string {
	conn, err := amqp.Dial(amqpDSN)
	panicOnErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	panicOnErr(err)
	defer ch.Close()

	_, err = ch.QueueDeclare(queue, false, false, false, false, nil)
	panicOnErr(err)

	events, err := ch.Consume(queue, "", true, false, false, false, nil)
	panicOnErr(err)

	var event amqp.Delivery
	select {
	case event = <-events:
	case <-time.After(15 * time.Second):
	}
	return string(event.Body)
}

func getPostgresEvent(eventID string) (event PostgresEvent) {
	db, err := sql.Open("postgres", postgresDSN)
	panicOnErr(err)
	defer db.Close()

	selectQuery := `SELECT id, title, starts_at, ends_at, description, user_id, notify_before FROM events WHERE id=$1`
	err = db.QueryRow(selectQuery, eventID).Scan(
		&event.ID,
		&event.Title,
		&event.StartsAt,
		&event.EndsAt,
		&event.Description,
		&event.UserID,
		&event.NotifyBefore,
	)
	panicOnErr(err)

	event.StartsAt = event.StartsAt.In(time.UTC)
	event.EndsAt = event.EndsAt.In(time.UTC)

	return event
}

func deletePostgresEvent(eventID string) {
	db, err := sql.Open("postgres", postgresDSN)
	panicOnErr(err)
	defer db.Close()

	deleteQuery := `DELETE FROM events WHERE id=$1`
	_, err = db.Exec(deleteQuery, eventID)
	panicOnErr(err)
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
