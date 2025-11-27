#!/bin/bash

# Shared utilities for API test scripts
# Source this file in test scripts: source "$(dirname "$0")/test_utils.sh"

# Default configuration
BASE_URL="${BASE_URL:-http://localhost:8080}"
ADMIN_EMAIL="${ADMIN_EMAIL:-admin@cyclingstream.local}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-admin123}"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0
TESTS_TOTAL=0
TESTS_SKIPPED=0

# Initialize test counters
reset_test_counters() {
    TESTS_PASSED=0
    TESTS_FAILED=0
    TESTS_TOTAL=0
    TESTS_SKIPPED=0
}

# Print test result
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

skip_test() {
    ((TESTS_SKIPPED++))
    ((TESTS_TOTAL++))
    echo -e "${YELLOW}⊘${NC} $1 (skipped)"
}

info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1" >&2
}

# Check if server is running
check_server() {
    if ! curl -s --max-time 5 "$BASE_URL/health" > /dev/null 2>&1; then
        error "Backend server is not running on $BASE_URL"
        echo "Please start the server with: make run-backend"
        exit 1
    fi
}

# Make API request and return response body and HTTP code
# Usage: response=$(api_request "GET" "/endpoint" '{"data":"value"}' "Bearer token")
# Returns: response body\nHTTP_CODE
api_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local auth_header=$4
    
    local headers=(-H "Content-Type: application/json")
    if [ -n "$auth_header" ]; then
        headers+=(-H "Authorization: $auth_header")
    fi
    
    if [ -n "$data" ]; then
        curl -s -w "\n%{http_code}" --max-time 10 -X "$method" "$BASE_URL$endpoint" \
            "${headers[@]}" \
            -d "$data" 2>/dev/null
    else
        curl -s -w "\n%{http_code}" --max-time 10 -X "$method" "$BASE_URL$endpoint" \
            "${headers[@]}" 2>/dev/null
    fi
}

# Extract HTTP code from response
get_http_code() {
    echo "$1" | tail -n1
}

# Extract response body from response
get_response_body() {
    echo "$1" | sed '$d'
}

# Get admin token
get_admin_token() {
    local response=$(api_request "POST" "/auth/login" \
        "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}")
    
    local body=$(get_response_body "$response")
    local token=$(echo "$body" | jq -r '.token // empty' 2>/dev/null)
    
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
    local register_response=$(api_request "POST" "/auth/register" \
        "{\"email\":\"$email\",\"password\":\"$password\"}")
    
    local register_code=$(get_http_code "$register_response")
    if [ "$register_code" != "201" ]; then
        echo ""
        return 1
    fi
    
    # Login
    local login_response=$(api_request "POST" "/auth/login" \
        "{\"email\":\"$email\",\"password\":\"$password\"}")
    
    local login_body=$(get_response_body "$login_response")
    local token=$(echo "$login_body" | jq -r '.token // empty' 2>/dev/null)
    
    if [ -z "$token" ] || [ "$token" = "null" ]; then
        echo ""
        return 1
    fi
    echo "$token"
}

# Validate JSON response
validate_json() {
    local json_string=$1
    if echo "$json_string" | jq . > /dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Check if jq is installed
check_jq() {
    if ! command -v jq > /dev/null 2>&1; then
        error "jq is required but not installed"
        echo "Install with: sudo apt install jq (Debian/Ubuntu) or brew install jq (macOS)"
        exit 1
    fi
}

# Print test summary
print_test_summary() {
    echo ""
    echo "=========================================="
    echo "Test Summary"
    echo "=========================================="
    echo -e "Total Tests: ${TESTS_TOTAL}"
    echo -e "${GREEN}Passed: ${TESTS_PASSED}${NC}"
    echo -e "${RED}Failed: ${TESTS_FAILED}${NC}"
    if [ $TESTS_SKIPPED -gt 0 ]; then
        echo -e "${YELLOW}Skipped: ${TESTS_SKIPPED}${NC}"
    fi
    echo ""
    
    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "${GREEN}All tests passed!${NC}"
        return 0
    else
        echo -e "${RED}Some tests failed. Please review the output above.${NC}"
        return 1
    fi
}

# Cleanup function (can be overridden in test scripts)
cleanup() {
    # Default cleanup - can be extended in individual test scripts
    true
}

# Trap to ensure cleanup runs on exit
trap cleanup EXIT

# Print script header
print_test_header() {
    local script_name=$1
    echo "=========================================="
    echo "$script_name"
    echo "=========================================="
    echo ""
    info "Base URL: $BASE_URL"
    echo ""
}

