package app

import (
	"context"
	"errors"
	"time"

	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/storage"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Error(msg string)
}

type Storage interface {
	Create(event storage.Event) error
	Update(event storage.Event) error
	Delete(id string) error
	ListDay(date time.Time) ([]storage.Event, error)
	ListWeek(start time.Time) ([]storage.Event, error)
	ListMonth(start time.Time) ([]storage.Event, error)
}

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, e storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	ListDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListWeekEvents(ctx context.Context, start time.Time) ([]storage.Event, error)
	ListMonthEvents(ctx context.Context, start time.Time) ([]storage.Event, error)
}

type App struct {
	logger  Logger
	storage Storage
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	a.logger.Debug("App: CreateEvent called")

	events, _ := a.storage.ListDay(event.StartTime)
	for _, e := range events {
		if eventsOverlap(e, event) {
			a.logger.Error("date busy: " + event.Title)
			return storage.ErrDateBusy
		}
	}

	if err := a.storage.Create(event); err != nil {
		if errors.Is(err, storage.ErrDateBusy) {
			a.logger.Info("Time slot is busy for event: " + event.Title)
			return err
		}
		a.logger.Error("Failed to create event: " + err.Error())
		return err
	}

	a.logger.Info("Event created: " + event.ID)
	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	a.logger.Debug("App: DeleteEvent called")

	if err := a.storage.Delete(id); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			a.logger.Info("Event not found: " + id)
			return err
		}
		a.logger.Error("Failed to delete event: " + err.Error())
		return err
	}

	a.logger.Info("Event deleted: " + id)
	return nil
}

func (a *App) UpdateEvent(ctx context.Context, e storage.Event) error {
	return a.storage.Update(e)
}

func (a *App) ListDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	a.logger.Debug("App: ListDay called")
	return a.storage.ListDay(date)
}

func (a *App) ListWeekEvents(ctx context.Context, start time.Time) ([]storage.Event, error) {
	a.logger.Debug("App: ListWeek called")
	return a.storage.ListWeek(start)
}

func (a *App) ListMonthEvents(ctx context.Context, start time.Time) ([]storage.Event, error) {
	a.logger.Debug("App: ListMonth called")
	return a.storage.ListMonth(start)
}

func eventsOverlap(a, b storage.Event) bool {
	endA := a.StartTime.Add(a.Duration)
	endB := b.StartTime.Add(b.Duration)
	return a.StartTime.Before(endB) && b.StartTime.Before(endA)
}
