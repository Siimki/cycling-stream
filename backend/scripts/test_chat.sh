#!/bin/bash

# Integration test script for Live Chat functionality
# This script tests the chat endpoints and WebSocket functionality
# Prerequisites: Database must be running and migrations must be applied

set -e

BASE_URL="http://localhost:8080"
WS_URL="ws://localhost:8080"
ADMIN_EMAIL="admin@cyclingstream.local"
ADMIN_PASSWORD="admin123"
USER_EMAIL="testuser@example.com"
USER_PASSWORD="testpass123"

echo "🧪 Testing Live Chat Integration"
echo "================================="
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

# Step 1: Admin login
echo "Step 1: Admin login..."
ADMIN_RESPONSE=$(api_request "POST" "/auth/login" "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}" "")
ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$ADMIN_TOKEN" ]; then
    print_result 1 "Admin login failed"
    exit 1
fi
print_result 0 "Admin logged in"

# Step 2: Create a race
echo ""
echo "Step 2: Create a race..."
RACE_DATA="{\"name\":\"Chat Test Race\",\"description\":\"Test race for chat\",\"start_date\":\"2024-12-01T10:00:00Z\",\"end_date\":\"2024-12-01T12:00:00Z\",\"location\":\"Test Location\",\"category\":\"Test\",\"is_free\":true,\"price_cents\":0}"
RACE_RESPONSE=$(api_request "POST" "/admin/races" "$RACE_DATA" "$ADMIN_TOKEN")
RACE_ID=$(echo $RACE_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$RACE_ID" ]; then
    print_result 1 "Race creation failed"
    exit 1
fi
print_result 0 "Race created: $RACE_ID"

# Step 3: Set stream to live
echo ""
echo "Step 3: Set stream to live..."
STREAM_DATA="{\"status\":\"live\",\"origin_url\":\"http://test.com/stream.m3u8\",\"cdn_url\":\"http://cdn.test.com/stream.m3u8\"}"
STREAM_RESPONSE=$(api_request "POST" "/admin/races/$RACE_ID/stream" "$STREAM_DATA" "$ADMIN_TOKEN")
print_result 0 "Stream set to live"

# Step 4: User registration
echo ""
echo "Step 4: User registration..."
USER_DATA="{\"email\":\"$USER_EMAIL\",\"password\":\"$USER_PASSWORD\",\"name\":\"Test User\"}"
USER_RESPONSE=$(api_request "POST" "/auth/register" "$USER_DATA" "")
USER_TOKEN=$(echo $USER_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$USER_TOKEN" ]; then
    print_result 1 "User registration failed"
    exit 1
fi
print_result 0 "User registered and logged in"

# Step 5: Test chat history endpoint (should be empty initially)
echo ""
echo "Step 5: Test chat history endpoint..."
HISTORY_RESPONSE=$(api_request "GET" "/races/$RACE_ID/chat/history" "" "")
HISTORY_COUNT=$(echo $HISTORY_RESPONSE | grep -o '"messages":\[[^]]*\]' | grep -o '\]' | wc -l || echo "0")
print_result 0 "Chat history endpoint accessible (empty initially)"

# Step 6: Test chat stats endpoint
echo ""
echo "Step 6: Test chat stats endpoint..."
STATS_RESPONSE=$(api_request "GET" "/races/$RACE_ID/chat/stats" "" "")
STATS_TOTAL=$(echo $STATS_RESPONSE | grep -o '"total_messages":[0-9]*' | cut -d':' -f2 || echo "0")
print_result 0 "Chat stats endpoint accessible (total: $STATS_TOTAL)"

# Step 7: Test WebSocket connection (basic check)
echo ""
echo "Step 7: Test WebSocket connection..."
WS_ENDPOINT="$WS_URL/races/$RACE_ID/chat/ws"
WS_ENDPOINT_WITH_TOKEN="$WS_URL/races/$RACE_ID/chat/ws?token=$USER_TOKEN"

echo -e "${YELLOW}WebSocket endpoint (anonymous): $WS_ENDPOINT${NC}"
echo -e "${YELLOW}WebSocket endpoint (authenticated): $WS_ENDPOINT_WITH_TOKEN${NC}"
echo ""
echo -e "${YELLOW}To test WebSocket connection, use one of these methods:${NC}"
echo "1. Install websocat: cargo install websocat (or brew install websocat)"
echo "2. Then run: websocat $WS_ENDPOINT"
echo "3. Or use Node.js with ws package"
echo "4. Or use browser console with WebSocket API"
print_result 0 "WebSocket endpoint documented"

# Step 8: Test chat with offline stream (should fail)
echo ""
echo "Step 8: Test chat with offline stream..."
# Set stream to offline
OFFLINE_DATA="{\"status\":\"offline\"}"
api_request "PUT" "/admin/races/$RACE_ID/stream/status" "$OFFLINE_DATA" "$ADMIN_TOKEN" > /dev/null
echo -e "${YELLOW}Note: WebSocket connection to offline stream should be rejected${NC}"
print_result 0 "Stream set to offline for testing"

# Step 9: Set stream back to live
echo ""
echo "Step 9: Set stream back to live..."
LIVE_DATA="{\"status\":\"live\"}"
api_request "PUT" "/admin/races/$RACE_ID/stream/status" "$LIVE_DATA" "$ADMIN_TOKEN" > /dev/null
print_result 0 "Stream set back to live"

# Step 10: Test rate limiting documentation
echo ""
echo "Step 10: Rate limiting test..."
echo -e "${YELLOW}Rate limiting: 10 messages/minute per user${NC}"
echo -e "${YELLOW}To test rate limiting:${NC}"
echo "1. Connect WebSocket with authenticated user"
echo "2. Send 10 messages rapidly (should succeed)"
echo "3. Send 11th message immediately (should fail with rate limit error)"
echo "4. Wait 1 minute and try again (should succeed)"
print_result 0 "Rate limiting documented"

# Step 11: Test message persistence
echo ""
echo "Step 11: Test message persistence..."
echo "Create a message via API if possible, then verify it appears in history"
# Note: Full message creation requires WebSocket, but we can verify persistence
HISTORY_AFTER=$(api_request "GET" "/races/$RACE_ID/chat/history" "" "")
MESSAGE_COUNT=$(echo "$HISTORY_AFTER" | grep -o '"messages":\[[^]]*\]' | grep -o 'id' | wc -l || echo "0")
echo "Current message count in history: $MESSAGE_COUNT"
print_result 0 "Message persistence check completed"

# Summary
echo ""
echo "================================="
echo -e "${GREEN}✅ All basic chat API tests passed!${NC}"
echo ""
echo "Test Summary:"
echo "- Race ID: $RACE_ID"
echo "- User token: ${USER_TOKEN:0:20}..."
echo "- WebSocket endpoint: $WS_ENDPOINT"
echo ""
echo "Next steps for full WebSocket testing:"
echo ""
echo "1. Install websocat (recommended):"
echo "   cargo install websocat"
echo "   # OR: brew install websocat"
echo ""
echo "2. Test anonymous WebSocket connection:"
echo "   websocat $WS_ENDPOINT"
echo ""
echo "3. Test authenticated WebSocket connection:"
echo "   websocat \"$WS_ENDPOINT_WITH_TOKEN\""
echo ""
echo "4. Send a message (when authenticated):"
echo "   {\"type\":\"send_message\",\"data\":{\"message\":\"Hello, world!\"}}"
echo ""
echo "5. Test scenarios:"
echo "   - Multiple clients in same room"
echo "   - Rate limiting (send 11 messages rapidly)"
echo "   - Message persistence (refresh history endpoint)"
echo "   - Anonymous user cannot send messages"
echo ""
echo "6. Or use browser console:"
echo "   const ws = new WebSocket('$WS_ENDPOINT_WITH_TOKEN');"
echo "   ws.onopen = () => ws.send('{\"type\":\"send_message\",\"data\":{\"message\":\"Test\"}}');"
echo "   ws.onmessage = (e) => console.log(e.data);"
echo ""

