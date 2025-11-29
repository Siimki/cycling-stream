-- Stream providers capture the concrete hosting backend (Bunny, YouTube, etc.)
CREATE TABLE IF NOT EXISTS stream_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,
    provider TEXT NOT NULL, -- e.g. bunny_stream, youtube_embed
    provider_video_id TEXT NOT NULL, -- provider-specific identifier
    provider_url TEXT, -- optional convenience URL
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (stream_id, provider)
);

CREATE INDEX idx_stream_providers_provider ON stream_providers(provider);
CREATE INDEX idx_stream_providers_provider_video_id ON stream_providers(provider_video_id);

-- Raw playback event log (append-only)
CREATE TABLE IF NOT EXISTS playback_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,
    viewer_session_id UUID REFERENCES viewer_sessions(id) ON DELETE SET NULL,
    client_id TEXT NOT NULL,
    event_type TEXT NOT NULL, -- play, pause, heartbeat, ended, error, buffer_start, buffer_end
    video_time_seconds INTEGER,
    country TEXT DEFAULT 'unknown',
    device_type TEXT DEFAULT 'unknown',
    extra JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_playback_events_stream_created_at ON playback_events(stream_id, created_at);
CREATE INDEX idx_playback_events_stream_client ON playback_events(stream_id, client_id);
CREATE INDEX idx_playback_events_event_type ON playback_events(event_type);
