package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/app"
	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/logger"
	"github.com/gorilla/mux"
)

type Server struct {
	logg   *logger.Logger
	srv    *http.Server
	app    app.Application
	router *mux.Router
}

func NewServer(logger *logger.Logger, application app.Application, host string, port int) *Server {
	r := mux.NewRouter()
	s := &Server{
		logg:   logger,
		app:    application,
		router: r,
	}

	s.registerHandlers()
	s.router.Use(LoggingMiddleware)

	s.srv = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Info("Starting HTTP server on " + s.srv.Addr)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logg.Error("HTTP Server ListenAndServe: " + err.Error())
		}
	}()

	<-ctx.Done()
	return s.Stop(context.Background())
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Info("Stopping HTTP server...")
	return s.srv.Shutdown(ctx)
}

func (s *Server) Router() http.Handler {
	return s.router
}
