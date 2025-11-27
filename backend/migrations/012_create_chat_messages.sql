-- Create chat_messages table for live chat functionality
CREATE TABLE IF NOT EXISTS chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    race_id UUID NOT NULL REFERENCES races(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    username VARCHAR(255) NOT NULL,
    message TEXT NOT NULL CHECK (char_length(message) > 0 AND char_length(message) <= 500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_chat_messages_race_id ON chat_messages(race_id);
CREATE INDEX idx_chat_messages_created_at ON chat_messages(created_at);
CREATE INDEX idx_chat_messages_race_created ON chat_messages(race_id, created_at DESC);
CREATE INDEX idx_chat_messages_user_id ON chat_messages(user_id);

