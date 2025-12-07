package storage

import (
	"errors"
	"time"
)

var (
	ErrDateBusy = errors.New("this time slot is already occupied")
	ErrNotFound = errors.New("event not found")
)

type Event struct {
	ID           string
	Title        string
	StartTime    time.Time
	Duration     time.Duration
	Description  string
	UserID       string
	NotifyBefore time.Duration
}

type Storage interface {
	Create(event Event) error
	Update(event Event) error
	Delete(id string) error
	ListDay(date time.Time) ([]Event, error)
	ListWeek(start time.Time) ([]Event, error)
	ListMonth(start time.Time) ([]Event, error)
}
