package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSQLStorage_CreateAndList(t *testing.T) {
	ctx := context.Background()
	s := sqlstorage.New("postgres://user:pass@localhost:5432/calendar_test?sslmode=disable")
	require.NoError(t, s.Connect(ctx))

	event := storage.Event{
		ID:           uuid.New().String(),
		Title:        "Test event",
		StartTime:    time.Now(),
		Duration:     time.Hour,
		Description:  "test",
		UserID:       "123",
		NotifyBefore: time.Minute * 10,
	}

	require.NoError(t, s.Create(event))

	events, err := s.ListDay(time.Now())
	require.NoError(t, err)
	require.NotEmpty(t, events)
}
