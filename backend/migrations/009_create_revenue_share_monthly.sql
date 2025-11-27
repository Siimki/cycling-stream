-- Create revenue_share_monthly table to store monthly revenue aggregations
CREATE TABLE IF NOT EXISTS revenue_share_monthly (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    race_id UUID NOT NULL REFERENCES races(id) ON DELETE CASCADE,
    year INTEGER NOT NULL,
    month INTEGER NOT NULL CHECK (month >= 1 AND month <= 12),
    total_revenue_cents INTEGER NOT NULL DEFAULT 0,
    total_watch_minutes DECIMAL(10, 2) NOT NULL DEFAULT 0,
    platform_share_cents INTEGER NOT NULL DEFAULT 0,
    organizer_share_cents INTEGER NOT NULL DEFAULT 0,
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(race_id, year, month)
);

CREATE INDEX idx_revenue_share_monthly_race_id ON revenue_share_monthly(race_id);
CREATE INDEX idx_revenue_share_monthly_year_month ON revenue_share_monthly(year, month);
-- Note: idx_revenue_share_monthly_race_year_month is redundant because
-- the UNIQUE constraint on (race_id, year, month) automatically creates an index.
-- However, we keep it explicit for clarity and to ensure the index exists even if
-- the unique constraint is modified in the future.
CREATE INDEX IF NOT EXISTS idx_revenue_share_monthly_race_year_month ON revenue_share_monthly(race_id, year, month);

-- Create view for revenue share details with race information
CREATE OR REPLACE VIEW revenue_share_details AS
SELECT 
    rsm.id,
    rsm.race_id,
    r.name as race_name,
    rsm.year,
    rsm.month,
    rsm.total_revenue_cents,
    rsm.total_revenue_cents / 100.0 as total_revenue_dollars,
    rsm.total_watch_minutes,
    rsm.platform_share_cents,
    rsm.platform_share_cents / 100.0 as platform_share_dollars,
    rsm.organizer_share_cents,
    rsm.organizer_share_cents / 100.0 as organizer_share_dollars,
    rsm.calculated_at,
    rsm.created_at,
    rsm.updated_at
FROM revenue_share_monthly rsm
JOIN races r ON r.id = rsm.race_id;

