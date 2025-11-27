-- Add requires_login field to races table
-- This field indicates whether a race stream requires user authentication to access

ALTER TABLE races
ADD COLUMN IF NOT EXISTS requires_login BOOLEAN DEFAULT false;

-- Add index for filtering races by login requirement
CREATE INDEX IF NOT EXISTS idx_races_requires_login ON races(requires_login);


