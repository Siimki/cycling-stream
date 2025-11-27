# Test Scripts Documentation

This directory contains test scripts for the CyclingStream backend API.

## Prerequisites

- Backend server running on `localhost:8080` (or set `BASE_URL` environment variable)
- `jq` installed for JSON parsing
- `curl` installed for HTTP requests
- Database running and migrations applied

## Available Scripts

### `test_all_api.sh`
Comprehensive API test suite that tests all critical backend endpoints.

**Usage:**
```bash
./test_all_api.sh
```

**Environment Variables:**
- `BASE_URL` - API base URL (default: `http://localhost:8080`)
- `ADMIN_EMAIL` - Admin email (default: `admin@cyclingstream.local`)
- `ADMIN_PASSWORD` - Admin password (default: `admin123`)

### `test_security.sh`
Tests security features including rate limiting, security headers, input validation, and CSRF protection.

**Usage:**
```bash
./test_security.sh
```

### `test_login.sh`
Tests login and authentication functionality.

**Usage:**
```bash
./test_login.sh
```

### `test_chat.sh`
Tests chat functionality and WebSocket connections.

**Usage:**
```bash
./test_chat.sh
```

### `test_analytics.sh`
Tests analytics endpoints.

**Usage:**
```bash
./test_analytics.sh
```

### `test_viewer_tracking.sh`
Tests viewer tracking and session management.

**Usage:**
```bash
./test_viewer_tracking.sh
```

### `test_revenue_integration.sh`
Tests revenue calculation and reporting endpoints.

**Usage:**
```bash
./test_revenue_integration.sh
```

## Shared Utilities

### `test_utils.sh`
Common utilities for test scripts. Source this file in your test scripts:

```bash
source "$(dirname "$0")/test_utils.sh"
```

**Functions:**
- `check_server()` - Verify server is running
- `api_request()` - Make HTTP requests
- `get_admin_token()` - Get admin JWT token
- `get_user_token()` - Create test user and get token
- `pass_test()` - Record passing test
- `fail_test()` - Record failing test
- `print_test_summary()` - Print test results summary

## Writing New Test Scripts

1. Start with shebang and source utilities:
```bash
#!/bin/bash
set -e
source "$(dirname "$0")/test_utils.sh"
```

2. Check prerequisites:
```bash
check_jq
check_server
```

3. Use helper functions:
```bash
response=$(api_request "GET" "/endpoint" "" "Bearer $token")
code=$(get_http_code "$response")
body=$(get_response_body "$response")
```

4. Record test results:
```bash
if [ "$code" = "200" ]; then
    pass_test "GET /endpoint returns 200"
else
    fail_test "GET /endpoint" "Expected 200, got $code"
fi
```

5. Print summary at end:
```bash
print_test_summary
exit $?
```

## Best Practices

1. **Error Handling**: Use `set -e` to exit on errors
2. **Cleanup**: Implement cleanup function for test data
3. **Isolation**: Each test should be independent
4. **Documentation**: Add comments explaining what each test does
5. **Output**: Use colored output for better readability
6. **Validation**: Validate JSON responses before parsing
7. **Rate Limiting**: Add delays between requests if needed
8. **Idempotency**: Tests should be safe to run multiple times

## Troubleshooting

### Server Not Running
```
Error: Backend server is not running on http://localhost:8080
```
**Solution**: Start the backend server with `make run-backend`

### jq Not Found
```
jq is required but not installed
```
**Solution**: Install jq:
- Debian/Ubuntu: `sudo apt install jq`
- macOS: `brew install jq`

### Tests Failing
- Check server logs for errors
- Verify database is running and migrations are applied
- Check environment variables are set correctly
- Review test output for specific error messages

## Continuous Integration

These scripts can be used in CI/CD pipelines. Set appropriate environment variables:

```bash
export BASE_URL="https://api.example.com"
export ADMIN_EMAIL="ci-admin@example.com"
export ADMIN_PASSWORD="secure-password"
./test_all_api.sh
```

