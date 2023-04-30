package grpc

import (
	"context"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/pkg/event_service_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Create(ctx context.Context, req *event_service_v1.CreateRequest) (*emptypb.Empty, error) {
	event := inEvent(req.GetEvent())
	err := s.app.CreateEvent(ctx, event)
	return &emptypb.Empty{}, err
}

func (s *Server) Update(ctx context.Context, req *event_service_v1.UpdateRequest) (*emptypb.Empty, error) {
	ID := storage.EventID(req.GetId())
	event := inEvent(req.GetEvent())
	err := s.app.UpdateEvent(ctx, ID, event)
	return &emptypb.Empty{}, err
}

func (s *Server) Delete(ctx context.Context, req *event_service_v1.IDRequest) (*emptypb.Empty, error) {
	ID := storage.EventID(req.GetId())
	err := s.app.DeleteEvent(ctx, ID)
	return &emptypb.Empty{}, err
}

func (s *Server) Get(ctx context.Context, req *event_service_v1.IDRequest) (*event_service_v1.GetResponse, error) {
	ID := storage.EventID(req.GetId())
	event, err := s.app.GetEvent(ctx, ID)
	if err != nil {
		return &event_service_v1.GetResponse{}, err
	}
	return &event_service_v1.GetResponse{
		Event: outEvent(event),
	}, nil
}

func (s *Server) ListEventForDay(
	ctx context.Context,
	req *event_service_v1.TimeRequest,
) (*event_service_v1.ListEventResponse, error) {
	at := inTime(req.GetAt())
	events := s.app.ListEventForDay(ctx, at)
	return &event_service_v1.ListEventResponse{
		Events: outEvents(events),
	}, nil
}

func (s *Server) ListEventForWeek(
	ctx context.Context,
	req *event_service_v1.TimeRequest,
) (*event_service_v1.ListEventResponse, error) {
	at := inTime(req.GetAt())
	events := s.app.ListEventForWeek(ctx, at)
	return &event_service_v1.ListEventResponse{
		Events: outEvents(events),
	}, nil
}

func (s *Server) ListEventForMonth(
	ctx context.Context,
	req *event_service_v1.TimeRequest,
) (*event_service_v1.ListEventResponse, error) {
	at := inTime(req.GetAt())
	events := s.app.ListEventForMonth(ctx, at)
	return &event_service_v1.ListEventResponse{
		Events: outEvents(events),
	}, nil
}

func inEvent(reqEvent *event_service_v1.Event) storage.Event {
	if reqEvent == nil {
		return storage.Event{}
	}
	return storage.Event{
		ID:           storage.EventID(reqEvent.GetId()),
		Title:        reqEvent.GetTitle(),
		StartsAt:     inTime(reqEvent.GetStartsAt()),
		EndsAt:       inTime(reqEvent.GetEndsAt()),
		Description:  reqEvent.GetDescription(),
		UserID:       reqEvent.GetUserId(),
		NotifyBefore: time.Second * time.Duration(reqEvent.GetNotifyBeforeSeconds()),
	}
}

func inTime(val string) time.Time {
	timeVal, err := time.Parse("2006-01-02T15:04:05", val)
	if err != nil {
		return time.Time{}
	}
	return timeVal
}

func outEvent(storageEvent storage.Event) *event_service_v1.Event {
	return &event_service_v1.Event{
		Id:                  string(storageEvent.ID),
		Title:               storageEvent.Title,
		StartsAt:            storageEvent.StartsAt.String(),
		EndsAt:              storageEvent.EndsAt.String(),
		Description:         storageEvent.Description,
		UserId:              storageEvent.UserID,
		NotifyBeforeSeconds: int32(storageEvent.NotifyBefore.Seconds()),
	}
}

func outEvents(events []storage.Event) []*event_service_v1.Event {
	res := make([]*event_service_v1.Event, len(events))
	for i := range events {
		res[i] = outEvent(events[i])
	}
	return res
}
