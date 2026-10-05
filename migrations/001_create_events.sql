CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    player_id BIGINT NOT NULL,
    type TEXT NOT NULL,
    payload TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_events_player_id
ON events(player_id);

CREATE INDEX IF NOT EXISTS idx_events_type
ON events(type);

CREATE INDEX IF NOT EXISTS idx_events_created_at
ON events(created_at);