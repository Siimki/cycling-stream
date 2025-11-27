#!/bin/bash

# Comprehensive API Test Suite for CyclingStream Platform
# Tests all critical backend endpoints
# Requires: backend server running on localhost:8080, database running, jq installed

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
ADMIN_EMAIL="${ADMIN_EMAIL:-admin@cyclingstream.local}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-admin123}"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0
TESTS_TOTAL=0

# Helper functions
pass_test() {
    ((TESTS_PASSED++))
    ((TESTS_TOTAL++))
    echo -e "${GREEN}✓${NC} $1"
}

fail_test() {
    ((TESTS_FAILED++))
    ((TESTS_TOTAL++))
    echo -e "${RED}✗${NC} $1"
    if [ -n "$2" ]; then
        echo -e "  ${YELLOW}Details: $2${NC}"
    fi
}

info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

# Check if server is running
check_server() {
    if ! curl -s "$BASE_URL/health" > /dev/null 2>&1; then
        echo -e "${RED}Error: Backend server is not running on $BASE_URL${NC}"
        echo "Please start the server with: make run-backend"
        exit 1
    fi
}

# Get admin token
get_admin_token() {
    local response=$(curl -s -X POST "$BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}")
    
    local token=$(echo "$response" | jq -r '.token // empty' 2>/dev/null)
    if [ -z "$token" ] || [ "$token" = "null" ]; then
        echo ""
        return 1
    fi
    echo "$token"
}

# Get user token (register new user)
get_user_token() {
    local email="testuser$(date +%s)@test.com"
    local password="Test123!@#"
    
    # Register
    local register_response=$(curl -s -X POST "$BASE_URL/auth/register" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$email\",\"password\":\"$password\"}")
    
    # Login
    local login_response=$(curl -s -X POST "$BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$email\",\"password\":\"$password\"}")
    
    local token=$(echo "$login_response" | jq -r '.token // empty' 2>/dev/null)
    if [ -z "$token" ] || [ "$token" = "null" ]; then
        echo ""
        return 1
    fi
    echo "$token"
}

echo "=========================================="
echo "CyclingStream API Test Suite"
echo "=========================================="
echo ""

check_server

# Get tokens
info "Getting admin token..."
ADMIN_TOKEN=$(get_admin_token)
if [ -z "$ADMIN_TOKEN" ]; then
    fail_test "Failed to get admin token" "Check admin credentials"
    exit 1
fi
pass_test "Admin token obtained"

info "Creating test user..."
USER_TOKEN=$(get_user_token)
if [ -z "$USER_TOKEN" ]; then
    fail_test "Failed to create test user"
    exit 1
fi
pass_test "Test user created and token obtained"

echo ""
echo "=========================================="
echo "1. Health & System Tests"
echo "=========================================="

# Test 1.1: Health check
response=$(curl -s "$BASE_URL/health")
if echo "$response" | jq -e '.status' > /dev/null 2>&1; then
    pass_test "GET /health returns valid response"
else
    fail_test "GET /health" "Invalid response format"
fi

# Test 1.2: Health check structure
if echo "$response" | jq -e '.services.database' > /dev/null 2>&1; then
    pass_test "GET /health includes database status"
else
    fail_test "GET /health missing database status"
fi

echo ""
echo "=========================================="
echo "2. Public Race Endpoints"
echo "=========================================="

# Test 2.1: List races
response=$(curl -s "$BASE_URL/races")
if echo "$response" | jq -e 'type == "array"' > /dev/null 2>&1; then
    pass_test "GET /races returns array"
else
    fail_test "GET /races" "Invalid response format"
fi

# Test 2.2: Get race by ID (if races exist)
race_id=$(echo "$response" | jq -r '.[0].id // empty' 2>/dev/null)
if [ -n "$race_id" ] && [ "$race_id" != "null" ]; then
    response=$(curl -s "$BASE_URL/races/$race_id")
    if echo "$response" | jq -e '.id' > /dev/null 2>&1; then
        pass_test "GET /races/:id returns race details"
    else
        fail_test "GET /races/:id" "Invalid response"
    fi
    
    # Test 2.3: Stream status
    response=$(curl -s "$BASE_URL/races/$race_id/stream/status")
    if echo "$response" | jq -e '.status' > /dev/null 2>&1; then
        pass_test "GET /races/:id/stream/status returns status"
    else
        fail_test "GET /races/:id/stream/status" "Invalid response"
    fi
else
    info "No races found, skipping race detail tests"
fi

echo ""
echo "=========================================="
echo "3. Authentication Tests"
echo "=========================================="

# Test 3.1: Register new user
test_email="newuser$(date +%s)@test.com"
response=$(curl -s -X POST "$BASE_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$test_email\",\"password\":\"Test123!@#\"}")
if echo "$response" | jq -e '.token' > /dev/null 2>&1; then
    pass_test "POST /auth/register creates user and returns token"
else
    fail_test "POST /auth/register" "Failed to register user"
fi

# Test 3.2: Register duplicate email
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$test_email\",\"password\":\"Test123!@#\"}")
http_code=$(echo "$response" | tail -n1)
if [ "$http_code" = "400" ] || [ "$http_code" = "409" ]; then
    pass_test "POST /auth/register rejects duplicate email"
else
    fail_test "POST /auth/register duplicate email" "Expected 400/409, got $http_code"
fi

# Test 3.3: Login
response=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$test_email\",\"password\":\"Test123!@#\"}")
if echo "$response" | jq -e '.token' > /dev/null 2>&1; then
    pass_test "POST /auth/login returns token"
else
    fail_test "POST /auth/login" "Failed to login"
fi

# Test 3.4: Login with wrong password
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$test_email\",\"password\":\"WrongPass123!\"}")
http_code=$(echo "$response" | tail -n1)
if [ "$http_code" = "401" ]; then
    pass_test "POST /auth/login rejects wrong password"
else
    fail_test "POST /auth/login wrong password" "Expected 401, got $http_code"
fi

# Test 3.5: Get user profile
response=$(curl -s -X GET "$BASE_URL/users/me" \
    -H "Authorization: Bearer $USER_TOKEN")
if echo "$response" | jq -e '.email' > /dev/null 2>&1; then
    pass_test "GET /users/me returns user profile"
else
    fail_test "GET /users/me" "Failed to get profile"
fi

# Test 3.6: Get profile without auth
response=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/users/me")
http_code=$(echo "$response" | tail -n1)
if [ "$http_code" = "401" ]; then
    pass_test "GET /users/me requires authentication"
else
    fail_test "GET /users/me auth check" "Expected 401, got $http_code"
fi

echo ""
echo "=========================================="
echo "4. Admin Race Management"
echo "=========================================="

# Test 4.1: Create race
race_data="{\"name\":\"Test Race $(date +%s)\",\"description\":\"Test race description\",\"date\":\"2024-12-31T12:00:00Z\",\"location\":\"Test Location\",\"price_cents\":1000,\"is_free\":false}"
response=$(curl -s -X POST "$BASE_URL/admin/races" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d "$race_data")
created_race_id=$(echo "$response" | jq -r '.id // empty' 2>/dev/null)
if [ -n "$created_race_id" ] && [ "$created_race_id" != "null" ]; then
    pass_test "POST /admin/races creates race"
else
    fail_test "POST /admin/races" "Failed to create race"
    created_race_id=""
fi

if [ -n "$created_race_id" ]; then
    # Test 4.2: Update race
    update_data="{\"name\":\"Updated Test Race\",\"description\":\"Updated description\"}"
    response=$(curl -s -X PUT "$BASE_URL/admin/races/$created_race_id" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "$update_data")
    if echo "$response" | jq -e '.name == "Updated Test Race"' > /dev/null 2>&1; then
        pass_test "PUT /admin/races/:id updates race"
    else
        fail_test "PUT /admin/races/:id" "Failed to update race"
    fi
    
    # Test 4.3: Add stream to race
    stream_data="{\"origin_url\":\"https://example.com/stream.m3u8\",\"cdn_url\":\"https://cdn.example.com/stream.m3u8\",\"status\":\"live\"}"
    response=$(curl -s -X POST "$BASE_URL/admin/races/$created_race_id/stream" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "$stream_data")
    if echo "$response" | jq -e '.origin_url' > /dev/null 2>&1; then
        pass_test "POST /admin/races/:id/stream adds stream"
    else
        fail_test "POST /admin/races/:id/stream" "Failed to add stream"
    fi
    
    # Test 4.4: Update stream status
    status_data="{\"status\":\"live\"}"
    response=$(curl -s -X PUT "$BASE_URL/admin/races/$created_race_id/stream/status" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "$status_data")
    if echo "$response" | jq -e '.message' > /dev/null 2>&1 || echo "$response" | jq -e '.status == "live"' > /dev/null 2>&1; then
        pass_test "PUT /admin/races/:id/stream/status updates status"
    else
        fail_test "PUT /admin/races/:id/stream/status" "Failed to update status: $response"
    fi
    
    # Test 4.5: Delete race
    response=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_URL/admin/races/$created_race_id" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    http_code=$(echo "$response" | tail -n1)
    if [ "$http_code" = "200" ] || [ "$http_code" = "204" ]; then
        pass_test "DELETE /admin/races/:id deletes race"
    else
        fail_test "DELETE /admin/races/:id" "Expected 200/204, got $http_code"
    fi
fi

# Test 4.6: Admin endpoint requires auth
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/admin/races" \
    -H "Content-Type: application/json" \
    -d "$race_data")
http_code=$(echo "$response" | tail -n1)
if [ "$http_code" = "401" ] || [ "$http_code" = "403" ]; then
    pass_test "POST /admin/races requires admin authentication"
else
    fail_test "POST /admin/races auth check" "Expected 401/403, got $http_code"
fi

echo ""
echo "=========================================="
echo "5. Viewer Tracking"
echo "=========================================="

if [ -n "$race_id" ] && [ "$race_id" != "null" ]; then
    # Test 5.1: Start viewer session (anonymous)
    response=$(curl -s -X POST "$BASE_URL/viewers/sessions/start" \
        -H "Content-Type: application/json" \
        -d "{\"race_id\":\"$race_id\"}")
    session_token=$(echo "$response" | jq -r '.session_token // empty' 2>/dev/null)
    if [ -n "$session_token" ] && [ "$session_token" != "null" ]; then
        pass_test "POST /viewers/sessions/start creates anonymous session"
        
        # Test 5.2: Heartbeat
        response=$(curl -s -X POST "$BASE_URL/viewers/sessions/heartbeat" \
            -H "Content-Type: application/json" \
            -d "{\"session_token\":\"$session_token\"}")
        if echo "$response" | jq -e '.success' > /dev/null 2>&1; then
            pass_test "POST /viewers/sessions/heartbeat updates session"
        else
            fail_test "POST /viewers/sessions/heartbeat" "Failed to update heartbeat"
        fi
        
        # Test 5.3: End session
        response=$(curl -s -X POST "$BASE_URL/viewers/sessions/end" \
            -H "Content-Type: application/json" \
            -d "{\"session_token\":\"$session_token\"}")
        if echo "$response" | jq -e '.success' > /dev/null 2>&1; then
            pass_test "POST /viewers/sessions/end ends session"
        else
            fail_test "POST /viewers/sessions/end" "Failed to end session"
        fi
    else
        fail_test "POST /viewers/sessions/start" "Failed to create session"
    fi
    
    # Test 5.4: Get concurrent viewers
    response=$(curl -s "$BASE_URL/races/$race_id/viewers/concurrent")
    if echo "$response" | jq -e '.concurrent_viewers' > /dev/null 2>&1; then
        pass_test "GET /races/:id/viewers/concurrent returns count"
    else
        fail_test "GET /races/:id/viewers/concurrent" "Invalid response"
    fi
    
    # Test 5.5: Get unique viewers
    response=$(curl -s "$BASE_URL/races/$race_id/viewers/unique")
    if echo "$response" | jq -e '.unique_viewers' > /dev/null 2>&1; then
        pass_test "GET /races/:id/viewers/unique returns count"
    else
        fail_test "GET /races/:id/viewers/unique" "Invalid response"
    fi
else
    info "No race ID available, skipping viewer tracking tests"
fi

echo ""
echo "=========================================="
echo "6. Watch Time Tracking"
echo "=========================================="

if [ -n "$race_id" ] && [ "$race_id" != "null" ]; then
    # Test 6.1: Start watch session
    response=$(curl -s -X POST "$BASE_URL/users/watch/sessions/start" \
        -H "Authorization: Bearer $USER_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"race_id\":\"$race_id\"}")
    watch_session_id=$(echo "$response" | jq -r '.session_id // empty' 2>/dev/null)
    if [ -n "$watch_session_id" ] && [ "$watch_session_id" != "null" ]; then
        pass_test "POST /users/watch/sessions/start creates watch session"
        
        # Test 6.2: End watch session
        sleep 1
        response=$(curl -s -X POST "$BASE_URL/users/watch/sessions/end" \
            -H "Authorization: Bearer $USER_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{\"session_id\":\"$watch_session_id\"}")
        if echo "$response" | jq -e '.success' > /dev/null 2>&1; then
            pass_test "POST /users/watch/sessions/end ends watch session"
        else
            fail_test "POST /users/watch/sessions/end" "Failed to end session"
        fi
        
        # Test 6.3: Get watch stats
        response=$(curl -s -X GET "$BASE_URL/users/watch/sessions/stats/$race_id" \
            -H "Authorization: Bearer $USER_TOKEN")
        if echo "$response" | jq -e '.total_minutes' > /dev/null 2>&1; then
            pass_test "GET /users/watch/sessions/stats/:race_id returns stats"
        else
            fail_test "GET /users/watch/sessions/stats/:race_id" "Invalid response"
        fi
    else
        fail_test "POST /users/watch/sessions/start" "Failed to create watch session"
    fi
else
    info "No race ID available, skipping watch time tracking tests"
fi

echo ""
echo "=========================================="
echo "7. Analytics Endpoints"
echo "=========================================="

# Test 7.1: Race analytics
response=$(curl -s -X GET "$BASE_URL/admin/analytics/races" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
if echo "$response" | jq -e 'type == "array"' > /dev/null 2>&1; then
    pass_test "GET /admin/analytics/races returns data"
else
    fail_test "GET /admin/analytics/races" "Invalid response"
fi

# Test 7.2: Watch time analytics
response=$(curl -s -X GET "$BASE_URL/admin/analytics/watch-time" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
if echo "$response" | jq -e 'type == "array"' > /dev/null 2>&1; then
    pass_test "GET /admin/analytics/watch-time returns data"
else
    fail_test "GET /admin/analytics/watch-time" "Invalid response"
fi

# Test 7.3: Revenue analytics
response=$(curl -s -X GET "$BASE_URL/admin/analytics/revenue" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
if echo "$response" | jq -e 'type == "array"' > /dev/null 2>&1; then
    pass_test "GET /admin/analytics/revenue returns data"
else
    fail_test "GET /admin/analytics/revenue" "Invalid response"
fi

echo ""
echo "=========================================="
echo "8. Cost Tracking"
echo "=========================================="

# Test 8.1: Create cost
cost_data="{\"cost_type\":\"cdn\",\"amount_cents\":5000,\"year\":2024,\"month\":12,\"description\":\"Test CDN cost\"}"
response=$(curl -s -X POST "$BASE_URL/admin/costs" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d "$cost_data")
cost_id=$(echo "$response" | jq -r '.id // empty' 2>/dev/null)
if [ -n "$cost_id" ] && [ "$cost_id" != "null" ]; then
    pass_test "POST /admin/costs creates cost"
    
    # Test 8.2: Get cost
    response=$(curl -s -X GET "$BASE_URL/admin/costs/$cost_id" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    if echo "$response" | jq -e '.id' > /dev/null 2>&1; then
        pass_test "GET /admin/costs/:id returns cost"
    else
        fail_test "GET /admin/costs/:id" "Invalid response"
    fi
    
    # Test 8.3: List costs
    response=$(curl -s -X GET "$BASE_URL/admin/costs" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    if echo "$response" | jq -e 'type == "array"' > /dev/null 2>&1; then
        pass_test "GET /admin/costs returns list"
    else
        fail_test "GET /admin/costs" "Invalid response"
    fi
    
    # Test 8.4: Delete cost
    response=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_URL/admin/costs/$cost_id" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    http_code=$(echo "$response" | tail -n1)
    if [ "$http_code" = "200" ] || [ "$http_code" = "204" ]; then
        pass_test "DELETE /admin/costs/:id deletes cost"
    else
        fail_test "DELETE /admin/costs/:id" "Expected 200/204, got $http_code"
    fi
else
    fail_test "POST /admin/costs" "Failed to create cost"
fi

echo ""
echo "=========================================="
echo "9. Revenue Endpoints"
echo "=========================================="

# Test 9.1: Get revenue
response=$(curl -s -X GET "$BASE_URL/admin/revenue" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
if echo "$response" | jq -e 'type == "array"' > /dev/null 2>&1; then
    pass_test "GET /admin/revenue returns data"
else
    fail_test "GET /admin/revenue" "Invalid response"
fi

echo ""
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo -e "Total Tests: ${TESTS_TOTAL}"
echo -e "${GREEN}Passed: ${TESTS_PASSED}${NC}"
echo -e "${RED}Failed: ${TESTS_FAILED}${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed. Please review the output above.${NC}"
    exit 1
fi

