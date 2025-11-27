-- Migration 006: Add index on password_hash
-- NOTE: This migration is intentionally empty because:
-- 1. We use email for login lookups, not password_hash
-- 2. Password hashes are looked up via email index (idx_users_email)
-- 3. Adding an index on password_hash would not improve query performance
-- 
-- If password_hash lookups are needed in the future, uncomment below:
-- CREATE INDEX IF NOT EXISTS idx_users_password_hash ON users(password_hash);

