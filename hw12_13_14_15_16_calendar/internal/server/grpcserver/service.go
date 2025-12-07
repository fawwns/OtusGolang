package grpcserver

import (
	"context"

	pb "github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/api/proto"
	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {

	e := storage.Event{
		ID:           "",
		Title:        req.Event.Title,
		Description:  req.Event.Description,
		UserID:       req.Event.UserId,
		StartTime:    req.Event.StartTime.AsTime(),
		Duration:     req.Event.Duration.AsDuration(),
		NotifyBefore: req.Event.NotifyBefore.AsDuration(),
	}

	err := s.app.CreateEvent(ctx, e)
	if err != nil {
		return nil, err
	}

	return &pb.CreateEventResponse{
		Event: convertToProtoEvent(e),
	}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	err := s.app.DeleteEvent(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteEventResponse{}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {

	e := storage.Event{
		ID:           req.Event.Id,
		Title:        req.Event.Title,
		Description:  req.Event.Description,
		UserID:       req.Event.UserId,
		StartTime:    req.Event.StartTime.AsTime(),
		Duration:     req.Event.Duration.AsDuration(),
		NotifyBefore: req.Event.NotifyBefore.AsDuration(),
	}

	err := s.app.UpdateEvent(ctx, e)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateEventResponse{
		Event: convertToProtoEvent(e),
	}, nil
}
func (s *Server) ListDay(ctx context.Context, req *pb.ListDayRequest) (*pb.ListEventsResponse, error) {
	date := req.Date.AsTime()
	events, err := s.app.ListDayEvents(ctx, date)
	if err != nil {
		return nil, err
	}

	return convertToPbList(events), nil
}

func (s *Server) ListWeek(ctx context.Context, req *pb.ListWeekRequest) (*pb.ListEventsResponse, error) {
	start := req.Start.AsTime()
	events, err := s.app.ListWeekEvents(ctx, start)
	if err != nil {
		return nil, err
	}

	return convertToPbList(events), nil
}

func (s *Server) ListMonth(ctx context.Context, req *pb.ListMonthRequest) (*pb.ListEventsResponse, error) {
	start := req.Start.AsTime()
	events, err := s.app.ListMonthEvents(ctx, start)
	if err != nil {
		return nil, err
	}

	return convertToPbList(events), nil
}

func convertToProtoEvent(e storage.Event) *pb.Event {
	return &pb.Event{
		Id:           e.ID,
		Title:        e.Title,
		Description:  e.Description,
		UserId:       e.UserID,
		StartTime:    timestamppb.New(e.StartTime),
		Duration:     durationpb.New(e.Duration),
		NotifyBefore: durationpb.New(e.NotifyBefore),
	}
}

func convertToPbList(events []storage.Event) *pb.ListEventsResponse {
	res := &pb.ListEventsResponse{}

	for _, e := range events {
		res.Events = append(res.Events, &pb.Event{
			Id:           e.ID,
			Title:        e.Title,
			Description:  e.Description,
			UserId:       e.UserID,
			StartTime:    timestamppb.New(e.StartTime),
			Duration:     durationpb.New(e.Duration),
			NotifyBefore: durationpb.New(e.NotifyBefore),
		})
	}

	return res
}
