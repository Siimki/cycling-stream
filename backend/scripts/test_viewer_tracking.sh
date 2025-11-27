#!/bin/bash

# Integration test script for Viewer Tracking functionality (Phase 7.1)
# This script tests the viewer tracking endpoints
# Prerequisites: Database must be running and migrations must be applied

set -e

BASE_URL="http://localhost:8080"
ADMIN_EMAIL="admin@cyclingstream.local"
ADMIN_PASSWORD="admin123"
USER_EMAIL="testuser@example.com"
USER_PASSWORD="testpass123"

echo "🧪 Testing Viewer Tracking Integration (Phase 7.1)"
echo "=================================================="
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
    local save_cookies=$5
    
    if [ -n "$token" ]; then
        if [ -n "$data" ]; then
            if [ "$save_cookies" = "true" ]; then
                curl -s -X "$method" "$BASE_URL$endpoint" \
                    -H "Content-Type: application/json" \
                    -H "Authorization: Bearer $token" \
                    -d "$data" \
                    -c /tmp/viewer_test_cookies.txt \
                    -b /tmp/viewer_test_cookies.txt
            else
                curl -s -X "$method" "$BASE_URL$endpoint" \
                    -H "Content-Type: application/json" \
                    -H "Authorization: Bearer $token" \
                    -d "$data"
            fi
        else
            if [ "$save_cookies" = "true" ]; then
                curl -s -X "$method" "$BASE_URL$endpoint" \
                    -H "Authorization: Bearer $token" \
                    -c /tmp/viewer_test_cookies.txt \
                    -b /tmp/viewer_test_cookies.txt
            else
                curl -s -X "$method" "$BASE_URL$endpoint" \
                    -H "Authorization: Bearer $token"
            fi
        fi
    else
        if [ -n "$data" ]; then
            if [ "$save_cookies" = "true" ]; then
                curl -s -X "$method" "$BASE_URL$endpoint" \
                    -H "Content-Type: application/json" \
                    -d "$data" \
                    -c /tmp/viewer_test_cookies.txt \
                    -b /tmp/viewer_test_cookies.txt
            else
                curl -s -X "$method" "$BASE_URL$endpoint" \
                    -H "Content-Type: application/json" \
                    -d "$data"
            fi
        else
            if [ "$save_cookies" = "true" ]; then
                curl -s -X "$method" "$BASE_URL$endpoint" \
                    -c /tmp/viewer_test_cookies.txt \
                    -b /tmp/viewer_test_cookies.txt
            else
                curl -s -X "$method" "$BASE_URL$endpoint"
            fi
        fi
    fi
}

# Clean up cookies file
rm -f /tmp/viewer_test_cookies.txt

# Test 1: Health check
echo "1. Testing health endpoint..."
response=$(api_request "GET" "/health")
if echo "$response" | grep -q "status"; then
    print_result 0 "Health check passed"
else
    print_result 1 "Health check failed"
fi

# Test 2: Admin login (to create a test race)
echo ""
echo "2. Testing admin login..."
login_response=$(api_request "POST" "/auth/login" "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}")
ADMIN_TOKEN=$(echo "$login_response" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -n "$ADMIN_TOKEN" ]; then
    print_result 0 "Admin login successful"
    echo "   Token: ${ADMIN_TOKEN:0:20}..."
else
    print_result 1 "Admin login failed"
    echo "   Response: $login_response"
    exit 1
fi

# Test 3: Create a test race
echo ""
echo "3. Creating test race for viewer tracking..."
race_data='{
    "name": "Test Viewer Tracking Race",
    "description": "Race for viewer tracking testing",
    "is_free": true,
    "price_cents": 0
}'
race_response=$(api_request "POST" "/admin/races" "$race_data" "$ADMIN_TOKEN")
RACE_ID=$(echo "$race_response" | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -n "$RACE_ID" ]; then
    print_result 0 "Test race created (ID: ${RACE_ID:0:8}...)"
else
    print_result 1 "Failed to create test race"
    echo "   Response: $race_response"
    exit 1
fi

# Test 4: Start anonymous viewer session
echo ""
echo "4. Testing anonymous viewer session start..."
anon_session_data="{\"race_id\":\"$RACE_ID\"}"
anon_session_response=$(api_request "POST" "/viewers/sessions/start" "$anon_session_data" "" "true")
ANON_SESSION_ID=$(echo "$anon_session_response" | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -n "$ANON_SESSION_ID" ]; then
    print_result 0 "Anonymous viewer session started (ID: ${ANON_SESSION_ID:0:8}...)"
else
    print_result 1 "Failed to start anonymous viewer session"
    echo "   Response: $anon_session_response"
    exit 1
fi

# Test 5: Get concurrent viewers (should be 1)
echo ""
echo "5. Testing GET /races/:id/viewers/concurrent..."
concurrent_response=$(api_request "GET" "/races/$RACE_ID/viewers/concurrent")
CONCURRENT_COUNT=$(echo "$concurrent_response" | grep -o '"concurrent_count":[0-9]*' | cut -d':' -f2)

if [ "$CONCURRENT_COUNT" = "1" ]; then
    print_result 0 "Concurrent viewers count correct (1)"
else
    print_result 1 "Concurrent viewers count incorrect (expected 1, got $CONCURRENT_COUNT)"
    echo "   Response: $concurrent_response"
fi

# Test 6: Get unique viewers (should be 1)
echo ""
echo "6. Testing GET /races/:id/viewers/unique..."
unique_response=$(api_request "GET" "/races/$RACE_ID/viewers/unique")
UNIQUE_COUNT=$(echo "$unique_response" | grep -o '"unique_viewer_count":[0-9]*' | cut -d':' -f2)

if [ "$UNIQUE_COUNT" = "1" ]; then
    print_result 0 "Unique viewers count correct (1)"
else
    print_result 1 "Unique viewers count incorrect (expected 1, got $UNIQUE_COUNT)"
    echo "   Response: $unique_response"
fi

# Test 7: Send heartbeat
echo ""
echo "7. Testing POST /viewers/sessions/heartbeat..."
heartbeat_data="{\"session_id\":\"$ANON_SESSION_ID\"}"
heartbeat_response=$(api_request "POST" "/viewers/sessions/heartbeat" "$heartbeat_data" "" "true")

if echo "$heartbeat_response" | grep -q "successfully"; then
    print_result 0 "Heartbeat updated successfully"
else
    print_result 1 "Heartbeat update failed"
    echo "   Response: $heartbeat_response"
fi

# Test 8: Register and login as regular user
echo ""
echo "8. Registering test user..."
register_data="{\"email\":\"$USER_EMAIL\",\"password\":\"$USER_PASSWORD\",\"name\":\"Test User\"}"
register_response=$(api_request "POST" "/auth/register" "$register_data")

if echo "$register_response" | grep -q "token"; then
    print_result 0 "User registered successfully"
    USER_TOKEN=$(echo "$register_response" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
else
    # User might already exist, try login
    echo "   User might already exist, trying login..."
    login_response=$(api_request "POST" "/auth/login" "{\"email\":\"$USER_EMAIL\",\"password\":\"$USER_PASSWORD\"}")
    USER_TOKEN=$(echo "$login_response" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    if [ -n "$USER_TOKEN" ]; then
        print_result 0 "User login successful"
    else
        print_result 1 "User registration/login failed"
        echo "   Response: $register_response"
        exit 1
    fi
fi

# Test 9: Start authenticated viewer session
echo ""
echo "9. Testing authenticated viewer session start..."
auth_session_data="{\"race_id\":\"$RACE_ID\"}"
auth_session_response=$(api_request "POST" "/viewers/sessions/start" "$auth_session_data" "$USER_TOKEN")
AUTH_SESSION_ID=$(echo "$auth_session_response" | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -n "$AUTH_SESSION_ID" ]; then
    print_result 0 "Authenticated viewer session started (ID: ${AUTH_SESSION_ID:0:8}...)"
else
    print_result 1 "Failed to start authenticated viewer session"
    echo "   Response: $auth_session_response"
    exit 1
fi

# Test 10: Get concurrent viewers (should be 2 now)
echo ""
echo "10. Testing concurrent viewers with 2 sessions..."
concurrent_response2=$(api_request "GET" "/races/$RACE_ID/viewers/concurrent")
CONCURRENT_COUNT2=$(echo "$concurrent_response2" | grep -o '"concurrent_count":[0-9]*' | cut -d':' -f2)
AUTH_COUNT=$(echo "$concurrent_response2" | grep -o '"authenticated_count":[0-9]*' | cut -d':' -f2)
ANON_COUNT=$(echo "$concurrent_response2" | grep -o '"anonymous_count":[0-9]*' | cut -d':' -f2)

if [ "$CONCURRENT_COUNT2" = "2" ] && [ "$AUTH_COUNT" = "1" ] && [ "$ANON_COUNT" = "1" ]; then
    print_result 0 "Concurrent viewers count correct (2 total: 1 auth, 1 anon)"
else
    print_result 1 "Concurrent viewers count incorrect (expected 2, got $CONCURRENT_COUNT2)"
    echo "   Response: $concurrent_response2"
fi

# Test 11: Try to start session again with same token (should return existing active session)
echo ""
echo "11. Testing session reuse (same token, before ending)..."
anon_session_response2=$(api_request "POST" "/viewers/sessions/start" "$anon_session_data" "" "true")
ANON_SESSION_ID2=$(echo "$anon_session_response2" | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ "$ANON_SESSION_ID" = "$ANON_SESSION_ID2" ]; then
    print_result 0 "Session reuse works correctly (same session ID returned)"
else
    print_result 1 "Session reuse failed (different session ID)"
    echo "   Original: $ANON_SESSION_ID"
    echo "   New: $ANON_SESSION_ID2"
fi

# Test 12: End anonymous session
echo ""
echo "12. Testing end viewer session..."
end_session_data="{\"session_id\":\"$ANON_SESSION_ID\"}"
end_session_response=$(api_request "POST" "/viewers/sessions/end" "$end_session_data" "" "true")

if echo "$end_session_response" | grep -q "successfully"; then
    print_result 0 "Viewer session ended successfully"
else
    print_result 1 "Failed to end viewer session"
    echo "   Response: $end_session_response"
fi

# Test 13: Get concurrent viewers (should be 1 now)
echo ""
echo "13. Testing concurrent viewers after ending session..."
concurrent_response3=$(api_request "GET" "/races/$RACE_ID/viewers/concurrent")
CONCURRENT_COUNT3=$(echo "$concurrent_response3" | grep -o '"concurrent_count":[0-9]*' | cut -d':' -f2)

if [ "$CONCURRENT_COUNT3" = "1" ]; then
    print_result 0 "Concurrent viewers count correct after ending session (1)"
else
    print_result 1 "Concurrent viewers count incorrect (expected 1, got $CONCURRENT_COUNT3)"
    echo "   Response: $concurrent_response3"
fi

# Test 14: Get unique viewers (should still be 2)
echo ""
echo "14. Testing unique viewers (should still be 2)..."
unique_response2=$(api_request "GET" "/races/$RACE_ID/viewers/unique")
UNIQUE_COUNT2=$(echo "$unique_response2" | grep -o '"unique_viewer_count":[0-9]*' | cut -d':' -f2)

if [ "$UNIQUE_COUNT2" = "2" ]; then
    print_result 0 "Unique viewers count correct (2)"
else
    print_result 1 "Unique viewers count incorrect (expected 2, got $UNIQUE_COUNT2)"
    echo "   Response: $unique_response2"
fi

# Cleanup: End all sessions
echo ""
echo "Cleaning up test sessions..."
# ANON_SESSION_ID was already ended, but try ending AUTH_SESSION_ID
api_request "POST" "/viewers/sessions/end" "{\"session_id\":\"$AUTH_SESSION_ID\"}" "$USER_TOKEN" > /dev/null

echo ""
echo "=================================================="
echo -e "${GREEN}✅ All viewer tracking tests passed!${NC}"
echo ""
echo "Phase 7.1: Viewer Tracking - Implementation Complete"
echo ""
echo "Tested features:"
echo "  ✅ Anonymous viewer session tracking"
echo "  ✅ Authenticated viewer session tracking"
echo "  ✅ Concurrent viewer counting"
echo "  ✅ Unique viewer counting"
echo "  ✅ Session heartbeat mechanism"
echo "  ✅ Session ending"
echo "  ✅ Session reuse (same token)"
echo ""

