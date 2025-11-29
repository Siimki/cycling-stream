-- Daily Bunny analytics snapshots per video.
CREATE TABLE IF NOT EXISTS bunny_video_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bunny_video_id TEXT NOT NULL,
    stream_id UUID REFERENCES streams(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    views INTEGER DEFAULT 0,
    watch_time_seconds BIGINT DEFAULT 0,
    geo_breakdown JSONB,
    raw_payload JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (bunny_video_id, date)
);

CREATE INDEX idx_bunny_stats_stream_id ON bunny_video_stats(stream_id);
CREATE INDEX idx_bunny_stats_date ON bunny_video_stats(date);
