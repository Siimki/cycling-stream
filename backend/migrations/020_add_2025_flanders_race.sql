-- Add 2025 Tour of Flanders race with YouTube stream
-- This race requires user login to access the stream

INSERT INTO races (
    id,
    name,
    description,
    start_date,
    location,
    category,
    is_free,
    price_cents,
    requires_login,
    stage_length_km,
    elevation_meters,
    created_at,
    updated_at
) VALUES (
    gen_random_uuid(),
    '2025 Tour of Flanders',
    'The 2025 Tour of Flanders, one of the five monuments of cycling, covering 269 km through the challenging cobbled climbs of Belgium.',
    '2025-04-06 10:00:00+00'::timestamptz,
    'Belgium',
    'Classic',
    true,
    0,
    true,
    269,
    1769,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

-- Create associated YouTube stream for the race
INSERT INTO streams (
    id,
    race_id,
    status,
    stream_type,
    source_id,
    created_at,
    updated_at
)
SELECT 
    gen_random_uuid(),
    r.id,
    'live',
    'youtube',
    'T1O7zxajgcU',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
FROM races r
WHERE r.name = '2025 Tour of Flanders'
LIMIT 1;


