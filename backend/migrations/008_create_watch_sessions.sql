-- Create watch_sessions table to track user watch time
CREATE TABLE IF NOT EXISTS watch_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    race_id UUID NOT NULL REFERENCES races(id) ON DELETE CASCADE,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP WITH TIME ZONE,
    duration_seconds INTEGER, -- Calculated duration in seconds
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_watch_sessions_user_id ON watch_sessions(user_id);
CREATE INDEX idx_watch_sessions_race_id ON watch_sessions(race_id);
CREATE INDEX idx_watch_sessions_started_at ON watch_sessions(started_at);
CREATE INDEX idx_watch_sessions_user_race ON watch_sessions(user_id, race_id);

-- Create view for aggregated watch time per user per race
CREATE OR REPLACE VIEW watch_time_aggregated AS
SELECT 
    user_id,
    race_id,
    COUNT(*) as session_count,
    SUM(duration_seconds) as total_seconds,
    SUM(duration_seconds) / 60.0 as total_minutes,
    MIN(started_at) as first_watched,
    MAX(ended_at) as last_watched
FROM watch_sessions
WHERE duration_seconds IS NOT NULL
GROUP BY user_id, race_id;

