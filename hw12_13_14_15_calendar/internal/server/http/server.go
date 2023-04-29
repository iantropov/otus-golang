package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	host, port string
	logger     Logger
	server     http.Server
}

type Logger interface {
	Info(string)
	Infof(string, ...any)
	Error(string)
	Errorf(string, ...any)
}

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, id storage.EventID, event storage.Event) error
	DeleteEvent(ctx context.Context, id storage.EventID) error
	GetEvent(ctx context.Context, id storage.EventID) (storage.Event, error)
	ListEventForDay(ctx context.Context, day time.Time) []storage.Event
	ListEventForMonth(ctx context.Context, monthStart time.Time) []storage.Event
	ListEventForWeek(ctx context.Context, weekStart time.Time) []storage.Event
}

type serverContext string

const statusCodeKey = serverContext("statusCode")

func NewServer(host, port string, logger Logger, app Application) *Server {
	return &Server{
		host:   host,
		port:   port,
		logger: logger,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", wrapHandler(s.getHello))

	s.server = http.Server{
		Addr:              net.JoinHostPort(s.host, s.port),
		Handler:           loggingMiddleware(s.logger, mux),
		ReadHeaderTimeout: time.Minute,
	}

	s.logger.Infof("listening http at %s:%s\n", s.host, s.port)

	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Infof("stopping http at %s:%s\n", s.host, s.port)
	return s.server.Shutdown(ctx)
}

func wrapHandler[Req any, Res any](handler func(ctx context.Context, req Req) (res Res, err error)) func(w http.ResponseWriter, r *http.Request) {
	return func(resWriter http.ResponseWriter, httpReq *http.Request) {
		ctx := httpReq.Context()

		limitedReader := io.LimitReader(httpReq.Body, 1_000_000)

		body, err := io.ReadAll(limitedReader)
		if err != nil {
			respondWithError(resWriter, httpReq, http.StatusBadRequest, "reading body", err)
			return
		}

		var request Req
		err = json.Unmarshal(body, &request)
		if err != nil {
			respondWithError(resWriter, httpReq, http.StatusBadRequest, "decoding JSON", err)
			return
		}

		response, err := handler(ctx, request)
		if err != nil {
			respondWithError(resWriter, httpReq, http.StatusInternalServerError, "running handler", err)
			return
		}

		rawJSON, err := json.Marshal(response)
		if err != nil {
			respondWithError(resWriter, httpReq, http.StatusInternalServerError, "encoding JSON", err)
			return
		}

		respondWithSuccess(resWriter, httpReq, rawJSON)
	}
}

func respondWithError(resWriter http.ResponseWriter, r *http.Request, status int, text string, err error) {
	writeErrorJSON(resWriter, text, err)
	resWriter.Header().Add("Content-Type", "application/json")
	setResponseCode(resWriter, r, 400)

	resWriter.WriteHeader(status)
}

func respondWithSuccess(resWriter http.ResponseWriter, r *http.Request, rawJSON []byte) {
	resWriter.Header().Add("Content-Type", "application/json")
	setResponseCode(resWriter, r, 200)
	_, _ = resWriter.Write(rawJSON)
}

func writeErrorJSON(w http.ResponseWriter, text string, err error) {
	buf := bytes.NewBufferString("{\"message\":\"")
	buf.WriteString(text)
	buf.WriteString("\",\"error\":\"")
	buf.WriteString(err.Error())
	buf.WriteString("\"}\n")

	w.Write(buf.Bytes())
}

func setResponseCode(w http.ResponseWriter, r *http.Request, statusCode int) {
	ctx := context.WithValue(r.Context(), statusCodeKey, statusCode)
	*r = *(r.WithContext(ctx))

	w.WriteHeader(statusCode)
}
