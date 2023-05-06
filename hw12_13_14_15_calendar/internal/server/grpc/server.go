package grpc

import (
	"context"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/server"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/pkg/event_service_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	event_service_v1.UnimplementedEventServiceV1Server

	host, port string
	logger     server.Logger
	app        server.Application
	server     *grpc.Server
}

func NewServer(host, port string, logger server.Logger, app server.Application) *Server {
	return &Server{
		host:   host,
		port:   port,
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start() error {
	grpcListener, err := net.Listen("tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		return err
	}

	s.server = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				LoggingInterceptor(s.logger),
			),
		),
	)
	reflection.Register(s.server)
	event_service_v1.RegisterEventServiceV1Server(s.server, s)
	s.logger.Infof("listening grpc at %s:%s\n", s.host, s.port)

	return s.server.Serve(grpcListener)
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Infof("stopping grpc at %s:%s\n", s.host, s.port)
	if s.server == nil {
		return nil
	}

	go func() {
		<-ctx.Done()
		s.server.Stop()
	}()
	s.server.GracefulStop()

	return nil
}
