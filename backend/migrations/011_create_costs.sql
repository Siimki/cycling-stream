-- Create costs table to track platform costs (CDN, server, storage, etc.)
CREATE TABLE IF NOT EXISTS costs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    race_id UUID REFERENCES races(id) ON DELETE CASCADE,
    cost_type VARCHAR(50) NOT NULL CHECK (cost_type IN ('cdn', 'server', 'storage', 'bandwidth', 'other')),
    amount_cents INTEGER NOT NULL,
    year INTEGER NOT NULL,
    month INTEGER NOT NULL CHECK (month >= 1 AND month <= 12),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_costs_race_id ON costs(race_id);
CREATE INDEX idx_costs_year_month ON costs(year, month);
CREATE INDEX idx_costs_race_year_month ON costs(race_id, year, month);
CREATE INDEX idx_costs_cost_type ON costs(cost_type);
CREATE INDEX idx_costs_race_year_month_type ON costs(race_id, year, month, cost_type);

-- Create view for cost details with race information
CREATE OR REPLACE VIEW cost_details AS
SELECT 
    c.id,
    c.race_id,
    r.name as race_name,
    c.cost_type,
    c.amount_cents,
    c.amount_cents / 100.0 as amount_dollars,
    c.year,
    c.month,
    c.description,
    c.created_at,
    c.updated_at
FROM costs c
LEFT JOIN races r ON r.id = c.race_id;

-- Create view for monthly cost aggregation by race
CREATE OR REPLACE VIEW cost_summary_monthly AS
SELECT 
    race_id,
    year,
    month,
    SUM(CASE WHEN cost_type = 'cdn' THEN amount_cents ELSE 0 END) as cdn_cents,
    SUM(CASE WHEN cost_type = 'server' THEN amount_cents ELSE 0 END) as server_cents,
    SUM(CASE WHEN cost_type = 'storage' THEN amount_cents ELSE 0 END) as storage_cents,
    SUM(CASE WHEN cost_type = 'bandwidth' THEN amount_cents ELSE 0 END) as bandwidth_cents,
    SUM(CASE WHEN cost_type = 'other' THEN amount_cents ELSE 0 END) as other_cents,
    SUM(amount_cents) as total_cents,
    SUM(amount_cents) / 100.0 as total_dollars
FROM costs
GROUP BY race_id, year, month;

