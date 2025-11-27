-- Create streams table
CREATE TABLE IF NOT EXISTS streams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    race_id UUID NOT NULL REFERENCES races(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'planned', -- planned, live, ended
    origin_url VARCHAR(500), -- Direct HLS URL from origin server
    cdn_url VARCHAR(500), -- CDN HLS URL
    stream_key VARCHAR(255), -- RTMP stream key
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_streams_race_id ON streams(race_id);
CREATE INDEX idx_streams_status ON streams(status);

