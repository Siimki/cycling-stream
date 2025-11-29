-- Update stream status values from old system to new system
-- Old: planned, live, ended, offline
-- New: live, offline, upcoming

-- Step 1: Update all 'ended' statuses to 'offline'
UPDATE streams
SET status = 'offline', updated_at = CURRENT_TIMESTAMP
WHERE status = 'ended';

-- Step 2: Update all 'planned' statuses to 'upcoming'
UPDATE streams
SET status = 'upcoming', updated_at = CURRENT_TIMESTAMP
WHERE status = 'planned';

-- Step 3: Add CHECK constraint to enforce only valid status values
-- First, drop existing constraint if it exists (PostgreSQL doesn't have IF EXISTS for constraints)
DO $$
BEGIN
    -- Remove any existing check constraint on status
    ALTER TABLE streams DROP CONSTRAINT IF EXISTS streams_status_check;
END $$;

-- Add new constraint
ALTER TABLE streams
ADD CONSTRAINT streams_status_check CHECK (status IN ('live', 'offline', 'upcoming'));

-- Step 4: Update default value
ALTER TABLE streams
ALTER COLUMN status SET DEFAULT 'upcoming';

-- Step 5: Update column comment
COMMENT ON COLUMN streams.status IS 'live, offline, upcoming';

