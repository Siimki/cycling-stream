-- Update existing races with realistic stage data
-- This migration populates the new stage fields with realistic Tour de France data

-- Update "Youtube test race" with Tour de France 2025 Stage 17 data
UPDATE races
SET 
    name = 'Tour de France 2025',
    stage_name = 'Stage 17',
    stage_type = 'Mountain',
    elevation_meters = 4800,
    estimated_finish_time = '17:45:00',
    stage_length_km = 166,
    updated_at = CURRENT_TIMESTAMP
WHERE id = 'EEC495D2-FB6A-4FCE-8CF9-E8E5EF8A2571';

-- Update "Spaceship Stream" with special stage data
UPDATE races
SET 
    stage_name = 'Special Stage',
    stage_type = 'Special',
    elevation_meters = 0,
    estimated_finish_time = '18:00:00',
    stage_length_km = 120,
    updated_at = CURRENT_TIMESTAMP
WHERE name = 'Spaceship Stream';

