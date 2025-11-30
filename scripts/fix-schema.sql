-- Fix Races Table
ALTER TABLE races ADD COLUMN IF NOT EXISTS requires_login BOOLEAN DEFAULT false;
ALTER TABLE races ADD COLUMN IF NOT EXISTS stage_name VARCHAR(100);
ALTER TABLE races ADD COLUMN IF NOT EXISTS stage_type VARCHAR(50);
ALTER TABLE races ADD COLUMN IF NOT EXISTS elevation_meters INTEGER;
ALTER TABLE races ADD COLUMN IF NOT EXISTS estimated_finish_time TIME;
ALTER TABLE races ADD COLUMN IF NOT EXISTS stage_length_km INTEGER;

-- Fix Chat Messages Table
ALTER TABLE chat_messages ADD COLUMN IF NOT EXISTS user_role VARCHAR(32);
ALTER TABLE chat_messages ADD COLUMN IF NOT EXISTS badges JSONB NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE chat_messages ADD COLUMN IF NOT EXISTS special_emote BOOLEAN NOT NULL DEFAULT FALSE;

-- Create Achievements Tables if not exist (Migration 035)
CREATE TABLE IF NOT EXISTS achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    description TEXT,
    icon TEXT,
    points INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_id UUID NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
    unlocked_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    UNIQUE (user_id, achievement_id)
);

-- Create User Preferences if not exist (Migration 021)
CREATE TABLE IF NOT EXISTS user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    data_mode VARCHAR(20) NOT NULL DEFAULT 'standard' CHECK (data_mode IN ('casual', 'standard', 'pro')),
    preferred_units VARCHAR(20) NOT NULL DEFAULT 'metric' CHECK (preferred_units IN ('metric', 'imperial')),
    theme VARCHAR(20) NOT NULL DEFAULT 'auto' CHECK (theme IN ('light', 'dark', 'auto')),
    accent_color VARCHAR(50),
    device_type VARCHAR(20) CHECK (device_type IN ('tv', 'desktop', 'mobile', 'tablet')),
    notification_preferences JSONB DEFAULT '{}'::jsonb,
    onboarding_completed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for user_preferences
CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON user_preferences(user_id);

-- Check if trigger exists for user_preferences (conditionally drop to recreate or ignore)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_update_user_preferences_updated_at') THEN
        CREATE OR REPLACE FUNCTION update_user_preferences_updated_at()
        RETURNS TRIGGER AS $func$
        BEGIN
            NEW.updated_at = CURRENT_TIMESTAMP;
            RETURN NEW;
        END;
        $func$ LANGUAGE plpgsql;

        CREATE TRIGGER trigger_update_user_preferences_updated_at
            BEFORE UPDATE ON user_preferences
            FOR EACH ROW
            EXECUTE FUNCTION update_user_preferences_updated_at();
    END IF;
END
$$;

-- Fix Users Table (ensure bio and points exist)
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS points INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS xp_total INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS level INTEGER DEFAULT 1;

-- DATA CLEANUP: Ensure no NULLs in critical boolean/int columns that map to Go primitives
UPDATE races SET is_free = false WHERE is_free IS NULL;
UPDATE races SET price_cents = 0 WHERE price_cents IS NULL;
UPDATE races SET requires_login = false WHERE requires_login IS NULL;

