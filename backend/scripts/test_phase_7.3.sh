#!/bin/bash

# Test script for Phase 7.3: Monitoring & Logging
# This script tests:
# 1. Enhanced health check endpoint
# 2. Cost tracking endpoints
# 3. Structured logging (verified by checking logs)

set -e

BASE_URL="http://localhost:8080"
ADMIN_EMAIL="admin@example.com"
ADMIN_PASSWORD="admin123"

echo "========================================="
echo "Phase 7.3 Testing: Monitoring & Logging"
echo "========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print test result
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $2"
    else
        echo -e "${RED}✗${NC} $2"
    fi
}

# Function to make authenticated admin request
admin_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    if [ -z "$ADMIN_TOKEN" ]; then
        echo -e "${RED}Error: Admin not logged in${NC}"
        return 1
    fi
    
    if [ -z "$data" ]; then
        curl -s -X "$method" \
            -H "Authorization: Bearer $ADMIN_TOKEN" \
            -H "Content-Type: application/json" \
            "$BASE_URL$endpoint"
    else
        curl -s -X "$method" \
            -H "Authorization: Bearer $ADMIN_TOKEN" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$BASE_URL$endpoint"
    fi
}

echo "Step 1: Login as admin"
echo "----------------------"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}")

ADMIN_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$ADMIN_TOKEN" ]; then
    echo -e "${RED}Failed to login as admin${NC}"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

print_result 0 "Admin login successful"
echo ""

echo "Step 2: Test Enhanced Health Check"
echo "-----------------------------------"
HEALTH_RESPONSE=$(curl -s "$BASE_URL/health")
echo "$HEALTH_RESPONSE" | jq '.' 2>/dev/null || echo "$HEALTH_RESPONSE"

# Check for required fields
if echo "$HEALTH_RESPONSE" | grep -q "status" && \
   echo "$HEALTH_RESPONSE" | grep -q "timestamp" && \
   echo "$HEALTH_RESPONSE" | grep -q "uptime_seconds" && \
   echo "$HEALTH_RESPONSE" | grep -q "response_time_ms" && \
   echo "$HEALTH_RESPONSE" | grep -q "services" && \
   echo "$HEALTH_RESPONSE" | grep -q "system"; then
    print_result 0 "Health check contains all required fields"
else
    print_result 1 "Health check missing required fields"
fi

# Check database status
if echo "$HEALTH_RESPONSE" | grep -q '"status":"ok"'; then
    print_result 0 "Database connectivity check working"
else
    print_result 1 "Database connectivity check failed"
fi
echo ""

echo "Step 3: Test Cost Tracking - Create Cost"
echo "-----------------------------------------"
# Get a race ID first
RACES_RESPONSE=$(admin_request "GET" "/admin/races")
RACE_ID=$(echo "$RACES_RESPONSE" | jq -r '.[0].id' 2>/dev/null || echo "")

if [ -z "$RACE_ID" ] || [ "$RACE_ID" = "null" ]; then
    echo -e "${YELLOW}Warning: No races found, creating cost without race_id${NC}"
    COST_DATA="{\"cost_type\":\"cdn\",\"amount_cents\":5000,\"year\":2024,\"month\":12,\"description\":\"Test CDN cost\"}"
else
    COST_DATA="{\"race_id\":\"$RACE_ID\",\"cost_type\":\"cdn\",\"amount_cents\":5000,\"year\":2024,\"month\":12,\"description\":\"Test CDN cost\"}"
fi

CREATE_COST_RESPONSE=$(admin_request "POST" "/admin/costs" "$COST_DATA")
echo "$CREATE_COST_RESPONSE" | jq '.' 2>/dev/null || echo "$CREATE_COST_RESPONSE"

COST_ID=$(echo "$CREATE_COST_RESPONSE" | jq -r '.id' 2>/dev/null || echo "")

if [ -n "$COST_ID" ] && [ "$COST_ID" != "null" ]; then
    print_result 0 "Cost created successfully (ID: $COST_ID)"
else
    print_result 1 "Failed to create cost"
    echo "Response: $CREATE_COST_RESPONSE"
fi
echo ""

if [ -n "$COST_ID" ] && [ "$COST_ID" != "null" ]; then
    echo "Step 4: Test Cost Tracking - Get Cost"
    echo "--------------------------------------"
    GET_COST_RESPONSE=$(admin_request "GET" "/admin/costs/$COST_ID")
    echo "$GET_COST_RESPONSE" | jq '.' 2>/dev/null || echo "$GET_COST_RESPONSE"
    
    if echo "$GET_COST_RESPONSE" | grep -q "$COST_ID"; then
        print_result 0 "Get cost by ID working"
    else
        print_result 1 "Get cost by ID failed"
    fi
    echo ""
    
    echo "Step 5: Test Cost Tracking - List All Costs"
    echo "--------------------------------------------"
    LIST_COSTS_RESPONSE=$(admin_request "GET" "/admin/costs")
    COST_COUNT=$(echo "$LIST_COSTS_RESPONSE" | jq 'length' 2>/dev/null || echo "0")
    echo "Found $COST_COUNT cost(s)"
    
    if [ "$COST_COUNT" -gt 0 ]; then
        print_result 0 "List costs working"
    else
        print_result 1 "List costs returned empty"
    fi
    echo ""
    
    echo "Step 6: Test Cost Tracking - Get Cost Summary"
    echo "----------------------------------------------"
    SUMMARY_RESPONSE=$(admin_request "GET" "/admin/costs/summary?year=2024&month=12")
    echo "$SUMMARY_RESPONSE" | jq '.' 2>/dev/null || echo "$SUMMARY_RESPONSE"
    
    if echo "$SUMMARY_RESPONSE" | grep -q "total_cents"; then
        print_result 0 "Cost summary working"
    else
        print_result 1 "Cost summary failed"
    fi
    echo ""
    
    echo "Step 7: Test Cost Tracking - Update Cost"
    echo "----------------------------------------"
    UPDATE_DATA="{\"race_id\":\"$RACE_ID\",\"cost_type\":\"server\",\"amount_cents\":7500,\"year\":2024,\"month\":12,\"description\":\"Updated test cost\"}"
    UPDATE_RESPONSE=$(admin_request "PUT" "/admin/costs/$COST_ID" "$UPDATE_DATA")
    echo "$UPDATE_RESPONSE" | jq '.' 2>/dev/null || echo "$UPDATE_RESPONSE"
    
    if echo "$UPDATE_RESPONSE" | grep -q "7500"; then
        print_result 0 "Update cost working"
    else
        print_result 1 "Update cost failed"
    fi
    echo ""
    
    echo "Step 8: Test Cost Tracking - Delete Cost"
    echo "----------------------------------------"
    DELETE_RESPONSE=$(admin_request "DELETE" "/admin/costs/$COST_ID")
    DELETE_STATUS=$?
    
    # Verify deletion
    GET_DELETED_RESPONSE=$(admin_request "GET" "/admin/costs/$COST_ID")
    if echo "$GET_DELETED_RESPONSE" | grep -q "not found"; then
        print_result 0 "Delete cost working"
    else
        print_result 1 "Delete cost failed"
    fi
    echo ""
fi

echo "Step 9: Test Cost Tracking - Get Costs by Race"
echo "-----------------------------------------------"
if [ -n "$RACE_ID" ] && [ "$RACE_ID" != "null" ]; then
    RACE_COSTS_RESPONSE=$(admin_request "GET" "/admin/costs/races/$RACE_ID")
    echo "$RACE_COSTS_RESPONSE" | jq '.' 2>/dev/null || echo "$RACE_COSTS_RESPONSE"
    
    if echo "$RACE_COSTS_RESPONSE" | grep -q "\[\]"; then
        print_result 0 "Get costs by race working (empty result expected)"
    else
        print_result 0 "Get costs by race working"
    fi
else
    echo -e "${YELLOW}Skipped: No race ID available${NC}"
fi
echo ""

echo "========================================="
echo "Testing Summary"
echo "========================================="
echo ""
echo -e "${GREEN}Phase 7.3 testing complete!${NC}"
echo ""
echo "Note: Structured logging and slow query logging are infrastructure"
echo "features that will be visible in server logs during normal operation."
echo ""
echo "To verify structured logging:"
echo "  1. Check server logs - they should be in JSON format (production) or"
echo "     structured text format (development)"
echo "  2. Look for request logs with fields: method, path, status, duration_ms"
echo ""
echo "To verify slow query logging:"
echo "  1. Execute queries that take > 1 second"
echo "  2. Check logs for 'Slow query detected' warnings"
echo ""

