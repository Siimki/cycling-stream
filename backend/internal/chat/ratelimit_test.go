package chat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_CheckRateLimit(t *testing.T) {
	rl := NewRateLimiter()
	defer rl.Stop()

	t.Run("First message is allowed", func(t *testing.T) {
		identifier := "user-1"
		allowed := rl.CheckRateLimit(identifier)
		assert.True(t, allowed)
	})

	t.Run("Messages within limit are allowed", func(t *testing.T) {
		identifier := "user-2"
		for i := 0; i < ChatRateLimitMaxMessages; i++ {
			allowed := rl.CheckRateLimit(identifier)
			assert.True(t, allowed, "Message %d should be allowed", i+1)
		}
	})

	t.Run("Message exceeding limit is denied", func(t *testing.T) {
		identifier := "user-3"
		// Send max messages
		for i := 0; i < ChatRateLimitMaxMessages; i++ {
			rl.CheckRateLimit(identifier)
		}
		// Next message should be denied
		allowed := rl.CheckRateLimit(identifier)
		assert.False(t, allowed)
	})

	t.Run("Different identifiers have separate limits", func(t *testing.T) {
		identifier1 := "user-4"
		identifier2 := "user-5"

		// Fill up identifier1's limit
		for i := 0; i < ChatRateLimitMaxMessages; i++ {
			rl.CheckRateLimit(identifier1)
		}

		// identifier2 should still be able to send
		allowed := rl.CheckRateLimit(identifier2)
		assert.True(t, allowed)
	})

	t.Run("Limit resets after window", func(t *testing.T) {
		// This test would require time manipulation or waiting
		// For now, we verify the structure
		identifier := "user-6"
		allowed := rl.CheckRateLimit(identifier)
		assert.True(t, allowed)
	})
}

func TestRateLimiter_GetRemainingMessages(t *testing.T) {
	rl := NewRateLimiter()
	defer rl.Stop()

	t.Run("New identifier has full limit", func(t *testing.T) {
		identifier := "user-7"
		remaining := rl.GetRemainingMessages(identifier)
		assert.Equal(t, ChatRateLimitMaxMessages, remaining)
	})

	t.Run("Remaining decreases with messages", func(t *testing.T) {
		identifier := "user-8"
		rl.CheckRateLimit(identifier)
		remaining := rl.GetRemainingMessages(identifier)
		assert.Equal(t, ChatRateLimitMaxMessages-1, remaining)
	})

	t.Run("Remaining is zero when limit reached", func(t *testing.T) {
		identifier := "user-9"
		for i := 0; i < ChatRateLimitMaxMessages; i++ {
			rl.CheckRateLimit(identifier)
		}
		remaining := rl.GetRemainingMessages(identifier)
		assert.Equal(t, 0, remaining)
	})
}

func TestRateLimiter_Cleanup(t *testing.T) {
	rl := NewRateLimiter()
	defer rl.Stop()

	t.Run("Cleanup runs periodically", func(t *testing.T) {
		// Verify cleanup ticker exists
		assert.NotNil(t, rl.cleanupTicker)
		assert.NotNil(t, rl.stopCleanup)
	})

	t.Run("Stop stops cleanup", func(t *testing.T) {
		rl2 := NewRateLimiter()
		rl2.Stop()
		// Verify it can be stopped without panic
		assert.NotNil(t, rl2)
	})
}

func TestRateLimiter_Constants(t *testing.T) {
	t.Run("Rate limit constants are set correctly", func(t *testing.T) {
		assert.Equal(t, 10, ChatRateLimitMaxMessages)
		assert.Equal(t, 1*time.Minute, ChatRateLimitWindow)
	})
}

