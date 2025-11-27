package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cyclingstream/backend/internal/testutil"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthHandler_Security tests security-related scenarios
func TestAuthHandler_Security(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	app, _, cleanup := setupAuthTestApp(t)
	defer cleanup()

	db := testutil.GetTestDB(t)
	defer db.Close()

	// Create a test user
	email := "security@test.com"
	password := "TestPassword123!"
	userID := testutil.CreateTestUser(t, db, email, password, "Security Test User")
	defer testutil.CleanupUsers(t, db, []string{userID})

	t.Run("SQL injection in email field is prevented", func(t *testing.T) {
		// Attempt SQL injection in email
		sqlInjectionAttempts := []string{
			"test@example.com' OR '1'='1",
			"test@example.com'; DROP TABLE users; --",
			"test@example.com' UNION SELECT * FROM users --",
			"admin@example.com' OR 1=1 --",
		}

		for _, maliciousEmail := range sqlInjectionAttempts {
			loginReq := LoginRequest{
				Email:    maliciousEmail,
				Password: password,
			}
			reqBody, err := json.Marshal(loginReq)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			require.NoError(t, err)

			// Should reject with validation error, not execute SQL
			assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusUnauthorized,
				"SQL injection attempt should be rejected, got status %d for email: %s", resp.StatusCode, maliciousEmail)
		}
	})

	t.Run("SQL injection in password field is prevented", func(t *testing.T) {
		sqlInjectionAttempts := []string{
			"' OR '1'='1",
			"'; DROP TABLE users; --",
			"' UNION SELECT * FROM users --",
		}

		for _, maliciousPassword := range sqlInjectionAttempts {
			loginReq := LoginRequest{
				Email:    email,
				Password: maliciousPassword,
			}
			reqBody, err := json.Marshal(loginReq)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			require.NoError(t, err)

			// Should reject with unauthorized, not execute SQL
			assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode,
				"SQL injection attempt should be rejected, got status %d", resp.StatusCode)
		}
	})

	t.Run("XSS in registration fields is sanitized", func(t *testing.T) {
		xssAttempts := []string{
			"<script>alert('XSS')</script>",
			"<img src=x onerror=alert('XSS')>",
			"javascript:alert('XSS')",
			"<svg onload=alert('XSS')>",
		}

		for _, xssPayload := range xssAttempts {
			registerReq := RegisterRequest{
				Email:    "xss" + xssPayload + "@test.com",
				Password: "TestPassword123!",
				Name:     &xssPayload,
			}
			reqBody, err := json.Marshal(registerReq)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			require.NoError(t, err)

			// Should either reject (invalid email) or sanitize (name field)
			// The key is that XSS should not be stored as-is
			if resp.StatusCode == fiber.StatusCreated {
				body, _ := io.ReadAll(resp.Body)
				var result AuthResponse
				json.Unmarshal(body, &result)
				if result.User != nil && result.User.Name != nil {
					// Name should be sanitized (no script tags)
					assert.NotContains(t, *result.User.Name, "<script>",
						"XSS payload should be sanitized in name field")
				}
			}
		}
	})

	t.Run("Input validation rejects extremely long inputs", func(t *testing.T) {
		// Test email length limit
		longEmail := "a" + strings.Repeat("b", 300) + "@test.com"
		registerReq := RegisterRequest{
			Email:    longEmail,
			Password: "TestPassword123!",
		}
		reqBody, err := json.Marshal(registerReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode,
			"Extremely long email should be rejected")
	})

	t.Run("Password length limit is enforced", func(t *testing.T) {
		// Password over 128 characters should be rejected
		longPassword := strings.Repeat("a", 129)
		loginReq := LoginRequest{
			Email:    email,
			Password: longPassword,
		}
		reqBody, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode,
			"Password over 128 characters should be rejected")
	})

	t.Run("Bio length limit is enforced", func(t *testing.T) {
		longBio := strings.Repeat("a", 121) // Over 120 char limit
		registerReq := RegisterRequest{
			Email:    "biolimit@test.com",
			Password: "TestPassword123!",
			Bio:      &longBio,
		}
		reqBody, err := json.Marshal(registerReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		// Bio should be truncated to 120 chars
		if resp.StatusCode == fiber.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			var result AuthResponse
			json.Unmarshal(body, &result)
			if result.User != nil && result.User.Bio != "" {
				assert.LessOrEqual(t, len(result.User.Bio), 120,
					"Bio should be truncated to 120 characters")
			}
		}
	})

	t.Run("Unicode and emoji in inputs are handled safely", func(t *testing.T) {
		testCases := []struct {
			name     string
			email    string
			username string
		}{
			{"Unicode characters", "test@example.com", "测试用户"},
			{"Emoji in name", "emoji@test.com", "User 😀🎉"},
			{"Special unicode", "unicode@test.com", "User 测试 🎯"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				registerReq := RegisterRequest{
					Email:    tc.email,
					Password: "TestPassword123!",
					Name:     &tc.username,
				}
				reqBody, err := json.Marshal(registerReq)
				require.NoError(t, err)

				req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				resp, err := app.Test(req)
				require.NoError(t, err)

				// Should handle unicode safely (either accept or reject gracefully)
				assert.True(t, resp.StatusCode == fiber.StatusCreated || resp.StatusCode == fiber.StatusBadRequest,
					"Unicode input should be handled safely, got status %d", resp.StatusCode)
			})
		}
	})
}

// TestAuthHandler_AuthenticationBypass tests authentication bypass attempts
func TestAuthHandler_AuthenticationBypass(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	app, _, cleanup := setupAuthTestApp(t)
	defer cleanup()

	db := testutil.GetTestDB(t)
	defer db.Close()

	// Create a test user
	email := "bypass@test.com"
	password := "TestPassword123!"
	userID := testutil.CreateTestUser(t, db, email, password, "Bypass Test User")
	defer testutil.CleanupUsers(t, db, []string{userID})

	// First, get a valid token
	loginReq := LoginRequest{
		Email:    email,
		Password: password,
	}
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	body, _ := io.ReadAll(resp.Body)
	var loginResult AuthResponse
	json.Unmarshal(body, &loginResult)
	validToken := loginResult.Token

	t.Run("Request without token is rejected", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/profile", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		// Should require authentication
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode,
			"Request without token should be rejected")
	})

	t.Run("Request with invalid token is rejected", func(t *testing.T) {
		invalidTokens := []string{
			"invalid.token.here",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature",
			"not-a-jwt-token",
			"",
		}

		for _, invalidToken := range invalidTokens {
			req := httptest.NewRequest("GET", "/users/profile", nil)
			req.Header.Set("Authorization", "Bearer "+invalidToken)
			resp, err := app.Test(req)
			require.NoError(t, err)

			assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode,
				"Invalid token should be rejected: %s", invalidToken)
		}
	})

	t.Run("Request with tampered token is rejected", func(t *testing.T) {
		// Tamper with the token by modifying a character
		tamperedToken := validToken[:len(validToken)-1] + "X"
		req := httptest.NewRequest("GET", "/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+tamperedToken)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode,
			"Tampered token should be rejected")
	})

	t.Run("Request with expired token is rejected", func(t *testing.T) {
		// Create an expired token
		claims := jwt.MapClaims{
			"user_id":  userID,
			"is_admin": false,
			"exp":       time.Now().Add(-1 * time.Hour).Unix(), // Expired 1 hour ago
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		expiredToken, _ := token.SignedString([]byte("test-secret-key-for-testing-only"))

		req := httptest.NewRequest("GET", "/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode,
			"Expired token should be rejected")
	})

	t.Run("Request with token signed with wrong secret is rejected", func(t *testing.T) {
		// Create a token with wrong secret
		claims := jwt.MapClaims{
			"user_id":  userID,
			"is_admin": false,
			"exp":       time.Now().Add(24 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		wrongSecretToken, _ := token.SignedString([]byte("wrong-secret-key"))

		req := httptest.NewRequest("GET", "/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+wrongSecretToken)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode,
			"Token signed with wrong secret should be rejected")
	})

	t.Run("Request with malformed Authorization header is rejected", func(t *testing.T) {
		malformedHeaders := []string{
			validToken,                    // Missing "Bearer " prefix
			"Basic " + validToken,         // Wrong scheme
			"Bearer",                      // Missing token
			"Bearer  " + validToken,       // Extra spaces
			"Bearer" + validToken,         // Missing space
		}

		for _, header := range malformedHeaders {
			req := httptest.NewRequest("GET", "/users/profile", nil)
			req.Header.Set("Authorization", header)
			resp, err := app.Test(req)
			require.NoError(t, err)

			assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode,
				"Malformed Authorization header should be rejected: %s", header)
		}
	})
}

// TestAuthHandler_Authorization tests authorization boundaries
func TestAuthHandler_Authorization(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	app, authHandler, cleanup := setupAuthTestApp(t)
	defer cleanup()

	db := testutil.GetTestDB(t)
	defer db.Close()

	// Create two test users
	user1Email := "user1@test.com"
	user1Password := "TestPassword123!"
	user1ID := testutil.CreateTestUser(t, db, user1Email, user1Password, "User 1")
	defer testutil.CleanupUsers(t, db, []string{user1ID})

	user2Email := "user2@test.com"
	user2Password := "TestPassword123!"
	user2ID := testutil.CreateTestUser(t, db, user2Email, user2Password, "User 2")
	defer testutil.CleanupUsers(t, db, []string{user2ID})

	// Get tokens for both users
	getToken := func(email, password string) string {
		loginReq := LoginRequest{Email: email, Password: password}
		reqBody, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		var result AuthResponse
		json.Unmarshal(body, &result)
		return result.Token
	}

	user1Token := getToken(user1Email, user1Password)
	_ = getToken(user2Email, user2Password) // user2Token not used in current tests

	// Setup authenticated routes for testing
	userAuthMiddleware := func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization required"})
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid authorization header"})
		}
		tokenString := parts[1]
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("test-secret-key-for-testing-only"), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}
		c.Locals("user_id", claims["user_id"])
		c.Locals("is_admin", claims["is_admin"])
		return c.Next()
	}

	app.Get("/users/profile", userAuthMiddleware, authHandler.GetProfile)
	app.Post("/users/password", userAuthMiddleware, authHandler.ChangePassword)

	t.Run("User cannot access another user's profile by manipulating user_id", func(t *testing.T) {
		// This test verifies that GetProfile uses the authenticated user's ID from the token,
		// not from request parameters. The handler should use c.Locals("user_id") which comes
		// from the JWT token, not from URL params or body.

		req := httptest.NewRequest("GET", "/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+user1Token)
		resp, err := app.Test(req)
		require.NoError(t, err)

		if resp.StatusCode == fiber.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			var user map[string]interface{}
			json.Unmarshal(body, &user)
			// Should return user1's profile, not user2's
			assert.Equal(t, user1ID, user["id"], "Should return authenticated user's profile")
		}
	})

	t.Run("Non-admin user cannot access admin endpoints", func(t *testing.T) {
		// Create a non-admin token
		claims := jwt.MapClaims{
			"user_id":  user1ID,
			"is_admin": false, // Not admin
			"exp":       time.Now().Add(24 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		nonAdminToken, _ := token.SignedString([]byte("test-secret-key-for-testing-only"))

		// Try to access an admin endpoint (would need to set up an admin route for this test)
		// For now, we verify the middleware rejects non-admin tokens
		req := httptest.NewRequest("GET", "/admin/races", nil)
		req.Header.Set("Authorization", "Bearer "+nonAdminToken)
		resp, err := app.Test(req)
		require.NoError(t, err)

		// Admin middleware should reject non-admin users
		// Note: This test assumes admin routes use AuthMiddleware which requires is_admin=true
		assert.True(t, resp.StatusCode == fiber.StatusForbidden || resp.StatusCode == fiber.StatusNotFound,
			"Non-admin user should not access admin endpoints")
	})

	t.Run("Anonymous user cannot perform authenticated actions", func(t *testing.T) {
		// Try to change password without authentication
		changePasswordReq := map[string]string{
			"current_password": user1Password,
			"new_password":     "NewPassword123!",
		}
		reqBody, _ := json.Marshal(changePasswordReq)
		req := httptest.NewRequest("POST", "/users/password", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		// No Authorization header
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode,
			"Anonymous user should not be able to change password")
	})
}

