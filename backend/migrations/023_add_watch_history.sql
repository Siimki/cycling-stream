-- Create view for user watch history (races watched > 10 minutes)
-- This extends the existing watch_sessions table with a convenient view
CREATE OR REPLACE VIEW user_watch_history AS
SELECT 
    ws.user_id,
    ws.race_id,
    r.name as race_name,
    r.category as race_category,
    r.start_date as race_start_date,
    COUNT(*) as session_count,
    SUM(ws.duration_seconds) as total_seconds,
    SUM(ws.duration_seconds) / 60.0 as total_minutes,
    MIN(ws.started_at) as first_watched,
    MAX(COALESCE(ws.ended_at, ws.started_at)) as last_watched,
    -- Calculate completion percentage (if race has estimated duration)
    -- For now, we'll use a simple heuristic: if watched > 60 minutes, consider it "watched"
    CASE 
        WHEN SUM(ws.duration_seconds) > 3600 THEN true
        ELSE false
    END as likely_completed
FROM watch_sessions ws
INNER JOIN races r ON ws.race_id = r.id
WHERE ws.duration_seconds IS NOT NULL
    AND ws.duration_seconds >= 600  -- Only include sessions >= 10 minutes
GROUP BY ws.user_id, ws.race_id, r.name, r.category, r.start_date;

-- Create index to support watch history queries
CREATE INDEX IF NOT EXISTS idx_watch_sessions_duration ON watch_sessions(duration_seconds) WHERE duration_seconds IS NOT NULL;

