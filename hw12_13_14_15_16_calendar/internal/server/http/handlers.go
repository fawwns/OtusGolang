package internalhttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/storage"
	"github.com/gorilla/mux"
)

func (s *Server) registerHandlers() {
	s.router.HandleFunc("/events", s.createEvent).Methods("POST")
	s.router.HandleFunc("/events/{id}", s.updateEvent).Methods("PUT")
	s.router.HandleFunc("/events/{id}", s.deleteEvent).Methods("DELETE")
	s.router.HandleFunc("/events/day", s.listDay).Methods("GET")
	s.router.HandleFunc("/events/week", s.listWeek).Methods("GET")
	s.router.HandleFunc("/events/month", s.listMonth).Methods("GET")
}

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	e := storage.Event{
		Title:        req.Title,
		Description:  req.Description,
		StartTime:    req.StartTime,
		Duration:     req.Duration.Duration,
		UserID:       req.UserID,
		NotifyBefore: req.NotifyBefore.Duration,
	}

	if err := s.app.CreateEvent(r.Context(), e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": e.ID})
}

func (s *Server) updateEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "missing event ID", http.StatusBadRequest)
		return
	}

	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse("2006-01-02T15:04:05", req.StartTime)
	if err != nil {
		http.Error(w, "invalid start_time format", http.StatusBadRequest)
		return
	}

	e := storage.Event{
		ID:           id,
		Title:        req.Title,
		Description:  req.Description,
		StartTime:    startTime,
		Duration:     time.Duration(req.Duration) * time.Second,
		NotifyBefore: time.Duration(req.NotifyBefore) * time.Second,
	}

	if err := s.app.UpdateEvent(r.Context(), e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "missing event id", http.StatusBadRequest)
		return
	}

	if err := s.app.DeleteEvent(r.Context(), id); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "event not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) listDay(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid date format", http.StatusBadRequest)
		return
	}

	events, err := s.app.ListDayEvents(r.Context(), date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}

func (s *Server) listWeek(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("start")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid start date format", http.StatusBadRequest)
		return
	}

	events, err := s.app.ListWeekEvents(r.Context(), date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}

func (s *Server) listMonth(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("start")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid start date format", http.StatusBadRequest)
		return
	}

	events, err := s.app.ListMonthEvents(r.Context(), date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}
