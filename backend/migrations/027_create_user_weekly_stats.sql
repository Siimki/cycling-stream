-- Create user_weekly_stats table to track weekly goal progress
CREATE TABLE IF NOT EXISTS user_weekly_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    week_number VARCHAR(10) NOT NULL, -- ISO week format: YYYY-WW (e.g., "2025-01")
    watch_minutes INTEGER NOT NULL DEFAULT 0,
    chat_messages INTEGER NOT NULL DEFAULT 0,
    weekly_goal_completed BOOLEAN NOT NULL DEFAULT false,
    weekly_reward_claimed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, week_number)
);

-- Create indexes for user_weekly_stats
CREATE INDEX idx_user_weekly_stats_user_week ON user_weekly_stats(user_id, week_number);
CREATE INDEX idx_user_weekly_stats_week ON user_weekly_stats(week_number);
CREATE INDEX idx_user_weekly_stats_user_id ON user_weekly_stats(user_id);

-- Add comment
COMMENT ON TABLE user_weekly_stats IS 'Tracks weekly goal progress (30 min watch + 3 messages) per user per week';


