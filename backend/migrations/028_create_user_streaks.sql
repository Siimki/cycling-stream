-- Create user_streaks table to track current streak state
CREATE TABLE IF NOT EXISTS user_streaks (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current_streak_weeks INTEGER NOT NULL DEFAULT 0,
    last_completed_week_number VARCHAR(10), -- ISO week format: YYYY-WW
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for user_streaks
CREATE INDEX idx_user_streaks_user_id ON user_streaks(user_id);

-- Add comment
COMMENT ON TABLE user_streaks IS 'Tracks current weekly streak state (one row per user, purely cosmetic)';


