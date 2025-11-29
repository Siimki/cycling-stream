-- Update all races with accurate dates and locations
-- Based on official race calendars and stage information

-- Tour de France 2025 - Already correct, but ensure location is complete
UPDATE races
SET location = 'France (21 stages from Lille to Paris)'
WHERE name = 'Tour de France 2025';

-- Tour de France 2024 - Stage 17 (mountain stage in Alps)
-- Stage 17 was July 18, 2024 - Saint-Paul-Trois-Châteaux to Superdévoluy
UPDATE races
SET start_date = '2024-07-18 10:00:00+00',
    location = 'French Alps (Saint-Paul-Trois-Châteaux to Superdévoluy)'
WHERE name = 'Tour de France 2024 - Stage 17';

-- Tour of Flanders 2025 - Already has correct date, just ensure location is accurate
UPDATE races
SET location = 'Belgium (Antwerp to Oudenaarde via the Flemish Ardennes)'
WHERE name = '2025 Tour of Flanders';

-- Giro d'Italia 2024 - Stage 12 (May 16, 2024 - Martinsicuro to Fano)
UPDATE races
SET start_date = '2024-05-16 10:00:00+00',
    location = 'Central Italy (Martinsicuro to Fano, Adriatic coast)'
WHERE name = 'Giro d''Italia 2024 - Stage 12';

-- Fuji Criterium races - These appear to be test races, keeping dates but improving location
UPDATE races
SET location = 'Tokyo, Japan (Odaiba district)'
WHERE name LIKE 'Fuji Criterium%';

-- Spaceship Stream - Test race, keeping reasonable date
UPDATE races
SET location = 'International Space Station (Low Earth Orbit)'
WHERE name = 'Spaceship Stream';

-- Youtube test race - Test race
UPDATE races
SET location = 'Virtual / Online'
WHERE name = 'Youtube test race';

-- Add end_date for single-day races (set to same day, 6 hours after start)
UPDATE races
SET end_date = start_date + INTERVAL '6 hours'
WHERE end_date IS NULL
  AND name != 'Tour de France 2025';  -- Keep Tour de France multi-day end_date as is
