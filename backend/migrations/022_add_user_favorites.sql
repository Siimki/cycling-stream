-- Create user_favorites table for tracking user favorites (riders, teams, races, series)
CREATE TABLE IF NOT EXISTS user_favorites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    favorite_type VARCHAR(20) NOT NULL CHECK (favorite_type IN ('rider', 'team', 'race', 'series')),
    favorite_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, favorite_type, favorite_id)
);

CREATE INDEX idx_user_favorites_user_id ON user_favorites(user_id);
CREATE INDEX idx_user_favorites_type ON user_favorites(favorite_type);
CREATE INDEX idx_user_favorites_user_type ON user_favorites(user_id, favorite_type);

