-- Create races table
CREATE TABLE IF NOT EXISTS races (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    location VARCHAR(255),
    category VARCHAR(100),
    is_free BOOLEAN DEFAULT false,
    price_cents INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_races_start_date ON races(start_date);
CREATE INDEX idx_races_category ON races(category);
CREATE INDEX idx_races_is_free ON races(is_free);

