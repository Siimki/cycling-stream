-- Create viewer_sessions table to track concurrent and unique viewers per race
CREATE TABLE IF NOT EXISTS viewer_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE, -- NULL for anonymous viewers
    race_id UUID NOT NULL REFERENCES races(id) ON DELETE CASCADE,
    session_token VARCHAR(255) NOT NULL, -- Unique token for anonymous viewers or user_id for authenticated
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Updated on heartbeat
    ended_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_viewer_sessions_race_id ON viewer_sessions(race_id);
CREATE INDEX idx_viewer_sessions_user_id ON viewer_sessions(user_id);
CREATE INDEX idx_viewer_sessions_session_token ON viewer_sessions(session_token);
CREATE INDEX idx_viewer_sessions_active ON viewer_sessions(race_id, is_active) WHERE is_active = TRUE;
CREATE INDEX idx_viewer_sessions_started_at ON viewer_sessions(started_at);

-- Create view for concurrent viewers per race
CREATE OR REPLACE VIEW concurrent_viewers AS
SELECT 
    race_id,
    COUNT(*) as concurrent_count,
    COUNT(DISTINCT user_id) FILTER (WHERE user_id IS NOT NULL) as authenticated_count,
    COUNT(*) FILTER (WHERE user_id IS NULL) as anonymous_count
FROM viewer_sessions
WHERE is_active = TRUE
GROUP BY race_id;

-- Create view for unique viewers per race (all time)
CREATE OR REPLACE VIEW unique_viewers AS
SELECT 
    race_id,
    COUNT(DISTINCT COALESCE(user_id::text, session_token)) as unique_viewer_count,
    COUNT(DISTINCT user_id) FILTER (WHERE user_id IS NOT NULL) as unique_authenticated_count,
    COUNT(DISTINCT session_token) FILTER (WHERE user_id IS NULL) as unique_anonymous_count
FROM viewer_sessions
GROUP BY race_id;

