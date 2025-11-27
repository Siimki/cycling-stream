-- Add bio column to users table for short user descriptions
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio VARCHAR(120) NOT NULL DEFAULT '';

