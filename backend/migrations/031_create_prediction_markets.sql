-- Create prediction_markets table
CREATE TABLE IF NOT EXISTS prediction_markets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    race_id UUID NOT NULL REFERENCES races(id) ON DELETE CASCADE,
    question VARCHAR(255) NOT NULL,
    options JSONB NOT NULL, -- Array of {id, text, odds}
    status VARCHAR(20) NOT NULL DEFAULT 'open', -- open, settled, cancelled
    settled_option_id UUID, -- Which option won (references id in options JSONB)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    settled_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for prediction_markets
CREATE INDEX idx_prediction_markets_race_id ON prediction_markets(race_id);
CREATE INDEX idx_prediction_markets_status ON prediction_markets(status);
CREATE INDEX idx_prediction_markets_race_status ON prediction_markets(race_id, status);

-- Add comment
COMMENT ON TABLE prediction_markets IS 'Prediction markets for races (e.g., "Will there be a breakaway?")';
COMMENT ON COLUMN prediction_markets.options IS 'JSONB array of prediction options with id, text, and odds';
COMMENT ON COLUMN prediction_markets.settled_option_id IS 'UUID of the winning option (matches id in options array)';


