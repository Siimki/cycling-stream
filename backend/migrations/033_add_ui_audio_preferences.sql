-- Add UI and audio preference JSONB columns for granular motion/sound controls
ALTER TABLE user_preferences
    ADD COLUMN IF NOT EXISTS ui_preferences JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS audio_preferences JSONB NOT NULL DEFAULT '{}'::jsonb;

-- Seed sensible defaults for existing preference rows
UPDATE user_preferences
SET ui_preferences = jsonb_strip_nulls(jsonb_build_object(
        'chat_animations', true,
        'reduced_motion', false,
        'button_pulse', true,
        'poll_animations', true
    ))
WHERE (ui_preferences IS NULL OR ui_preferences = '{}'::jsonb);

UPDATE user_preferences
SET audio_preferences = jsonb_strip_nulls(jsonb_build_object(
        'button_clicks', true,
        'notification_sounds', true,
        'mention_pings', true,
        'master_volume', 0.15
    ))
WHERE (audio_preferences IS NULL OR audio_preferences = '{}'::jsonb);

