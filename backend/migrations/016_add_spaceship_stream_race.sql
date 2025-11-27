-- Add Spaceship Stream race with YouTube stream
-- This migration creates a new race and associates a YouTube stream with it

DO $$
DECLARE
    race_id_val UUID;
BEGIN
    -- Insert the Spaceship Stream race
    INSERT INTO races (name, description, category, is_free, price_cents)
    VALUES (
        'Spaceship Stream',
        'Watch the amazing spaceship stream!',
        'Special',
        true,
        0
    )
    RETURNING id INTO race_id_val;

    -- Insert the YouTube stream for this race
    INSERT INTO streams (race_id, status, stream_type, source_id)
    VALUES (
        race_id_val,
        'live',
        'youtube',
        'fO9e9jnhYK8'
    );

    RAISE NOTICE 'Spaceship Stream race created with ID: %', race_id_val;
END $$;

