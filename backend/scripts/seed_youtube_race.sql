INSERT INTO races (id, name, is_free, price_cents, start_date, created_at, updated_at)
VALUES ('EEC495D2-FB6A-4FCE-8CF9-E8E5EF8A2571', 'Youtube test race', true, 0, NOW(), NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO streams (race_id, status, stream_type, source_id, created_at, updated_at)
VALUES ('EEC495D2-FB6A-4FCE-8CF9-E8E5EF8A2571', 'live', 'youtube', 'jfKfPfyJRdk', NOW(), NOW())
ON CONFLICT (race_id) DO UPDATE
SET status = 'live', stream_type = 'youtube', source_id = 'jfKfPfyJRdk', updated_at = NOW();

