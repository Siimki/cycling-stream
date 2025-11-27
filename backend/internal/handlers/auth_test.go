package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/cyclingstream/backend/internal/repository"
	"github.com/cyclingstream/backend/internal/testutil"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupAuthTestApp creates a Fiber app with auth routes for testing
func setupAuthTestApp(t *testing.T) (*fiber.App, *AuthHandler, func()) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := testutil.GetTestDB(t)
	userRepo := repository.NewUserRepository(db)
	authHandler := NewAuthHandler(userRepo, "test-secret-key-for-testing-only")

	app := fiber.New()
	app.Post("/auth/login", authHandler.Login)
	app.Post("/auth/register", authHandler.Register)

	cleanup := func() {
		db.Close()
	}

	return app, authHandler, cleanup
}

// TestAuthHandler_Login tests login functionality
func TestAuthHandler_Login(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	app, _, cleanup := setupAuthTestApp(t)
	defer cleanup()

	db := testutil.GetTestDB(t)
	defer db.Close()

	// Create a test user
	email := "testuser@example.com"
	password := "TestPassword123!"
	userID := testutil.CreateTestUser(t, db, email, password, "Test User")
	defer testutil.CleanupUsers(t, db, []string{userID})

	t.Run("Successful login with valid credentials", func(t *testing.T) {
		loginReq := LoginRequest{
			Email:    email,
			Password: password,
		}
		reqBody, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result AuthResponse
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.NotEmpty(t, result.Token)
		assert.NotNil(t, result.User)
		assert.Equal(t, email, result.User.Email)
		assert.Empty(t, result.User.PasswordHash, "Password hash should not be returned")
	})

	t.Run("Login fails with wrong password", func(t *testing.T) {
		loginReq := LoginRequest{
			Email:    email,
			Password: "WrongPassword123!",
		}
		reqBody, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, "Invalid credentials", result["error"])
	})

	t.Run("Login fails with non-existent email", func(t *testing.T) {
		loginReq := LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "SomePassword123!",
		}
		reqBody, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, "Invalid credentials", result["error"])
	})

	t.Run("Login fails with empty email", func(t *testing.T) {
		loginReq := LoginRequest{
			Email:    "",
			Password: password,
		}
		reqBody, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, "Email is required", result["error"])
	})

	t.Run("Login fails with empty password", func(t *testing.T) {
		loginReq := LoginRequest{
			Email:    email,
			Password: "",
		}
		reqBody, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, "Password is required", result["error"])
	})

	t.Run("Login fails with invalid email format", func(t *testing.T) {
		loginReq := LoginRequest{
			Email:    "not-an-email",
			Password: password,
		}
		reqBody, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, "Invalid email format", result["error"])
	})

	t.Run("Login fails with password too long", func(t *testing.T) {
		longPassword := string(make([]byte, 129)) // 129 characters
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

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, "Invalid credentials", result["error"])
	})

	t.Run("Login fails with invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, "Invalid request body", result["error"])
	})

	t.Run("Admin login with hardcoded credentials", func(t *testing.T) {
		loginReq := LoginRequest{
			Email:    "admin@cyclingstream.local",
			Password: "admin123",
		}
		reqBody, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result AuthResponse
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.NotEmpty(t, result.Token)
		assert.NotNil(t, result.User)
		assert.Equal(t, "admin@cyclingstream.local", result.User.Email)
		assert.Equal(t, "admin", result.User.ID)
	})

	t.Run("Admin login fails with wrong password", func(t *testing.T) {
		loginReq := LoginRequest{
			Email:    "admin@cyclingstream.local",
			Password: "wrongpassword",
		}
		reqBody, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		// Should check database and fail
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}

