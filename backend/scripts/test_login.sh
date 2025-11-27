#!/bin/bash

# Test script for login functionality
# Prerequisites: Backend server must be running on localhost:8080

set -e

BASE_URL="http://localhost:8080"

echo "🧪 Testing Login Functionality"
echo "=============================="
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
    fi
}

# Function to make API request
api_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    if [ -n "$data" ]; then
        curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data"
    else
        curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint"
    fi
}

# Test 1: Check if server is running
echo "Test 1: Checking if server is running..."
HEALTH_RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/health")
HEALTH_CODE=$(echo "$HEALTH_RESPONSE" | tail -1)
if [ "$HEALTH_CODE" = "200" ]; then
    print_result 0 "Server is running"
else
    print_result 1 "Server is not running (HTTP $HEALTH_CODE)"
    echo "Please start the backend server first: cd backend && go run cmd/api/main.go"
    exit 1
fi

# Test 2: Login with admin credentials
echo ""
echo "Test 2: Admin login..."
ADMIN_RESPONSE=$(api_request "POST" "/auth/login" '{"email":"admin@cyclingstream.local","password":"admin123"}')
ADMIN_CODE=$(echo "$ADMIN_RESPONSE" | tail -1)
ADMIN_BODY=$(echo "$ADMIN_RESPONSE" | sed '$d')

if [ "$ADMIN_CODE" = "200" ]; then
    ADMIN_TOKEN=$(echo "$ADMIN_BODY" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    if [ -n "$ADMIN_TOKEN" ]; then
        print_result 0 "Admin login successful"
        echo "  Token: ${ADMIN_TOKEN:0:20}..."
    else
        print_result 1 "Admin login failed - no token in response"
    fi
else
    print_result 1 "Admin login failed (HTTP $ADMIN_CODE)"
    echo "  Response: $ADMIN_BODY"
fi

# Small delay to avoid rate limiting
sleep 1

# Test 3: Login with invalid credentials
echo ""
echo "Test 3: Login with invalid credentials..."
INVALID_RESPONSE=$(api_request "POST" "/auth/login" '{"email":"admin@cyclingstream.local","password":"wrongpassword"}')
INVALID_CODE=$(echo "$INVALID_RESPONSE" | tail -1)
INVALID_BODY=$(echo "$INVALID_RESPONSE" | sed '$d')

if [ "$INVALID_CODE" = "401" ]; then
    print_result 0 "Invalid credentials correctly rejected (401)"
else
    print_result 1 "Invalid credentials test failed (expected 401, got $INVALID_CODE)"
    echo "  Response: $INVALID_BODY"
fi

sleep 1

# Test 4: Login with empty email
echo ""
echo "Test 4: Login with empty email..."
EMPTY_EMAIL_RESPONSE=$(api_request "POST" "/auth/login" '{"email":"","password":"test123"}')
EMPTY_EMAIL_CODE=$(echo "$EMPTY_EMAIL_RESPONSE" | tail -1)
EMPTY_EMAIL_BODY=$(echo "$EMPTY_EMAIL_RESPONSE" | sed '$d')

if [ "$EMPTY_EMAIL_CODE" = "400" ]; then
    print_result 0 "Empty email correctly rejected (400)"
else
    print_result 1 "Empty email test failed (expected 400, got $EMPTY_EMAIL_CODE)"
    echo "  Response: $EMPTY_EMAIL_BODY"
fi

sleep 1

# Test 5: Login with empty password
echo ""
echo "Test 5: Login with empty password..."
EMPTY_PASS_RESPONSE=$(api_request "POST" "/auth/login" '{"email":"test@example.com","password":""}')
EMPTY_PASS_CODE=$(echo "$EMPTY_PASS_RESPONSE" | tail -1)
EMPTY_PASS_BODY=$(echo "$EMPTY_PASS_RESPONSE" | sed '$d')

if [ "$EMPTY_PASS_CODE" = "400" ]; then
    print_result 0 "Empty password correctly rejected (400)"
else
    print_result 1 "Empty password test failed (expected 400, got $EMPTY_PASS_CODE)"
    echo "  Response: $EMPTY_PASS_BODY"
fi

sleep 1

# Test 6: Login with invalid email format
echo ""
echo "Test 6: Login with invalid email format..."
INVALID_EMAIL_RESPONSE=$(api_request "POST" "/auth/login" '{"email":"not-an-email","password":"test123"}')
INVALID_EMAIL_CODE=$(echo "$INVALID_EMAIL_RESPONSE" | tail -1)
INVALID_EMAIL_BODY=$(echo "$INVALID_EMAIL_RESPONSE" | sed '$d')

if [ "$INVALID_EMAIL_CODE" = "400" ]; then
    print_result 0 "Invalid email format correctly rejected (400)"
else
    print_result 1 "Invalid email format test failed (expected 400, got $INVALID_EMAIL_CODE)"
    echo "  Response: $INVALID_EMAIL_BODY"
fi

sleep 1

# Test 7: Login with non-existent user
echo ""
echo "Test 7: Login with non-existent user..."
NONEXISTENT_RESPONSE=$(api_request "POST" "/auth/login" '{"email":"nonexistent@example.com","password":"test123"}')
NONEXISTENT_CODE=$(echo "$NONEXISTENT_RESPONSE" | tail -1)
NONEXISTENT_BODY=$(echo "$NONEXISTENT_RESPONSE" | sed '$d')

if [ "$NONEXISTENT_CODE" = "401" ]; then
    print_result 0 "Non-existent user correctly rejected (401)"
else
    print_result 1 "Non-existent user test failed (expected 401, got $NONEXISTENT_CODE)"
    echo "  Response: $NONEXISTENT_BODY"
fi

sleep 1

# Test 8: Login with invalid JSON
echo ""
echo "Test 8: Login with invalid JSON..."
INVALID_JSON_RESPONSE=$(curl -s -w "\n%{http_code}" -X "POST" "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "invalid json")
INVALID_JSON_CODE=$(echo "$INVALID_JSON_RESPONSE" | tail -1)
INVALID_JSON_BODY=$(echo "$INVALID_JSON_RESPONSE" | sed '$d')

if [ "$INVALID_JSON_CODE" = "400" ]; then
    print_result 0 "Invalid JSON correctly rejected (400)"
else
    print_result 1 "Invalid JSON test failed (expected 400, got $INVALID_JSON_CODE)"
    echo "  Response: $INVALID_JSON_BODY"
fi

sleep 1

# Test 9: Register a user and then login
echo ""
echo "Test 9: Register and login with new user..."
REGISTER_RESPONSE=$(api_request "POST" "/auth/register" '{"email":"testuser'$(date +%s)'@example.com","password":"TestPassword123!","name":"Test User"}')
REGISTER_CODE=$(echo "$REGISTER_RESPONSE" | tail -1)
REGISTER_BODY=$(echo "$REGISTER_RESPONSE" | sed '$d')

if [ "$REGISTER_CODE" = "201" ]; then
    REGISTER_EMAIL=$(echo "$REGISTER_BODY" | grep -o '"email":"[^"]*' | cut -d'"' -f4)
    print_result 0 "User registration successful"
    
    # Now try to login
    echo "  Attempting login with registered user..."
    LOGIN_RESPONSE=$(api_request "POST" "/auth/login" "{\"email\":\"$REGISTER_EMAIL\",\"password\":\"TestPassword123!\"}")
    LOGIN_CODE=$(echo "$LOGIN_RESPONSE" | tail -1)
    LOGIN_BODY=$(echo "$LOGIN_RESPONSE" | sed '$d')
    
    if [ "$LOGIN_CODE" = "200" ]; then
        LOGIN_TOKEN=$(echo "$LOGIN_BODY" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
        if [ -n "$LOGIN_TOKEN" ]; then
            print_result 0 "Login with registered user successful"
            echo "  Token: ${LOGIN_TOKEN:0:20}..."
        else
            print_result 1 "Login failed - no token in response"
        fi
    else
        print_result 1 "Login with registered user failed (HTTP $LOGIN_CODE)"
        echo "  Response: $LOGIN_BODY"
    fi
else
    print_result 1 "User registration failed (HTTP $REGISTER_CODE)"
    echo "  Response: $REGISTER_BODY"
    echo "  Note: This might fail if user already exists or database is not accessible"
fi

# Summary
echo ""
echo "=============================="
echo "Login Testing Complete"
echo ""

