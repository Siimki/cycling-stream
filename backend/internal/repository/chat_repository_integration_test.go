// +build integration

package repository

import (
	"testing"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestChatRepository_Integration tests chat repository with real database
func TestChatRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := testutil.GetTestDB(t)
	defer db.Close()

	repo := NewChatRepository(db)
	raceID := testutil.CreateTestRace(t, db, "Chat Test Race")
	defer testutil.CleanupRaces(t, db, []string{raceID})

	userID := testutil.CreateTestUser(t, db, "test@example.com", "hashedpassword", "Test User")
	defer testutil.CleanupUsers(t, db, []string{userID})

	t.Run("Create message with authenticated user", func(t *testing.T) {
		msg := &models.ChatMessage{
			RaceID:   raceID,
			UserID:   &userID,
			Username: "Test User",
			Message:  "Hello, world!",
		}

		err := repo.Create(msg)
		require.NoError(t, err)
		assert.NotEmpty(t, msg.ID)
		assert.NotZero(t, msg.CreatedAt)
		assert.Equal(t, raceID, msg.RaceID)
		assert.Equal(t, userID, *msg.UserID)
		assert.Equal(t, "Test User", msg.Username)
		assert.Equal(t, "Hello, world!", msg.Message)
	})

	t.Run("Create anonymous message", func(t *testing.T) {
		msg := &models.ChatMessage{
			RaceID:   raceID,
			UserID:   nil,
			Username: "Anonymous",
			Message:  "Anonymous message",
		}

		err := repo.Create(msg)
		require.NoError(t, err)
		assert.NotEmpty(t, msg.ID)
		assert.Nil(t, msg.UserID)
		assert.Equal(t, "Anonymous", msg.Username)
	})

	t.Run("Get messages with pagination", func(t *testing.T) {
		// Create multiple messages
		for i := 0; i < 5; i++ {
			msg := &models.ChatMessage{
				RaceID:   raceID,
				UserID:   &userID,
				Username: "Test User",
				Message:  "Message " + string(rune('A'+i)),
			}
			err := repo.Create(msg)
			require.NoError(t, err)
		}

		// Get first 3 messages
		messages, err := repo.GetByRaceID(raceID, 3, 0)
		require.NoError(t, err)
		assert.Len(t, messages, 3)

		// Get next 3 messages (with offset)
		messages2, err := repo.GetByRaceID(raceID, 3, 3)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(messages2), 2) // At least 2 more from previous test

		// Verify no overlap
		for _, m1 := range messages {
			for _, m2 := range messages2 {
				assert.NotEqual(t, m1.ID, m2.ID)
			}
		}
	})

	t.Run("Get messages in chronological order", func(t *testing.T) {
		// Clean previous messages for this test
		testutil.CleanupChatMessages(t, db)

		// Create messages with slight delays
		var messageIDs []string
		for i := 0; i < 3; i++ {
			msg := &models.ChatMessage{
				RaceID:   raceID,
				UserID:   &userID,
				Username: "Test User",
				Message:  "Ordered message " + string(rune('A'+i)),
			}
			err := repo.Create(msg)
			require.NoError(t, err)
			messageIDs = append(messageIDs, msg.ID)
		}

		// Get messages
		messages, err := repo.GetByRaceID(raceID, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(messages), 3)

		// Verify chronological order (oldest first)
		for i := 1; i < len(messages); i++ {
			assert.True(t, messages[i].CreatedAt.After(messages[i-1].CreatedAt) || messages[i].CreatedAt.Equal(messages[i-1].CreatedAt),
				"Messages should be in chronological order")
		}
	})

	t.Run("Get recent messages", func(t *testing.T) {
		testutil.CleanupChatMessages(t, db)

		// Create 10 messages
		for i := 0; i < 10; i++ {
			msg := &models.ChatMessage{
				RaceID:   raceID,
				UserID:   &userID,
				Username: "Test User",
				Message:  "Recent message " + string(rune('0'+i)),
			}
			err := repo.Create(msg)
			require.NoError(t, err)
		}

		// Get 5 most recent
		messages, err := repo.GetRecentByRaceID(raceID, 5)
		require.NoError(t, err)
		assert.Len(t, messages, 5)

		// Verify they're in chronological order (oldest first after reversal)
		for i := 1; i < len(messages); i++ {
			assert.True(t, messages[i].CreatedAt.After(messages[i-1].CreatedAt) || messages[i].CreatedAt.Equal(messages[i-1].CreatedAt))
		}
	})

	t.Run("Count messages", func(t *testing.T) {
		testutil.CleanupChatMessages(t, db)

		// Initially should be 0
		count, err := repo.CountByRaceID(raceID)
		require.NoError(t, err)
		assert.Equal(t, 0, count)

		// Create 3 messages
		for i := 0; i < 3; i++ {
			msg := &models.ChatMessage{
				RaceID:   raceID,
				UserID:   &userID,
				Username: "Test User",
				Message:  "Count message " + string(rune('A'+i)),
			}
			err := repo.Create(msg)
			require.NoError(t, err)
		}

		// Count should be 3
		count, err = repo.CountByRaceID(raceID)
		require.NoError(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("Get empty result for non-existent race", func(t *testing.T) {
		nonExistentRaceID := uuid.New().String()
		messages, err := repo.GetByRaceID(nonExistentRaceID, 10, 0)
		require.NoError(t, err)
		assert.Empty(t, messages)

		count, err := repo.CountByRaceID(nonExistentRaceID)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("Messages are race-specific", func(t *testing.T) {
		testutil.CleanupChatMessages(t, db)

		raceID2 := testutil.CreateTestRace(t, db, "Another Test Race")
		defer testutil.CleanupRaces(t, db, []string{raceID2})

		// Create message for race1
		msg1 := &models.ChatMessage{
			RaceID:   raceID,
			UserID:   &userID,
			Username: "Test User",
			Message:  "Race 1 message",
		}
		err := repo.Create(msg1)
		require.NoError(t, err)

		// Create message for race2
		msg2 := &models.ChatMessage{
			RaceID:   raceID2,
			UserID:   &userID,
			Username: "Test User",
			Message:  "Race 2 message",
		}
		err = repo.Create(msg2)
		require.NoError(t, err)

		// Get messages for race1 - should only get race1 message
		messages, err := repo.GetByRaceID(raceID, 10, 0)
		require.NoError(t, err)
		assert.Len(t, messages, 1)
		assert.Equal(t, raceID, messages[0].RaceID)
		assert.Equal(t, "Race 1 message", messages[0].Message)

		// Get messages for race2 - should only get race2 message
		messages2, err := repo.GetByRaceID(raceID2, 10, 0)
		require.NoError(t, err)
		assert.Len(t, messages2, 1)
		assert.Equal(t, raceID2, messages2[0].RaceID)
		assert.Equal(t, "Race 2 message", messages2[0].Message)
	})

	t.Run("Create message with invalid race ID is handled", func(t *testing.T) {
		invalidRaceID := "00000000-0000-0000-0000-000000000000"
		msg := &models.ChatMessage{
			RaceID:   invalidRaceID,
			UserID:   &userID,
			Username: "Test User",
			Message:  "Message for invalid race",
		}

		// Should either succeed (if foreign key constraint is not enforced)
		// or fail gracefully (if constraint is enforced)
		err := repo.Create(msg)
		// Both outcomes are acceptable - the key is that it doesn't panic
		if err != nil {
			// If foreign key constraint is enforced, error is expected
			assert.NotNil(t, err)
		} else {
			// If no constraint, message should be created
			assert.NotEmpty(t, msg.ID)
		}
	})

	t.Run("GetByRaceID with invalid UUID format is handled", func(t *testing.T) {
		invalidRaceID := "not-a-valid-uuid"
		messages, err := repo.GetByRaceID(invalidRaceID, 10, 0)

		// Should either return empty array or error, but not panic
		if err != nil {
			assert.NotNil(t, err)
		} else {
			assert.Empty(t, messages)
		}
	})

	t.Run("CountByRaceID with invalid UUID format is handled", func(t *testing.T) {
		invalidRaceID := "not-a-valid-uuid"
		count, err := repo.CountByRaceID(invalidRaceID)

		// Should either return 0 or error, but not panic
		if err != nil {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, 0, count)
		}
	})
}

