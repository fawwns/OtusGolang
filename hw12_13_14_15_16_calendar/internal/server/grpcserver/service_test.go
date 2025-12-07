package grpcserver_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/api/proto"
	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/app"
	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/logger"
	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/server/grpcserver"
	memorystorage "github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/storage/memory"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const bufSize = 1024 * 1024

func startTestServer(t *testing.T) (proto.CalendarServiceClient, func()) {
	lis := bufconn.Listen(bufSize)

	logg := logger.New("debug")
	st := memorystorage.New()
	ap := app.New(logg, st)

	s := grpc.NewServer()
	proto.RegisterCalendarServiceServer(s, grpcserver.New(ap))

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Fatalf("server exited with error: %v", err)
		}
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}

	return proto.NewCalendarServiceClient(conn), func() {
		s.Stop()
		conn.Close()
	}
}

func TestCreateEvent(t *testing.T) {
	client, close := startTestServer(t)
	defer close()

	resp, err := client.CreateEvent(context.Background(), &proto.CreateEventRequest{
		Event: &proto.Event{
			Title:        "Test",
			Description:  "Desc",
			UserId:       "u1",
			StartTime:    timestamppb.Now(),
			Duration:     durationpb.New(time.Hour),
			NotifyBefore: durationpb.New(time.Minute),
		},
	})
	if err != nil {
		t.Fatalf("CreateEvent error: %v", err)
	}

	if resp.Event.Title != "Test" {
		t.Fatalf("wrong title")
	}
}

func TestListDay(t *testing.T) {
	client, close := startTestServer(t)
	defer close()

	now := time.Now()

	_, _ = client.CreateEvent(context.Background(), &proto.CreateEventRequest{
		Event: &proto.Event{
			Title:        "DayEvent",
			UserId:       "u1",
			StartTime:    timestamppb.New(now),
			Duration:     durationpb.New(time.Hour),
			NotifyBefore: durationpb.New(time.Minute),
		},
	})

	resp, err := client.ListDay(context.Background(), &proto.ListDayRequest{
		Date: timestamppb.New(now),
	})
	if err != nil {
		t.Fatalf("ListDay error: %v", err)
	}

	if len(resp.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(resp.Events))
	}
}
