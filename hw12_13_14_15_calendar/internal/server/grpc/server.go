package grpc

import (
	"context"
	"net"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/pkg/event_service_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Logger interface {
	Info(string)
	Infof(string, ...any)
	Error(string)
	Errorf(string, ...any)
}

type Application interface {
	Create(ctx context.Context, event storage.Event) error
	Update(ctx context.Context, id storage.EventID, event storage.Event) error
	Delete(ctx context.Context, id storage.EventID) error
	Get(ctx context.Context, id storage.EventID) (storage.Event, error)
	ListEventForDay(ctx context.Context, day time.Time) []storage.Event
	ListEventForMonth(ctx context.Context, monthStart time.Time) []storage.Event
	ListEventForWeek(ctx context.Context, weekStart time.Time) []storage.Event
}

type Server struct {
	event_service_v1.UnimplementedEventServiceV1Server

	host, port string
	logger     Logger
	app        Application
	server     *grpc.Server
}

func NewServer(host, port string, logger Logger, app Application) *Server {
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
	// grpc.UnaryInterceptor(
	// grpcMiddleware.ChainUnaryServer(
	// 	interceptors.LoggingInterceptor,
	// ),
	// ),
	)
	reflection.Register(s.server)
	event_service_v1.RegisterEventServiceV1Server(s.server, s)
	s.logger.Infof("listening grpc at %s:%s\n", s.host, s.port)

	return s.server.Serve(grpcListener)
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Infof("stopping grpc at %s:%s\n", s.host, s.port)
	go func() {
		<-ctx.Done()
		s.server.Stop()
	}()
	s.server.GracefulStop()

	return nil
}

func (s *Server) Create(context.Context, *event_service_v1.CreateRequest) (*emptypb.Empty, error) {
	s.logger.Info("HEELOO from CREATE")
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (s *Server) Update(context.Context, *event_service_v1.UpdateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (s *Server) Delete(context.Context, *event_service_v1.IDRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (s *Server) Get(context.Context, *event_service_v1.IDRequest) (*event_service_v1.GetResponse, error) {
	s.logger.Info("HEELOO from GET")
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (s *Server) ListEventForDay(context.Context, *event_service_v1.TimeRequest) (*event_service_v1.ListEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventForDay not implemented")
}
func (s *Server) ListEventForWeek(context.Context, *event_service_v1.TimeRequest) (*event_service_v1.ListEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventForWeek not implemented")
}
func (s *Server) ListEventForMonth(context.Context, *event_service_v1.TimeRequest) (*event_service_v1.ListEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventForMonth not implemented")
}
