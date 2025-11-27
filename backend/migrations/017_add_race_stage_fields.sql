-- Add stage-specific fields to races table
-- These fields support displaying detailed race information in the RaceStats component

ALTER TABLE races
ADD COLUMN IF NOT EXISTS stage_name VARCHAR(100),
ADD COLUMN IF NOT EXISTS stage_type VARCHAR(50),
ADD COLUMN IF NOT EXISTS elevation_meters INTEGER,
ADD COLUMN IF NOT EXISTS estimated_finish_time TIME,
ADD COLUMN IF NOT EXISTS stage_length_km INTEGER;

-- Add index on stage_type for filtering
CREATE INDEX IF NOT EXISTS idx_races_stage_type ON races(stage_type);

