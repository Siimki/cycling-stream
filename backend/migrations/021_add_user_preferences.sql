-- Create user_preferences table for personalization settings
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

CREATE INDEX idx_user_preferences_user_id ON user_preferences(user_id);

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_user_preferences_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_user_preferences_updated_at
    BEFORE UPDATE ON user_preferences
    FOR EACH ROW
    EXECUTE FUNCTION update_user_preferences_updated_at();

