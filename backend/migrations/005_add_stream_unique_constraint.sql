-- Add unique constraint on race_id for streams (one stream per race)
-- Check if constraint already exists before adding (idempotent)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'streams_race_id_unique'
    ) THEN
        ALTER TABLE streams ADD CONSTRAINT streams_race_id_unique UNIQUE (race_id);
    END IF;
END $$;

