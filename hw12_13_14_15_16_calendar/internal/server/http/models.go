package internalhttp

import (
	"fmt"
	"time"
)

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.Duration.String())), nil
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	str := string(b)
	if str[0] == '"' {
		str = str[1 : len(str)-1]
	}
	duration, err := time.ParseDuration(str)
	if err != nil {
		return err
	}
	d.Duration = duration
	return nil
}

// CreateEventRequest модель для запроса на создание события
type CreateEventRequest struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	StartTime    time.Time `json:"start_time"`
	Duration     Duration  `json:"duration"`
	UserID       string    `json:"user_id"`
	NotifyBefore Duration  `json:"notify_before"`
}

type CreateEventResponse struct {
	ID string `json:"id"`
}

type UpdateEventRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	StartTime    string `json:"start_time"`
	Duration     int64  `json:"duration"`
	NotifyBefore int64  `json:"notify_before"`
}

type UpdateEventResponse struct {
	Message string `json:"message"`
}

type DeleteEventRequest struct {
	ID string `json:"id"`
}

type DeleteEventResponse struct {
	Message string `json:"message"`
}

type ListDayEventsRequest struct {
	Date string `json:"date"`
}

type ListDayEventsResponse struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	StartTime    time.Time     `json:"start_time"`
	Duration     time.Duration `json:"duration"`
	NotifyBefore time.Duration `json:"notify_before"`
}

type ListWeekEventsRequest struct {
	StartDate string `json:"start"`
}

type ListWeekEventsResponse struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	StartTime    time.Time     `json:"start_time"`
	Duration     time.Duration `json:"duration"`
	NotifyBefore time.Duration `json:"notify_before"`
}

type ListMonthEventsRequest struct {
	StartDate string `json:"start"`
}

type ListMonthEventsResponse struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	StartTime    time.Time     `json:"start_time"`
	Duration     time.Duration `json:"duration"`
	NotifyBefore time.Duration `json:"notify_before"`
}
