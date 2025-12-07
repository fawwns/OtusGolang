package memorystorage

import (
	"sync"
	"time"

	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex //nolint:unused
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) Create(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range s.events {
		if eventsOverlap(e, event) {
			return storage.ErrDateBusy
		}
	}

	s.events[event.ID] = event
	return nil
}

func (s *Storage) Update(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.events[event.ID]
	if !exists {
		return storage.ErrNotFound
	}

	for id, e := range s.events {
		if id == event.ID {
			continue
		}
		if eventsOverlap(e, event) {
			return storage.ErrDateBusy
		}
	}

	s.events[event.ID] = event

	return nil
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return storage.ErrNotFound
	}

	delete(s.events, id)
	return nil
}

func (s *Storage) ListDay(date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := []storage.Event{}

	for _, event := range s.events {
		if sameDay(event.StartTime, date) {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) ListWeek(start time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := []storage.Event{}
	weekEnd := start.AddDate(0, 0, 7)

	for _, e := range s.events {
		if e.StartTime.After(start) && e.StartTime.Before(weekEnd) {
			result = append(result, e)
		}
	}
	return result, nil
}

func (s *Storage) ListMonth(start time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := []storage.Event{}
	monthEnd := start.AddDate(0, 1, 0)

	for _, e := range s.events {
		if e.StartTime.After(start) && e.StartTime.Before(monthEnd) {
			result = append(result, e)
		}
	}
	return result, nil
}

func eventsOverlap(a, b storage.Event) bool {
	endA := a.StartTime.Add(a.Duration)
	endB := b.StartTime.Add(b.Duration)
	return a.StartTime.Before(endB) && b.StartTime.Before(endA)
}

func sameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
