-- Add tier_number, xp_reward, category, and requirement_json columns to missions table
ALTER TABLE missions ADD COLUMN IF NOT EXISTS tier_number INTEGER DEFAULT 1;
ALTER TABLE missions ADD COLUMN IF NOT EXISTS xp_reward INTEGER DEFAULT 0;
ALTER TABLE missions ADD COLUMN IF NOT EXISTS category VARCHAR(50) DEFAULT 'career';
ALTER TABLE missions ADD COLUMN IF NOT EXISTS requirement_json JSONB;

-- Update existing missions to have category='career' and tier_number=1
UPDATE missions SET category = 'career', tier_number = 1 WHERE category IS NULL OR category = '';

-- Create index on category for faster queries
CREATE INDEX IF NOT EXISTS idx_missions_category ON missions(category);
CREATE INDEX IF NOT EXISTS idx_missions_category_type_tier ON missions(category, mission_type, tier_number);

-- Add comment
COMMENT ON COLUMN missions.tier_number IS 'Tier number for career missions (1, 2, 3, etc.)';
COMMENT ON COLUMN missions.xp_reward IS 'XP reward for completing this mission';
COMMENT ON COLUMN missions.category IS 'Mission category: career or weekly';
COMMENT ON COLUMN missions.requirement_json IS 'Flexible requirement data in JSON format';


