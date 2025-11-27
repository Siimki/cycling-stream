-- Create entitlements table
CREATE TABLE IF NOT EXISTS entitlements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    race_id UUID NOT NULL REFERENCES races(id) ON DELETE CASCADE,
    type VARCHAR(50) DEFAULT 'ticket', -- ticket, subscription
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, race_id)
);

CREATE INDEX idx_entitlements_user_id ON entitlements(user_id);
CREATE INDEX idx_entitlements_race_id ON entitlements(race_id);
CREATE INDEX idx_entitlements_expires_at ON entitlements(expires_at);

