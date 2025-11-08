package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	httpServer *http.Server
	logger     app.Logger
	app        app.Application
}

func NewServer(logger app.Logger, application app.Application, host string, port int) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		fmt.Fprintf(w, "Hello, Calendar!\n")
		duration := time.Since(start)

		clientIP := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		proto := r.Proto
		userAgent := r.UserAgent()

		logger.Info(fmt.Sprintf("%s %s %s %s %s %d %v %s",
			clientIP,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			method, path, proto, 200, duration, userAgent))
	})

	addr := fmt.Sprintf("%s:%d", host, port)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &Server{
		httpServer: srv,
		logger:     logger,
		app:        application,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Starting HTTP server on " + s.httpServer.Addr)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("HTTP Server ListenAndServe: " + err.Error())
		}
	}()

	<-ctx.Done()
	return s.Stop(context.Background())
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server...")
	return s.httpServer.Shutdown(ctx)
}
