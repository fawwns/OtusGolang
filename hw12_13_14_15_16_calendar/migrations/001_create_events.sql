CREATE TABLE IF NOT EXISTS events (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    duration BIGINT NOT NULL,
    description TEXT,
    user_id TEXT NOT NULL,
    notify_before BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_events_start_time ON events (start_time);