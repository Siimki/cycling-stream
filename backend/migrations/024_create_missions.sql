-- Create mission_type enum
CREATE TYPE mission_type AS ENUM (
    'watch_time',
    'chat_message',
    'watch_race',
    'follow_series',
    'streak',
    'predict_winner'
);

-- Create missions table
CREATE TABLE IF NOT EXISTS missions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mission_type mission_type NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    points_reward INTEGER NOT NULL DEFAULT 0,
    target_value INTEGER NOT NULL DEFAULT 1,
    valid_from TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    valid_until TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for missions
CREATE INDEX idx_missions_mission_type ON missions(mission_type);
CREATE INDEX idx_missions_is_active ON missions(is_active);
CREATE INDEX idx_missions_valid_dates ON missions(valid_from, valid_until) WHERE is_active = true;

-- Add comment
COMMENT ON TABLE missions IS 'Available missions that users can complete to earn points';

