#!/bin/bash

# Test script for analytics ingestion endpoint
# Usage: ./test_analytics_ingestion.sh

API_URL="http://localhost:8080"
STREAM_ID="6f447475-d86f-42de-9be7-bf304dbb4d78"
CLIENT_ID="test-client-$(date +%s)"

echo "=== Testing Analytics Ingestion Endpoint ==="
echo ""

# Test 1: Valid event batch
echo "Test 1: Valid event batch (play, heartbeat, ended)"
curl -X POST "$API_URL/analytics/events" \
  -H "Content-Type: application/json" \
  -d "{
    \"streamId\": \"$STREAM_ID\",
    \"clientId\": \"$CLIENT_ID\",
    \"events\": [
      {\"type\": \"play\", \"videoTime\": 0},
      {\"type\": \"heartbeat\", \"videoTime\": 15},
      {\"type\": \"heartbeat\", \"videoTime\": 30},
      {\"type\": \"ended\", \"videoTime\": 45}
    ]
  }" | jq .
echo ""
echo ""

# Test 2: Invalid stream ID
echo "Test 2: Invalid stream ID (should return 404)"
curl -X POST "$API_URL/analytics/events" \
  -H "Content-Type: application/json" \
  -d "{
    \"streamId\": \"00000000-0000-0000-0000-000000000000\",
    \"clientId\": \"$CLIENT_ID\",
    \"events\": [{\"type\": \"play\"}]
  }" | jq .
echo ""
echo ""

# Test 3: Missing streamId
echo "Test 3: Missing streamId (should return 400)"
curl -X POST "$API_URL/analytics/events" \
  -H "Content-Type: application/json" \
  -d "{
    \"clientId\": \"$CLIENT_ID\",
    \"events\": [{\"type\": \"play\"}]
  }" | jq .
echo ""
echo ""

# Test 4: Missing clientId
echo "Test 4: Missing clientId (should return 400)"
curl -X POST "$API_URL/analytics/events" \
  -H "Content-Type: application/json" \
  -d "{
    \"streamId\": \"$STREAM_ID\",
    \"events\": [{\"type\": \"play\"}]
  }" | jq .
echo ""
echo ""

# Test 5: Empty events array
echo "Test 5: Empty events array (should return 400)"
curl -X POST "$API_URL/analytics/events" \
  -H "Content-Type: application/json" \
  -d "{
    \"streamId\": \"$STREAM_ID\",
    \"clientId\": \"$CLIENT_ID\",
    \"events\": []
  }" | jq .
echo ""
echo ""

# Test 6: Batch size limit (101 events, should return 400)
echo "Test 6: Batch size limit exceeded (101 events, should return 400)"
EVENTS=$(for i in {1..101}; do echo "{\"type\": \"heartbeat\", \"videoTime\": $i}"; done | jq -s .)
curl -X POST "$API_URL/analytics/events" \
  -H "Content-Type: application/json" \
  -d "{
    \"streamId\": \"$STREAM_ID\",
    \"clientId\": \"$CLIENT_ID\",
    \"events\": $EVENTS
  }" | jq .
echo ""
echo ""

# Test 7: Invalid event type
echo "Test 7: Invalid event type (should return 400)"
curl -X POST "$API_URL/analytics/events" \
  -H "Content-Type: application/json" \
  -d "{
    \"streamId\": \"$STREAM_ID\",
    \"clientId\": \"$CLIENT_ID\",
    \"events\": [{\"type\": \"invalid_event_type\"}]
  }" | jq .
echo ""
echo ""

# Test 8: Negative videoTime
echo "Test 8: Negative videoTime (should return 400)"
curl -X POST "$API_URL/analytics/events" \
  -H "Content-Type: application/json" \
  -d "{
    \"streamId\": \"$STREAM_ID\",
    \"clientId\": \"$CLIENT_ID\",
    \"events\": [{\"type\": \"play\", \"videoTime\": -1}]
  }" | jq .
echo ""
echo ""

# Test 9: Valid batch with all event types
echo "Test 9: Valid batch with all event types (play, pause, heartbeat, ended, error, buffer_start, buffer_end)"
curl -X POST "$API_URL/analytics/events" \
  -H "Content-Type: application/json" \
  -d "{
    \"streamId\": \"$STREAM_ID\",
    \"clientId\": \"$CLIENT_ID-all-events\",
    \"events\": [
      {\"type\": \"play\", \"videoTime\": 0},
      {\"type\": \"pause\", \"videoTime\": 10},
      {\"type\": \"play\", \"videoTime\": 10},
      {\"type\": \"buffer_start\", \"videoTime\": 20},
      {\"type\": \"buffer_end\", \"videoTime\": 25},
      {\"type\": \"error\", \"videoTime\": 30, \"extra\": {\"error_code\": \"NETWORK_ERROR\"}},
      {\"type\": \"heartbeat\", \"videoTime\": 45},
      {\"type\": \"ended\", \"videoTime\": 60}
    ]
  }" | jq .
echo ""
echo ""

# Test 10: Device type detection (User-Agent)
echo "Test 10: Device type detection from User-Agent"
curl -X POST "$API_URL/analytics/events" \
  -H "Content-Type: application/json" \
  -H "User-Agent: Mozilla/5.0 (iPad; CPU OS 14_0 like Mac OS X)" \
  -d "{
    \"streamId\": \"$STREAM_ID\",
    \"clientId\": \"$CLIENT_ID-tablet\",
    \"events\": [{\"type\": \"play\"}]
  }" | jq .
echo ""
echo ""

# Test 11: Country detection from headers
echo "Test 11: Country detection from CF-IPCountry header"
curl -X POST "$API_URL/analytics/events" \
  -H "Content-Type: application/json" \
  -H "CF-IPCountry: US" \
  -d "{
    \"streamId\": \"$STREAM_ID\",
    \"clientId\": \"$CLIENT_ID-country\",
    \"events\": [{\"type\": \"play\"}]
  }" | jq .
echo ""
echo ""

# Test 12: Verify events were persisted
echo "Test 12: Verify events were persisted to database"
echo "Checking playback_events table..."
docker exec -i cyclingstream_postgres psql -U cyclingstream -d cyclingstream -c "SELECT COUNT(*) as event_count, COUNT(DISTINCT client_id) as unique_clients FROM playback_events WHERE stream_id = '$STREAM_ID';"
echo ""

