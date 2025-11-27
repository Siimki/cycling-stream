-- Add points column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS points INTEGER NOT NULL DEFAULT 0;

-- Create index on points for potential leaderboard queries
CREATE INDEX IF NOT EXISTS idx_users_points ON users(points DESC);

