package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) *Storage {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return &Storage{db: db}
}

func (s *Storage) Connect(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) Create(event storage.Event) error {
	_, err := s.db.Exec(
		`INSERT INTO events (id, title, start_time, duration, description, user_id, notify_before)
         VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		event.ID, event.Title, event.StartTime, event.Duration,
		event.Description, event.UserID, event.NotifyBefore,
	)
	return err
}

func (s *Storage) Update(event storage.Event) error {
	res, err := s.db.Exec(`
		UPDATE events
		SET title=$2, start_time=$3, duration=$4, description=$5, user_id=$6, notify_before=$7
		WHERE id=$1
	`, event.ID, event.Title, event.StartTime, int64(event.Duration.Seconds()), event.Description, event.UserID, int64(event.NotifyBefore.Seconds()))
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (s *Storage) Delete(id string) error {
	res, err := s.db.Exec(`DELETE FROM events WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (s *Storage) ListDay(date time.Time) ([]storage.Event, error) {
	rows, err := s.db.Query(`
		SELECT id, title, start_time, duration, description, user_id, notify_before
		FROM events
		WHERE start_time::date = $1
	`, date.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var e storage.Event
		var durationSec, notifySec int64
		if err := rows.Scan(&e.ID, &e.Title, &e.StartTime, &durationSec, &e.Description, &e.UserID, &notifySec); err != nil {
			return nil, err
		}
		e.Duration = time.Duration(durationSec) * time.Second
		e.NotifyBefore = time.Duration(notifySec) * time.Second
		events = append(events, e)
	}
	return events, nil
}

func (s *Storage) ListWeek(start time.Time) ([]storage.Event, error) {
	end := start.AddDate(0, 0, 7)
	rows, err := s.db.Query(`
		SELECT id, title, start_time, duration, description, user_id, notify_before
		FROM events
		WHERE start_time >= $1 AND start_time < $2
	`, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var e storage.Event
		var durationSec, notifySec int64
		if err := rows.Scan(&e.ID, &e.Title, &e.StartTime, &durationSec, &e.Description, &e.UserID, &notifySec); err != nil {
			return nil, err
		}
		e.Duration = time.Duration(durationSec) * time.Second
		e.NotifyBefore = time.Duration(notifySec) * time.Second
		events = append(events, e)
	}
	return events, nil
}

func (s *Storage) ListMonth(start time.Time) ([]storage.Event, error) {
	end := start.AddDate(0, 1, 0)
	rows, err := s.db.Query(`
		SELECT id, title, start_time, duration, description, user_id, notify_before
		FROM events
		WHERE start_time >= $1 AND start_time < $2
	`, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var e storage.Event
		var durationSec, notifySec int64
		if err := rows.Scan(&e.ID, &e.Title, &e.StartTime, &durationSec, &e.Description, &e.UserID, &notifySec); err != nil {
			return nil, err
		}
		e.Duration = time.Duration(durationSec) * time.Second
		e.NotifyBefore = time.Duration(notifySec) * time.Second
		events = append(events, e)
	}
	return events, nil
}
