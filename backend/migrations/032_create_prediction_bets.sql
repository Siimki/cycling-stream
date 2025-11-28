-- Create prediction_bets table
CREATE TABLE IF NOT EXISTS prediction_bets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    market_id UUID NOT NULL REFERENCES prediction_markets(id) ON DELETE CASCADE,
    option_id UUID NOT NULL, -- Which option user bet on (references id in market.options JSONB)
    stake_points INTEGER NOT NULL, -- Points staked
    potential_payout INTEGER NOT NULL, -- stake * odds (calculated at bet time)
    result VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, won, lost
    payout_points INTEGER, -- Actual payout if won (same as potential_payout)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    settled_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for prediction_bets
CREATE INDEX idx_prediction_bets_user_id ON prediction_bets(user_id);
CREATE INDEX idx_prediction_bets_market_id ON prediction_bets(market_id);
CREATE INDEX idx_prediction_bets_result ON prediction_bets(result);
CREATE INDEX idx_prediction_bets_user_result ON prediction_bets(user_id, result);
CREATE INDEX idx_prediction_bets_market_result ON prediction_bets(market_id, result);

-- Add comment
COMMENT ON TABLE prediction_bets IS 'User bets on prediction markets';
COMMENT ON COLUMN prediction_bets.option_id IS 'UUID of the option bet on (matches id in market.options JSONB)';
COMMENT ON COLUMN prediction_bets.potential_payout IS 'Calculated payout (stake * odds) at time of bet';
COMMENT ON COLUMN prediction_bets.payout_points IS 'Actual payout if bet wins (set when market is settled)';


