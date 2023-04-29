package internalhttp

import (
	"context"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type helloRequest struct {
	ID string `json:"id"`
}

type helloResponse struct {
	Success bool `json:"success"`
}

type event struct {
	ID                  string    `json:"id"`
	Title               string    `json:"title"`
	StartsAt            time.Time `json:"starts_at"`
	EndsAt              time.Time `json:"ends_at"`
	Description         string    `json:"description"`
	UserID              string    `json:"user_id"`
	NotifyBeforeSeconds int       `json:"notify_before_seconds"`
}

type createEventRequest struct {
	Event event `json:"event"`
}

type updateEventRequest struct {
	ID    string `json:"id"`
	Event event  `json:"event"`
}

type eventIDRequest struct {
	ID string `json:"id"`
}

type listEventsRequest struct {
	At time.Time `json:"at"`
}

type emptyResponse struct {
}

type getEventResponse struct {
	Event event `json:"event"`
}

type listEventResponse struct {
	Events []event `json:"events"`
}

func (s *Server) getHello(ctx context.Context, req helloRequest) (helloResponse, error) {
	s.logger.Info("got /hello request: " + req.ID)

	return helloResponse{
		Success: true,
	}, nil
}

func (s *Server) createEvent(ctx context.Context, req createEventRequest) (emptyResponse, error) {
	event := inEvent(req.Event)
	err := s.app.CreateEvent(ctx, event)
	return emptyResponse{}, err
}

func (s *Server) updateEvent(ctx context.Context, req updateEventRequest) (emptyResponse, error) {
	ID := storage.EventID(req.ID)
	event := inEvent(req.Event)
	err := s.app.UpdateEvent(ctx, ID, event)
	return emptyResponse{}, err
}

func (s *Server) deleteEvent(ctx context.Context, req eventIDRequest) (emptyResponse, error) {
	ID := storage.EventID(req.ID)
	err := s.app.DeleteEvent(ctx, ID)
	return emptyResponse{}, err
}

func (s *Server) getEvent(ctx context.Context, req eventIDRequest) (getEventResponse, error) {
	ID := storage.EventID(req.ID)
	event, err := s.app.GetEvent(ctx, ID)
	if err != nil {
		return getEventResponse{}, err
	}
	return getEventResponse{
		Event: outEvent(event),
	}, nil
}

func (s *Server) listEventForDay(ctx context.Context, req listEventsRequest) (listEventResponse, error) {
	events := s.app.ListEventForDay(ctx, req.At)
	return listEventResponse{
		Events: outEvents(events),
	}, nil
}

func (s *Server) listEventForWeek(ctx context.Context, req listEventsRequest) (listEventResponse, error) {
	events := s.app.ListEventForWeek(ctx, req.At)
	return listEventResponse{
		Events: outEvents(events),
	}, nil
}

func (s *Server) listEventForMonth(ctx context.Context, req listEventsRequest) (listEventResponse, error) {
	events := s.app.ListEventForMonth(ctx, req.At)
	return listEventResponse{
		Events: outEvents(events),
	}, nil
}

func inEvent(reqEvent event) storage.Event {
	return storage.Event{
		ID:           storage.EventID(reqEvent.ID),
		Title:        reqEvent.Title,
		StartsAt:     reqEvent.StartsAt,
		EndsAt:       reqEvent.EndsAt,
		Description:  reqEvent.Description,
		UserID:       reqEvent.UserID,
		NotifyBefore: time.Second * time.Duration(reqEvent.NotifyBeforeSeconds),
	}
}

func outEvent(storageEvent storage.Event) event {
	return event{
		ID:                  string(storageEvent.ID),
		Title:               storageEvent.Title,
		StartsAt:            storageEvent.StartsAt,
		EndsAt:              storageEvent.EndsAt,
		Description:         storageEvent.Description,
		UserID:              storageEvent.UserID,
		NotifyBeforeSeconds: int(storageEvent.NotifyBefore.Seconds()),
	}
}

func outEvents(events []storage.Event) []event {
	res := make([]event, len(events))
	for i := range events {
		res[i] = outEvent(events[i])
	}
	return res
}
