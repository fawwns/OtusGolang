package internalhttp_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/logger"
	httpserver "github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/server/http"
	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/storage"
)

type mockApp struct {
	events []storage.Event
	err    error
}

func (m *mockApp) CreateEvent(ctx context.Context, e storage.Event) error {
	if m.err != nil {
		return m.err
	}
	m.events = append(m.events, e)
	return nil
}

func (m *mockApp) UpdateEvent(ctx context.Context, e storage.Event) error {
	if m.err != nil {
		return m.err
	}
	for i := range m.events {
		if m.events[i].ID == e.ID {
			m.events[i] = e
			return nil
		}
	}
	return errors.New("not found")
}

func (m *mockApp) DeleteEvent(ctx context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	for i := range m.events {
		if m.events[i].ID == id {
			m.events = append(m.events[:i], m.events[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func (m *mockApp) ListDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.events, nil
}

func (m *mockApp) ListWeekEvents(ctx context.Context, start time.Time) ([]storage.Event, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.events, nil
}

func (m *mockApp) ListMonthEvents(ctx context.Context, start time.Time) ([]storage.Event, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.events, nil
}

func newTestServer(app *mockApp) *httpserver.Server {
	logg := logger.New("debug")
	return httpserver.NewServer(logg, app, "127.0.0.1", 8080)
}

func TestCreateEvent(t *testing.T) {
	mock := &mockApp{}
	srv := newTestServer(mock)

	// Подготовка данных для события
	body := map[string]any{
		"id":            "1",
		"title":         "Test event",
		"start_time":    time.Now().Format(time.RFC3339),
		"duration":      "1h", // Используем строку без экранирования
		"description":   "desc",
		"user_id":       "u1",
		"notify_before": "10m", // Аналогично
	}
	b, _ := json.Marshal(body)

	// Создание запроса
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(b))
	w := httptest.NewRecorder()

	// Выполнение запроса
	srv.Router().ServeHTTP(w, req)

	// Проверка ответа
	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d", w.Code)
	}

	// Проверка, что событие было создано
	if len(mock.events) != 1 {
		t.Fatalf("event was not created")
	}
}

func TestDeleteEvent(t *testing.T) {
	mock := &mockApp{
		events: []storage.Event{
			{ID: "123", Title: "test"},
		},
	}
	srv := newTestServer(mock)

	req := httptest.NewRequest(http.MethodDelete, "/events/123", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}

	if len(mock.events) != 0 {
		t.Fatalf("event was not deleted")
	}
}

func TestListDay(t *testing.T) {
	mock := &mockApp{
		events: []storage.Event{
			{ID: "1", Title: "event"},
		},
	}
	srv := newTestServer(mock)

	req := httptest.NewRequest(http.MethodGet, "/events/day?date=2025-01-01", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}

	var resp []storage.Event
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	if len(resp) != 1 {
		t.Fatalf("expected 1 event, got %d", len(resp))
	}
}
