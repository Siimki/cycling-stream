#!/bin/bash

# Test script for Analytics endpoints (Phase 7.2)
# Requires backend server to be running on localhost:8080
# Requires admin JWT token

set -e

BASE_URL="http://localhost:8080"
ADMIN_EMAIL="admin@cyclingstream.local"
ADMIN_PASSWORD="admin123"

echo "=== Testing Analytics Endpoints ==="
echo ""

# Step 1: Login as admin
echo "1. Logging in as admin..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}")

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "ERROR: Failed to get admin token"
  echo "Response: $LOGIN_RESPONSE"
  exit 1
fi

echo "✓ Admin token obtained"
echo ""

# Step 2: Test GET /admin/analytics/races
echo "2. Testing GET /admin/analytics/races..."
RACE_ANALYTICS=$(curl -s -X GET "$BASE_URL/admin/analytics/races" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "Response:"
echo "$RACE_ANALYTICS" | jq '.' 2>/dev/null || echo "$RACE_ANALYTICS"
echo ""

# Step 3: Test GET /admin/analytics/watch-time
echo "3. Testing GET /admin/analytics/watch-time..."
WATCH_TIME=$(curl -s -X GET "$BASE_URL/admin/analytics/watch-time" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "Response:"
echo "$WATCH_TIME" | jq '.' 2>/dev/null || echo "$WATCH_TIME"
echo ""

# Step 4: Test GET /admin/analytics/watch-time with year filter
echo "4. Testing GET /admin/analytics/watch-time?year=2024..."
WATCH_TIME_YEAR=$(curl -s -X GET "$BASE_URL/admin/analytics/watch-time?year=2024" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "Response:"
echo "$WATCH_TIME_YEAR" | jq '.' 2>/dev/null || echo "$WATCH_TIME_YEAR"
echo ""

# Step 5: Test GET /admin/analytics/revenue
echo "5. Testing GET /admin/analytics/revenue..."
REVENUE=$(curl -s -X GET "$BASE_URL/admin/analytics/revenue" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "Response:"
echo "$REVENUE" | jq '.' 2>/dev/null || echo "$REVENUE"
echo ""

# Step 6: Test unauthorized access
echo "6. Testing unauthorized access (should fail)..."
UNAUTHORIZED=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET "$BASE_URL/admin/analytics/races" \
  -H "Content-Type: application/json")

HTTP_CODE=$(echo "$UNAUTHORIZED" | grep "HTTP_CODE" | cut -d: -f2)
if [ "$HTTP_CODE" == "401" ] || [ "$HTTP_CODE" == "403" ]; then
  echo "✓ Unauthorized access correctly rejected (HTTP $HTTP_CODE)"
else
  echo "✗ Unauthorized access test failed (HTTP $HTTP_CODE)"
fi
echo ""

echo "=== Analytics Endpoints Test Complete ==="

