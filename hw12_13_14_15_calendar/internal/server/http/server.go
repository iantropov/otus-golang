package internalhttp

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
)

type Server struct { // TODO
	host, port string
	logger     Logger
}

type Logger interface {
	Info(string)
	Infof(string, ...any)
	Error(string)
	Errorf(string, ...any)
}

type Application interface { // TODO
}

func NewServer(host, port string, logger Logger, app Application) *Server {
	return &Server{host, port, logger}
}

type loggerKey string

const logger loggerKey = "logger"

func (s *Server) Start(startCtx context.Context) error {

	go func() {
		mux := http.NewServeMux()

		mux.HandleFunc("/hello", getHello)

		server := &http.Server{
			Addr:    net.JoinHostPort(s.host, s.port),
			Handler: mux,
			BaseContext: func(l net.Listener) context.Context {
				requestCtx := context.WithValue(startCtx, logger, s.logger)
				return requestCtx
			},
		}

		s.logger.Infof("listening http at %s:%s\n", s.host, s.port)

		server.ListenAndServe()
	}()

	<-startCtx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}
