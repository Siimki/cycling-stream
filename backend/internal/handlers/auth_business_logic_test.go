package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cyclingstream/backend/internal/testutil"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthHandler_Register tests registration business logic
func TestAuthHandler_Register(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	app, _, cleanup := setupAuthTestApp(t)
	defer cleanup()

	db := testutil.GetTestDB(t)
	defer db.Close()

	t.Run("Register with duplicate email is rejected", func(t *testing.T) {
		email := "duplicate@test.com"
		password := "TestPassword123!"

		// Register first user
		registerReq := RegisterRequest{
			Email:    email,
			Password: password,
		}
		reqBody, err := json.Marshal(registerReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		// Try to register again with same email
		req2 := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
		req2.Header.Set("Content-Type", "application/json")
		resp2, err := app.Test(req2)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusConflict, resp2.StatusCode)

		body, _ := io.ReadAll(resp2.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		assert.Equal(t, "User with this email already exists", result["error"])

		// Cleanup
		testutil.CleanupUsers(t, db, []string{email})
	})

	t.Run("Password validation rules are enforced", func(t *testing.T) {
		testCases := []struct {
			name          string
			password      string
			expectedError string
		}{
			{"Too short", "Short1!", "Password must be at least"},
			{"No uppercase", "lowercase123!", "Password must contain"},
			{"No lowercase", "UPPERCASE123!", "Password must contain"},
			{"No number", "NoNumber!", "Password must contain"},
			{"No special char", "NoSpecial123", "Password must contain"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				registerReq := RegisterRequest{
					Email:    "validation@test.com",
					Password: tc.password,
				}
				reqBody, err := json.Marshal(registerReq)
				require.NoError(t, err)

				req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				resp, err := app.Test(req)
				require.NoError(t, err)

				assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

				body, _ := io.ReadAll(resp.Body)
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				assert.Contains(t, result["error"].(string), tc.expectedError)
			})
		}
	})

	t.Run("Email validation edge cases", func(t *testing.T) {
		testCases := []struct {
			name     string
			email    string
			shouldPass bool
		}{
			{"Valid email", "valid@example.com", true},
			{"Invalid format", "not-an-email", false},
			{"Missing @", "invalidemail.com", false},
			{"Missing domain", "invalid@", false},
			{"Missing local part", "@example.com", false},
			{"Multiple @", "invalid@@example.com", false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				registerReq := RegisterRequest{
					Email:    tc.email,
					Password: "TestPassword123!",
				}
				reqBody, err := json.Marshal(registerReq)
				require.NoError(t, err)

				req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				resp, err := app.Test(req)
				require.NoError(t, err)

				if tc.shouldPass {
					assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
					// Cleanup
					testutil.CleanupUsers(t, db, []string{tc.email})
				} else {
					assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
				}
			})
		}
	})

	t.Run("Bio length limit is enforced", func(t *testing.T) {
		// Bio should be max 120 characters
		bio120 := string(make([]byte, 120))
		for i := range bio120 {
			bio120 = bio120[:i] + "a" + bio120[i+1:]
		}
		bio121 := bio120 + "a" // 121 characters

		registerReq := RegisterRequest{
			Email:    "biolimit@test.com",
			Password: "TestPassword123!",
			Bio:      &bio121,
		}
		reqBody, err := json.Marshal(registerReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		if resp.StatusCode == fiber.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			var result AuthResponse
			json.Unmarshal(body, &result)
			if result.User != nil && result.User.Bio != "" {
				assert.LessOrEqual(t, len(result.User.Bio), 120,
					"Bio should be truncated to 120 characters")
			}
			testutil.CleanupUsers(t, db, []string{"biolimit@test.com"})
		}
	})
}

// TestAuthHandler_ChangePassword tests password change business logic
func TestAuthHandler_ChangePassword(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	app, authHandler, cleanup := setupAuthTestApp(t)
	defer cleanup()

	db := testutil.GetTestDB(t)
	defer db.Close()

	// Create a test user
	email := "changepass@test.com"
	password := "TestPassword123!"
	userID := testutil.CreateTestUser(t, db, email, password, "Change Password User")
	defer testutil.CleanupUsers(t, db, []string{userID})

	// Get a valid token
	loginReq := LoginRequest{Email: email, Password: password}
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	body, _ := io.ReadAll(resp.Body)
	var loginResult AuthResponse
	json.Unmarshal(body, &loginResult)
	token := loginResult.Token

	// Setup middleware and route
	userAuthMiddleware := func(c *fiber.Ctx) error {
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
		jwtToken, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("test-secret-key-for-testing-only"), nil
		})
		if err != nil || !jwtToken.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}
		c.Locals("user_id", claims["user_id"])
		c.Locals("is_admin", claims["is_admin"])
		return c.Next()
	}

	app.Post("/users/password", userAuthMiddleware, authHandler.ChangePassword)

	t.Run("Change password with wrong current password is rejected", func(t *testing.T) {
		changePasswordReq := ChangePasswordRequest{
			CurrentPassword: "WrongPassword123!",
			NewPassword:     "NewPassword123!",
		}
		reqBody, err := json.Marshal(changePasswordReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/users/password", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		assert.Equal(t, "Current password is incorrect", result["error"])
	})

	t.Run("Change password with same password as new password is allowed", func(t *testing.T) {
		// This is technically allowed - user can "change" to the same password
		// Though in practice, you might want to reject this
		changePasswordReq := ChangePasswordRequest{
			CurrentPassword: password,
			NewPassword:     password, // Same password
		}
		reqBody, err := json.Marshal(changePasswordReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/users/password", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		require.NoError(t, err)

		// Should succeed (same password is technically valid)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Change password validates new password rules", func(t *testing.T) {
		testCases := []struct {
			name          string
			newPassword   string
			expectedError string
		}{
			{"Too short", "Short1!", "Password must be at least"},
			{"No uppercase", "lowercase123!", "Password must contain"},
			{"No number", "NoNumber!", "Password must contain"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				changePasswordReq := ChangePasswordRequest{
					CurrentPassword: password,
					NewPassword:     tc.newPassword,
				}
				reqBody, err := json.Marshal(changePasswordReq)
				require.NoError(t, err)

				req := httptest.NewRequest("POST", "/users/password", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)
				resp, err := app.Test(req)
				require.NoError(t, err)

				assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

				body, _ := io.ReadAll(resp.Body)
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				assert.Contains(t, result["error"].(string), tc.expectedError)
			})
		}
	})

	t.Run("Change password with valid credentials succeeds", func(t *testing.T) {
		newPassword := "NewValidPassword123!"
		changePasswordReq := ChangePasswordRequest{
			CurrentPassword: password,
			NewPassword:     newPassword,
		}
		reqBody, err := json.Marshal(changePasswordReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/users/password", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Verify new password works
		loginReq := LoginRequest{Email: email, Password: newPassword}
		reqBody2, _ := json.Marshal(loginReq)
		req2 := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody2))
		req2.Header.Set("Content-Type", "application/json")
		resp2, _ := app.Test(req2)
		assert.Equal(t, fiber.StatusOK, resp2.StatusCode)
	})
}

// TestAuthHandler_AwardBonusPoints tests bonus points business logic
func TestAuthHandler_AwardBonusPoints(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	app, authHandler, cleanup := setupAuthTestApp(t)
	defer cleanup()

	db := testutil.GetTestDB(t)
	defer db.Close()

	// Create a test user
	email := "points@test.com"
	password := "TestPassword123!"
	userID := testutil.CreateTestUser(t, db, email, password, "Points User")
	defer testutil.CleanupUsers(t, db, []string{userID})

	// Get a valid token
	loginReq := LoginRequest{Email: email, Password: password}
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	body, _ := io.ReadAll(resp.Body)
	var loginResult AuthResponse
	json.Unmarshal(body, &loginResult)
	token := loginResult.Token

	// Setup middleware and route
	userAuthMiddleware := func(c *fiber.Ctx) error {
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
		jwtToken, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("test-secret-key-for-testing-only"), nil
		})
		if err != nil || !jwtToken.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}
		c.Locals("user_id", claims["user_id"])
		c.Locals("is_admin", claims["is_admin"])
		return c.Next()
	}

	app.Post("/users/bonus-points", userAuthMiddleware, authHandler.AwardBonusPoints)

	t.Run("AwardBonusPoints adds points to user", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/users/bonus-points", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		assert.Equal(t, "Bonus points awarded", result["message"])
		assert.Equal(t, float64(50), result["bonus_points"])
		assert.GreaterOrEqual(t, result["total_points"].(float64), float64(50))
	})

	t.Run("AwardBonusPoints can be called multiple times", func(t *testing.T) {
		// Award points multiple times
		for i := 0; i < 3; i++ {
			req := httptest.NewRequest("POST", "/users/bonus-points", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		}

		// Verify points accumulated
		req := httptest.NewRequest("POST", "/users/bonus-points", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		require.NoError(t, err)

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		// Should have at least 200 points (50 * 4 calls)
		assert.GreaterOrEqual(t, result["total_points"].(float64), float64(200))
	})
}

