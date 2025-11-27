// +build integration

package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cyclingstream/backend/internal/chat"
	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/cyclingstream/backend/internal/testutil"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestApp creates a Fiber app with chat routes for testing
func setupTestApp(t *testing.T) (*fiber.App, *ChatHandler, func()) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := testutil.GetTestDB(t)
	hub := chat.NewHub()
	go hub.Run()

	chatRepo := repository.NewChatRepository(db)
	raceRepo := repository.NewRaceRepository(db)
	streamRepo := repository.NewStreamRepository(db)
	userRepo := repository.NewUserRepository(db)
	rateLimiter := chat.NewRateLimiter()
	defer rateLimiter.Stop()

	handler := NewChatHandler(chatRepo, raceRepo, streamRepo, userRepo, hub, rateLimiter)

	app := fiber.New()
	
	// Setup routes (simplified for testing)
	app.Get("/races/:id/chat/history", handler.GetChatHistory)
	app.Get("/races/:id/chat/stats", handler.GetChatStats)

	cleanup := func() {
		db.Close()
	}

	return app, handler, cleanup
}

// TestChatHandler_GetChatHistory_Integration tests chat history endpoint with real database
func TestChatHandler_GetChatHistory_Integration(t *testing.T) {
	app, _, cleanup := setupTestApp(t)
	defer cleanup()

	// Setup test data
	db := testutil.GetTestDB(t)
	defer db.Close()

	raceID := testutil.CreateTestRace(t, db, "History Test Race")
	testutil.CreateTestStream(t, db, raceID, "live")
	defer testutil.CleanupRaces(t, db, []string{raceID})

	userID := testutil.CreateTestUser(t, db, "history@test.com", "password123", "History User")
	defer testutil.CleanupUsers(t, db, []string{userID})

	chatRepo := repository.NewChatRepository(db)

	t.Run("Get empty chat history", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/races/"+raceID+"/chat/history", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		messages, ok := result["messages"].([]interface{})
		require.True(t, ok)
		assert.Empty(t, messages)
	})

	t.Run("Get chat history with messages", func(t *testing.T) {
		// Create test messages
		for i := 0; i < 3; i++ {
			msg := &models.ChatMessage{
				RaceID:   raceID,
				UserID:   &userID,
				Username: "History User",
				Message:  "Test message " + string(rune('A'+i)),
			}
			err := chatRepo.Create(msg)
			require.NoError(t, err)
		}

		req := httptest.NewRequest("GET", "/races/"+raceID+"/chat/history", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		messages, ok := result["messages"].([]interface{})
		require.True(t, ok)
		assert.GreaterOrEqual(t, len(messages), 3)
	})

	t.Run("Get chat history with pagination", func(t *testing.T) {
		testutil.CleanupChatMessages(t, db)

		// Create 5 messages
		for i := 0; i < 5; i++ {
			msg := &models.ChatMessage{
				RaceID:   raceID,
				UserID:   &userID,
				Username: "History User",
				Message:  "Pagination message " + string(rune('0'+i)),
			}
			err := chatRepo.Create(msg)
			require.NoError(t, err)
		}

		// Get first 3
		req := httptest.NewRequest("GET", "/races/"+raceID+"/chat/history?limit=3&offset=0", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		messages := result["messages"].([]interface{})
		assert.LessOrEqual(t, len(messages), 3)
		assert.Equal(t, float64(3), result["limit"])
		assert.Equal(t, float64(0), result["offset"])
	})

	t.Run("Invalid race ID returns empty array", func(t *testing.T) {
		invalidRaceID := "00000000-0000-0000-0000-000000000000"
		req := httptest.NewRequest("GET", "/races/"+invalidRaceID+"/chat/history", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		messages := result["messages"].([]interface{})
		assert.Empty(t, messages)
	})
}

// TestChatHandler_GetChatStats_Integration tests chat stats endpoint with real database
func TestChatHandler_GetChatStats_Integration(t *testing.T) {
	app, handler, cleanup := setupTestApp(t)
	defer cleanup()

	// Setup test data
	db := testutil.GetTestDB(t)
	defer db.Close()

	raceID := testutil.CreateTestRace(t, db, "Stats Test Race")
	testutil.CreateTestStream(t, db, raceID, "live")
	defer testutil.CleanupRaces(t, db, []string{raceID})

	userID := testutil.CreateTestUser(t, db, "stats@test.com", "password123", "Stats User")
	defer testutil.CleanupUsers(t, db, []string{userID})

	chatRepo := repository.NewChatRepository(db)

	t.Run("Get stats with no messages", func(t *testing.T) {
		testutil.CleanupChatMessages(t, db)

		req := httptest.NewRequest("GET", "/races/"+raceID+"/chat/stats", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, float64(0), result["total_messages"])
		assert.Equal(t, float64(0), result["concurrent_connections"])
	})

	t.Run("Get stats with messages", func(t *testing.T) {
		testutil.CleanupChatMessages(t, db)

		// Create messages
		for i := 0; i < 5; i++ {
			msg := &models.ChatMessage{
				RaceID:   raceID,
				UserID:   &userID,
				Username: "Stats User",
				Message:  "Stats message " + string(rune('0'+i)),
			}
			err := chatRepo.Create(msg)
			require.NoError(t, err)
		}

		req := httptest.NewRequest("GET", "/races/"+raceID+"/chat/stats", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, float64(5), result["total_messages"])
		// concurrent_connections will be 0 in tests since no WebSocket connections
		assert.Equal(t, float64(0), result["concurrent_connections"])
	})

	t.Run("Stats reflect message count accurately", func(t *testing.T) {
		testutil.CleanupChatMessages(t, db)

		// Add messages one by one and verify count
		for i := 1; i <= 3; i++ {
			msg := &models.ChatMessage{
				RaceID:   raceID,
				UserID:   &userID,
				Username: "Stats User",
				Message:  "Count message",
			}
			err := chatRepo.Create(msg)
			require.NoError(t, err)

			req := httptest.NewRequest("GET", "/races/"+raceID+"/chat/stats", nil)
			resp, err := app.Test(req)
			require.NoError(t, err)

			body, _ := io.ReadAll(resp.Body)
			var result map[string]interface{}
			json.Unmarshal(body, &result)

			assert.Equal(t, float64(i), result["total_messages"])
		}
	})
}

// Note: WebSocket integration tests require a running server and WebSocket client.
// These should be tested with:
// 1. A real Fiber server running on a test port
// 2. WebSocket client library (github.com/gorilla/websocket recommended)
// 3. Or use the enhanced shell script (test_chat_integration.sh) with websocat or similar tool
//
// Example structure for WebSocket integration tests (requires running server):
//
// func TestChatHandler_HandleWebSocket_Integration(t *testing.T) {
//     // Setup: Start test server
//     // Connect WebSocket client
//     // Test connection acceptance/rejection
//     // Test message sending/receiving
//     // Test broadcasting
//     // Test rate limiting
// }

