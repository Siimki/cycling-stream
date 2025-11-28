-- Seed career missions with tiers
-- Watch Time Career Missions
INSERT INTO missions (mission_type, title, description, points_reward, xp_reward, target_value, tier_number, category, is_active)
VALUES 
    ('watch_time', 'Watch for 180 Minutes (Tier 1)', 'Watch 180 minutes of live races', 500, 200, 180, 1, 'career', true),
    ('watch_time', 'Watch for 600 Minutes (Tier 2)', 'Watch 600 minutes of live races', 1000, 300, 600, 2, 'career', true),
    ('watch_time', 'Watch for 1500 Minutes (Tier 3)', 'Watch 1500 minutes of live races', 2000, 500, 1500, 3, 'career', true)
ON CONFLICT DO NOTHING;

-- Chat Messages Career Missions
INSERT INTO missions (mission_type, title, description, points_reward, xp_reward, target_value, tier_number, category, is_active)
VALUES 
    ('chat_message', 'Send 3 Messages (Tier 1)', 'Send 3 chat messages in live races', 100, 50, 3, 1, 'career', true),
    ('chat_message', 'Send 25 Messages (Tier 2)', 'Send 25 chat messages in live races', 300, 100, 25, 2, 'career', true),
    ('chat_message', 'Send 100 Messages (Tier 3)', 'Send 100 chat messages in live races', 800, 200, 100, 3, 'career', true)
ON CONFLICT DO NOTHING;

-- Predictions Career Missions
INSERT INTO missions (mission_type, title, description, points_reward, xp_reward, target_value, tier_number, category, is_active)
VALUES 
    ('predict_winner', 'Place 3 Predictions (Tier 1)', 'Place 3 predictions', 300, 100, 3, 1, 'career', true),
    ('predict_winner', 'Win 5 Predictions (Tier 2)', 'Win 5 predictions', 800, 200, 5, 2, 'career', true),
    ('predict_winner', 'Place 25 Predictions (Tier 3)', 'Place 25 predictions', 1500, 400, 25, 3, 'career', true)
ON CONFLICT DO NOTHING;


