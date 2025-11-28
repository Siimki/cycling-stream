-- Add metadata columns for richer chat visuals
ALTER TABLE chat_messages
    ADD COLUMN IF NOT EXISTS user_role VARCHAR(32),
    ADD COLUMN IF NOT EXISTS badges JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS special_emote BOOLEAN NOT NULL DEFAULT FALSE;

