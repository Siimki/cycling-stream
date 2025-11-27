# Database Migration Audit

## Overview
This document audits all database migrations for consistency, safety, and best practices.

## Migration Files Review

### ✅ Migration 001: Create Users
- **Status**: Good
- **Idempotent**: Yes (uses `IF NOT EXISTS`)
- **Foreign Keys**: N/A (base table)
- **Indexes**: 
  - `idx_users_email` - Good (unique email lookups)
- **Rollback**: Would need to drop table (not implemented)

### ✅ Migration 002: Create Races
- **Status**: Good
- **Idempotent**: Yes (uses `IF NOT EXISTS`)
- **Foreign Keys**: N/A (base table)
- **Indexes**: 
  - `idx_races_start_date` - Good (for date filtering)
  - `idx_races_category` - Good (for category filtering)
  - `idx_races_is_free` - Good (for free/paid filtering)
- **Rollback**: Would need to drop table (not implemented)

### ✅ Migration 003: Create Streams
- **Status**: Good
- **Idempotent**: Yes (uses `IF NOT EXISTS`)
- **Foreign Keys**: 
  - `race_id` → `races(id)` with `ON DELETE CASCADE` ✅
- **Indexes**: 
  - `idx_streams_race_id` - Good (FK lookup)
  - `idx_streams_status` - Good (status filtering)
- **Rollback**: Would need to drop table (not implemented)

### ✅ Migration 004: Create Entitlements
- **Status**: Good
- **Idempotent**: Yes (uses `IF NOT EXISTS`)
- **Foreign Keys**: 
  - `user_id` → `users(id)` with `ON DELETE CASCADE` ✅
  - `race_id` → `races(id)` with `ON DELETE CASCADE` ✅
- **Indexes**: 
  - `idx_entitlements_user_id` - Good (user entitlements lookup)
  - `idx_entitlements_race_id` - Good (race access checks)
  - `idx_entitlements_expires_at` - Good (expiration checks)
- **Unique Constraint**: `(user_id, race_id)` ✅ Prevents duplicates
- **Rollback**: Would need to drop table (not implemented)

### ✅ Migration 005: Add Stream Unique Constraint
- **Status**: Good
- **Idempotent**: No (would fail if constraint exists)
- **Issue**: Should use `IF NOT EXISTS` or check first
- **Recommendation**: Update to handle existing constraint gracefully
- **Rollback**: Would need to drop constraint (not implemented)

### ⚠️ Migration 006: Add Password Hash Index
- **Status**: Empty/No-op
- **Issue**: Contains only comments, no actual migration
- **Recommendation**: Remove this migration or add the index if needed
- **Note**: Email is used for login, so password_hash index may not be necessary

### ✅ Migration 007: Create Payments
- **Status**: Good
- **Idempotent**: Yes (uses `IF NOT EXISTS`)
- **Foreign Keys**: 
  - `user_id` → `users(id)` with `ON DELETE CASCADE` ✅
  - `race_id` → `races(id)` with `ON DELETE SET NULL` ✅ (preserves payment record if race deleted)
- **Indexes**: 
  - `idx_payments_user_id` - Good (user payment history)
  - `idx_payments_race_id` - Good (race revenue)
  - `idx_payments_status` - Good (status filtering)
  - `idx_payments_stripe_payment_intent_id` - Good (Stripe webhook lookups)
- **Unique Constraints**: 
  - `stripe_payment_intent_id` ✅
  - `stripe_checkout_session_id` ✅
- **Rollback**: Would need to drop table (not implemented)

### ✅ Migration 008: Create Watch Sessions
- **Status**: Good
- **Idempotent**: Yes (uses `IF NOT EXISTS`)
- **Foreign Keys**: 
  - `user_id` → `users(id)` with `ON DELETE CASCADE` ✅
  - `race_id` → `races(id)` with `ON DELETE CASCADE` ✅
- **Indexes**: 
  - `idx_watch_sessions_user_id` - Good
  - `idx_watch_sessions_race_id` - Good
  - `idx_watch_sessions_started_at` - Good (time-based queries)
  - `idx_watch_sessions_user_race` - Good (composite for user+race queries)
- **Views**: 
  - `watch_time_aggregated` - Good (pre-computed aggregations)
- **Rollback**: Would need to drop table and view (not implemented)

### ✅ Migration 009: Create Revenue Share Monthly
- **Status**: Good
- **Idempotent**: Yes (uses `IF NOT EXISTS`)
- **Foreign Keys**: 
  - `race_id` → `races(id)` with `ON DELETE CASCADE` ✅
- **Indexes**: 
  - `idx_revenue_share_monthly_race_id` - Good
  - `idx_revenue_share_monthly_year_month` - Good (time-based queries)
  - `idx_revenue_share_monthly_race_year_month` - Good (composite, may be redundant with unique constraint)
- **Unique Constraint**: `(race_id, year, month)` ✅ Prevents duplicates
- **Check Constraints**: `month >= 1 AND month <= 12` ✅
- **Views**: 
  - `revenue_share_details` - Good (joins with races)
- **Rollback**: Would need to drop table and view (not implemented)

### ✅ Migration 010: Create Viewer Sessions
- **Status**: Good
- **Idempotent**: Yes (uses `IF NOT EXISTS`)
- **Foreign Keys**: 
  - `user_id` → `users(id)` with `ON DELETE CASCADE` ✅ (nullable for anonymous)
  - `race_id` → `races(id)` with `ON DELETE CASCADE` ✅
- **Indexes**: 
  - `idx_viewer_sessions_race_id` - Good
  - `idx_viewer_sessions_user_id` - Good (nullable, but useful for authenticated viewers)
  - `idx_viewer_sessions_session_token` - Good (anonymous viewer tracking)
  - `idx_viewer_sessions_active` - Excellent (partial index for active sessions)
  - `idx_viewer_sessions_started_at` - Good (time-based queries)
- **Views**: 
  - `concurrent_viewers` - Good (real-time viewer counts)
  - `unique_viewers` - Good (all-time unique viewer counts)
- **Rollback**: Would need to drop table and views (not implemented)

### ✅ Migration 011: Create Costs
- **Status**: Good
- **Idempotent**: Yes (uses `IF NOT EXISTS`)
- **Foreign Keys**: 
  - `race_id` → `races(id)` with `ON DELETE CASCADE` ✅
- **Indexes**: 
  - `idx_costs_race_id` - Good
  - `idx_costs_year_month` - Good (time-based queries)
  - `idx_costs_race_year_month` - Good (composite)
  - `idx_costs_cost_type` - Good (type filtering)
  - `idx_costs_race_year_month_type` - Good (composite for detailed queries)
- **Check Constraints**: 
  - `cost_type IN (...)` ✅
  - `month >= 1 AND month <= 12` ✅
- **Views**: 
  - `cost_details` - Good (joins with races)
  - `cost_summary_monthly` - Good (monthly aggregations)
- **Rollback**: Would need to drop table and views (not implemented)

### ✅ Migration 012: Create Chat Messages
- **Status**: Good
- **Idempotent**: Yes (uses `IF NOT EXISTS`)
- **Foreign Keys**: 
  - `race_id` → `races(id)` with `ON DELETE CASCADE` ✅
  - `user_id` → `users(id)` with `ON DELETE SET NULL` ✅ (preserves messages if user deleted)
- **Indexes**: 
  - `idx_chat_messages_race_id` - Good
  - `idx_chat_messages_created_at` - Good (time-based queries)
  - `idx_chat_messages_race_created` - Excellent (composite for chat history)
  - `idx_chat_messages_user_id` - Good (user message history)
- **Check Constraints**: 
  - `char_length(message) > 0 AND char_length(message) <= 500` ✅
- **Rollback**: Would need to drop table (not implemented)

### ✅ Migration 013: Add Points to Users
- **Status**: Good
- **Idempotent**: Yes (uses `IF NOT EXISTS` and `ADD COLUMN IF NOT EXISTS`)
- **Indexes**: 
  - `idx_users_points` - Good (leaderboard queries)
- **Rollback**: Would need to drop column and index (not implemented)

### ✅ Migration 014: Add Bio to Users
- **Status**: Good
- **Idempotent**: Yes (uses `ADD COLUMN IF NOT EXISTS`)
- **Rollback**: Would need to drop column (not implemented)

### ✅ Migration 015: Add Stream Type to Streams
- **Status**: Good
- **Idempotent**: No (standard ALTER TABLE)
- **Indexes**: 
  - `idx_streams_stream_type` - Good (filtering by type)
- **Rollback**: Would need to drop columns and index (not implemented)

## Summary of Issues

### Critical Issues
None found.

### Minor Issues

1. **Migration 005**: Should handle existing constraint gracefully
   - **Fix**: Add check or use `IF NOT EXISTS` equivalent
   - **Impact**: Low (only affects re-running migrations)

2. **Migration 006**: Empty migration file
   - **Fix**: Remove or implement the intended index
   - **Impact**: Low (no-op migration)

### Recommendations

1. **Rollback Migrations**: Consider adding down migrations for all up migrations to enable safe rollbacks
   - Current state: Migrations are forward-only
   - Impact: Cannot rollback schema changes if issues are discovered

2. **Index Optimization**: 
   - `idx_revenue_share_monthly_race_year_month` may be redundant since `(race_id, year, month)` is unique
   - However, the unique constraint creates an index automatically, so the explicit index is redundant
   - **Recommendation**: Remove explicit index, rely on unique constraint index

3. **Migration Naming**: All migrations follow consistent naming pattern ✅

4. **Foreign Key Strategy**: 
   - All foreign keys properly defined ✅
   - `ON DELETE CASCADE` used appropriately for dependent data ✅
   - `ON DELETE SET NULL` used appropriately for optional references ✅

5. **Data Integrity**: 
   - Check constraints properly defined ✅
   - Unique constraints prevent duplicates ✅
   - NOT NULL constraints where appropriate ✅

## Best Practices Compliance

- ✅ Idempotent migrations (most use `IF NOT EXISTS`)
- ✅ Proper foreign key relationships
- ✅ Appropriate indexes for query patterns
- ✅ Check constraints for data validation
- ✅ Unique constraints where needed
- ⚠️ No rollback migrations (consider adding)
- ⚠️ One empty migration (006)

## Action Items

1. [ ] Fix Migration 005 to handle existing constraints
2. [ ] Remove or implement Migration 006
3. [ ] Consider removing redundant index in Migration 009
4. [ ] Consider adding down migrations for production safety
