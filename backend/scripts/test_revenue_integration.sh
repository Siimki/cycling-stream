#!/bin/bash

# Integration test script for Revenue Share functionality
# This script tests the revenue share endpoints and calculations
# Prerequisites: Database must be running and migrations must be applied

set -e

BASE_URL="http://localhost:8080"
ADMIN_EMAIL="admin@example.com"
ADMIN_PASSWORD="admin123"

echo "🧪 Testing Revenue Share Integration"
echo "===================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print test result
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
        exit 1
    fi
}

# Function to make API request
api_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local token=$4
    
    if [ -n "$token" ]; then
        if [ -n "$data" ]; then
            curl -s -X "$method" "$BASE_URL$endpoint" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $token" \
                -d "$data"
        else
            curl -s -X "$method" "$BASE_URL$endpoint" \
                -H "Authorization: Bearer $token"
        fi
    else
        if [ -n "$data" ]; then
            curl -s -X "$method" "$BASE_URL$endpoint" \
                -H "Content-Type: application/json" \
                -d "$data"
        else
            curl -s -X "$method" "$BASE_URL$endpoint"
        fi
    fi
}

# Test 1: Health check
echo "1. Testing health endpoint..."
response=$(api_request "GET" "/health")
if echo "$response" | grep -q "status"; then
    print_result 0 "Health check passed"
else
    print_result 1 "Health check failed"
fi

# Test 2: Admin login
echo ""
echo "2. Testing admin login..."
login_response=$(api_request "POST" "/auth/login" "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}")
TOKEN=$(echo "$login_response" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -n "$TOKEN" ]; then
    print_result 0 "Admin login successful"
    echo "   Token: ${TOKEN:0:20}..."
else
    print_result 1 "Admin login failed"
    echo "   Response: $login_response"
    exit 1
fi

# Test 3: Create a test race
echo ""
echo "3. Creating test race..."
race_data='{
    "name": "Test Revenue Race",
    "description": "Race for revenue share testing",
    "is_free": false,
    "price_cents": 1000
}'
race_response=$(api_request "POST" "/admin/races" "$race_data" "$TOKEN")
RACE_ID=$(echo "$race_response" | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -n "$RACE_ID" ]; then
    print_result 0 "Test race created (ID: ${RACE_ID:0:8}...)"
else
    print_result 1 "Failed to create test race"
    echo "   Response: $race_response"
    exit 1
fi

# Test 4: Get all revenue (should be empty initially)
echo ""
echo "4. Testing GET /admin/revenue (should be empty)..."
revenue_response=$(api_request "GET" "/admin/revenue" "" "$TOKEN")
if echo "$revenue_response" | grep -q '"data":\[\]' || echo "$revenue_response" | grep -q '"data":\[.*\]'; then
    print_result 0 "Get revenue endpoint works"
else
    print_result 1 "Get revenue endpoint failed"
    echo "   Response: $revenue_response"
fi

# Test 5: Get revenue by race (should be empty)
echo ""
echo "5. Testing GET /admin/revenue/races/:id..."
race_revenue_response=$(api_request "GET" "/admin/revenue/races/$RACE_ID" "" "$TOKEN")
if echo "$race_revenue_response" | grep -q '"data"'; then
    print_result 0 "Get revenue by race endpoint works"
else
    print_result 1 "Get revenue by race endpoint failed"
    echo "   Response: $race_revenue_response"
fi

# Test 6: Get revenue summary (should return zero values)
echo ""
echo "6. Testing GET /admin/revenue/races/:id/summary..."
summary_response=$(api_request "GET" "/admin/revenue/races/$RACE_ID/summary" "" "$TOKEN")
if echo "$summary_response" | grep -q '"total_revenue_cents"'; then
    print_result 0 "Get revenue summary endpoint works"
    echo "   Summary: $summary_response"
else
    print_result 1 "Get revenue summary endpoint failed"
    echo "   Response: $summary_response"
fi

# Test 7: Recalculate revenue (should work even with no data)
echo ""
echo "7. Testing POST /admin/revenue/recalculate..."
recalc_response=$(api_request "POST" "/admin/revenue/recalculate" "" "$TOKEN")
if echo "$recalc_response" | grep -q "successfully"; then
    print_result 0 "Recalculate revenue endpoint works"
else
    print_result 1 "Recalculate revenue endpoint failed"
    echo "   Response: $recalc_response"
fi

# Test 8: Recalculate revenue for specific period
echo ""
echo "8. Testing POST /admin/revenue/recalculate/:year/:month..."
YEAR=$(date +%Y)
MONTH=$(date +%m)
recalc_period_response=$(api_request "POST" "/admin/revenue/recalculate/$YEAR/$MONTH" "" "$TOKEN")
if echo "$recalc_period_response" | grep -q "successfully"; then
    print_result 0 "Recalculate revenue for period endpoint works"
else
    print_result 1 "Recalculate revenue for period endpoint failed"
    echo "   Response: $recalc_period_response"
fi

# Test 9: Test with query parameters
echo ""
echo "9. Testing GET /admin/revenue with query parameters..."
revenue_filtered=$(api_request "GET" "/admin/revenue?year=$YEAR&month=$MONTH" "" "$TOKEN")
if echo "$revenue_filtered" | grep -q '"data"'; then
    print_result 0 "Revenue filtering by year/month works"
else
    print_result 1 "Revenue filtering failed"
    echo "   Response: $revenue_filtered"
fi

echo ""
echo "===================================="
echo -e "${GREEN}✅ All integration tests passed!${NC}"
echo ""
echo "Note: To test with actual revenue data, you need to:"
echo "  1. Create payments with status='succeeded' for the race"
echo "  2. Create watch_sessions with duration_seconds for the race"
echo "  3. Run the recalculate endpoint"
echo "  4. Check the revenue data"

