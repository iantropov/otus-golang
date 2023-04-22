package internalhttp

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
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

type Application interface { // TODO
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

func (s *Server) Start(startCtx context.Context) error {
	go func() {
		mux := http.NewServeMux()

		mux.HandleFunc("/hello", s.getHello)

		s.server = http.Server{
			Addr:              net.JoinHostPort(s.host, s.port),
			Handler:           loggingMiddleware(s.logger, mux),
			ReadHeaderTimeout: time.Minute,
		}

		s.logger.Infof("listening http at %s:%s\n", s.host, s.port)
		s.server.ListenAndServe()
	}()

	<-startCtx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Infof("stopping http at %s:%s\n", s.host, s.port)
	return s.server.Shutdown(ctx)
}

func (s *Server) getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")

	s.setResponseCode(w, r, 200)
	io.WriteString(w, "Hello, HTTP!\n")
}

func (s *Server) setResponseCode(w http.ResponseWriter, r *http.Request, statusCode int) {
	ctx := context.WithValue(r.Context(), statusCodeKey, statusCode)
	*r = *(r.WithContext(ctx))

	w.WriteHeader(statusCode)
}
