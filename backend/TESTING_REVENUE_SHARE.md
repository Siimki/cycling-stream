# Testing Revenue Share (Phase 6.4)

This document describes how to test the revenue share functionality.

## Prerequisites

1. Database is running (`make docker-up`)
2. Migrations are up to date (`make migrate-up`)
3. Backend server is running (`make run-backend`)
4. Admin JWT token (obtained from `/auth/login` with admin credentials)

## Test Data Setup

Before testing, you need to set up test data:

1. **Create a race** (via admin endpoint):
```bash
POST /admin/races
{
  "name": "Test Race",
  "description": "Test race for revenue share",
  "is_free": false,
  "price_cents": 1000
}
```

2. **Create payments** (via payment endpoint or directly in database):
   - Create payments with `status = 'succeeded'` for the race
   - Payments should be in different months to test monthly aggregation

3. **Create watch sessions** (via watch tracking):
   - Users should watch the race
   - Watch sessions should have `duration_seconds` set
   - Sessions should be in the same months as payments

## Testing Revenue Calculation

### 1. Recalculate All Monthly Revenue

This will calculate revenue for all races with payments:

```bash
POST /admin/revenue/recalculate
Authorization: Bearer <admin_jwt_token>
```

**Expected Response:**
```json
{
  "message": "Revenue data recalculated successfully"
}
```

### 2. Recalculate Revenue for Specific Period

Calculate revenue for a specific year and month:

```bash
POST /admin/revenue/recalculate/2024/11
Authorization: Bearer <admin_jwt_token>
```

**Expected Response:**
```json
{
  "message": "Revenue data recalculated successfully"
}
```

### 3. Get All Revenue Data

Get all monthly revenue data:

```bash
GET /admin/revenue
Authorization: Bearer <admin_jwt_token>
```

**Query Parameters (optional):**
- `year`: Filter by year (e.g., `?year=2024`)
- `month`: Filter by month (e.g., `?month=11` or `?year=2024&month=11`)

**Expected Response:**
```json
{
  "data": [
    {
      "id": "...",
      "race_id": "...",
      "race_name": "Test Race",
      "year": 2024,
      "month": 11,
      "total_revenue_cents": 1000,
      "total_revenue_dollars": 10.0,
      "total_watch_minutes": 120.5,
      "platform_share_cents": 500,
      "platform_share_dollars": 5.0,
      "organizer_share_cents": 500,
      "organizer_share_dollars": 5.0,
      "calculated_at": "...",
      "created_at": "...",
      "updated_at": "..."
    }
  ]
}
```

### 4. Get Revenue by Race

Get monthly revenue data for a specific race:

```bash
GET /admin/revenue/races/{race_id}
Authorization: Bearer <admin_jwt_token>
```

**Expected Response:**
```json
{
  "data": [
    {
      "id": "...",
      "race_id": "...",
      "race_name": "Test Race",
      "year": 2024,
      "month": 11,
      ...
    }
  ]
}
```

### 5. Get Revenue Summary by Race

Get aggregated revenue summary for a specific race:

```bash
GET /admin/revenue/races/{race_id}/summary
Authorization: Bearer <admin_jwt_token>
```

**Expected Response:**
```json
{
  "race_id": "...",
  "race_name": "Test Race",
  "total_revenue_cents": 5000,
  "total_revenue_dollars": 50.0,
  "total_watch_minutes": 600.0,
  "platform_share_cents": 2500,
  "platform_share_dollars": 25.0,
  "organizer_share_cents": 2500,
  "organizer_share_dollars": 25.0,
  "month_count": 3
}
```

## Revenue Split Calculation

The revenue is split 50/50 between platform and organizer:
- **Platform Share**: `total_revenue_cents / 2`
- **Organizer Share**: `total_revenue_cents - platform_share_cents` (handles odd cents)

**Examples:**
- $10.00 (1000 cents) → Platform: $5.00, Organizer: $5.00
- $10.01 (1001 cents) → Platform: $5.00, Organizer: $5.01
- $0.01 (1 cent) → Platform: $0.00, Organizer: $0.01

## Test Scenarios

### Scenario 1: Single Payment, Single Month
1. Create a race with price $10.00
2. Create a payment of $10.00 for November 2024
3. Create watch sessions totaling 60 minutes in November 2024
4. Recalculate revenue for November 2024
5. Verify:
   - Total revenue: $10.00
   - Platform share: $5.00
   - Organizer share: $5.00
   - Watch minutes: 60.0

### Scenario 2: Multiple Payments, Multiple Months
1. Create a race
2. Create payments in different months (e.g., November and December 2024)
3. Create watch sessions in corresponding months
4. Recalculate all revenue
5. Verify monthly breakdowns are correct
6. Verify summary aggregates correctly

### Scenario 3: Odd Amount Revenue Split
1. Create a payment of $10.01 (1001 cents)
2. Recalculate revenue
3. Verify:
   - Platform share: $5.00 (500 cents)
   - Organizer share: $5.01 (501 cents)
   - Total: $10.01

### Scenario 4: No Revenue Data
1. Get revenue for a race with no payments
2. Verify empty array or zero values are returned appropriately

### Scenario 5: Filter by Year/Month
1. Create revenue data for multiple months
2. Filter by specific year and month
3. Verify only matching records are returned

## Database Verification

You can also verify the data directly in the database:

```sql
-- View all revenue share data
SELECT * FROM revenue_share_details ORDER BY year DESC, month DESC;

-- View revenue for a specific race
SELECT * FROM revenue_share_details WHERE race_id = '<race_id>';

-- Verify revenue calculation
SELECT 
    rsm.race_id,
    r.name,
    rsm.year,
    rsm.month,
    rsm.total_revenue_cents,
    rsm.platform_share_cents,
    rsm.organizer_share_cents,
    (rsm.platform_share_cents + rsm.organizer_share_cents) as total_check
FROM revenue_share_monthly rsm
JOIN races r ON r.id = rsm.race_id;
```

## Error Cases to Test

1. **Invalid race ID**: Should return 404 or empty data
2. **Invalid year/month**: Should return 400 Bad Request
3. **Unauthorized access**: Should return 401 Unauthorized
4. **Missing admin token**: Should return 401 Unauthorized

## Automated Tests

Run the unit tests:

```bash
cd backend
go test -v ./internal/repository -run TestRevenueRepository_RevenueSplit
```

This tests the revenue split calculation logic (50/50 split with proper handling of odd cents).

