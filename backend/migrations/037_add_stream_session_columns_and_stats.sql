-- Extend viewer_sessions to capture stream-level metadata for analytics.
ALTER TABLE viewer_sessions
    ADD COLUMN IF NOT EXISTS stream_id UUID REFERENCES streams(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS country TEXT DEFAULT 'unknown',
    ADD COLUMN IF NOT EXISTS device_type TEXT DEFAULT 'unknown',
    ADD COLUMN IF NOT EXISTS total_watch_seconds INTEGER,
    ADD COLUMN IF NOT EXISTS buffer_seconds INTEGER,
    ADD COLUMN IF NOT EXISTS error_count INTEGER;

CREATE INDEX IF NOT EXISTS idx_viewer_sessions_stream_id ON viewer_sessions(stream_id);
CREATE INDEX IF NOT EXISTS idx_viewer_sessions_country ON viewer_sessions(country);
CREATE INDEX IF NOT EXISTS idx_viewer_sessions_device_type ON viewer_sessions(device_type);

-- Aggregated per-stream stats (overall; can extend with date if needed later).
CREATE TABLE IF NOT EXISTS stream_stats (
    stream_id UUID PRIMARY KEY REFERENCES streams(id) ON DELETE CASCADE,
    unique_viewers INTEGER DEFAULT 0,
    total_watch_seconds BIGINT DEFAULT 0,
    avg_watch_seconds INTEGER DEFAULT 0,
    peak_concurrent_viewers INTEGER DEFAULT 0,
    top_countries JSONB,
    device_breakdown JSONB,
    last_calculated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
