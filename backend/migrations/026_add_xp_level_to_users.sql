-- Add XP, Level, and best_streak_weeks columns to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS xp_total INTEGER NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS level INTEGER NOT NULL DEFAULT 1;
ALTER TABLE users ADD COLUMN IF NOT EXISTS best_streak_weeks INTEGER NOT NULL DEFAULT 0;

-- Create indexes for XP and Level (useful for leaderboards and queries)
CREATE INDEX IF NOT EXISTS idx_users_xp ON users(xp_total DESC);
CREATE INDEX IF NOT EXISTS idx_users_level ON users(level DESC);

-- Add comment
COMMENT ON COLUMN users.xp_total IS 'Total experience points accumulated by the user';
COMMENT ON COLUMN users.level IS 'Current level based on XP total';
COMMENT ON COLUMN users.best_streak_weeks IS 'Best weekly streak achieved (cosmetic stat)';


