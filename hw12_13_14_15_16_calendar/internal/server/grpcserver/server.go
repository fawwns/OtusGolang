package grpcserver

import (
	pb "github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/api/proto"
	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	pb.UnimplementedCalendarServiceServer
	app *app.App
}

func New(app *app.App) *Server {
	return &Server{app: app}
}
