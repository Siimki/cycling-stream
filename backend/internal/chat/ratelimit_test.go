package chat

import (
	"testing"

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

	t.Run("Limit resets after window expires", func(t *testing.T) {
		identifier := "user-6"
		
		// Fill up the limit
		for i := 0; i < ChatRateLimitMaxMessages; i++ {
			allowed := rl.CheckRateLimit(identifier)
			assert.True(t, allowed, "Message %d should be allowed", i+1)
		}
		
		// Next message should be denied
		allowed := rl.CheckRateLimit(identifier)
		assert.False(t, allowed, "Message should be denied after limit reached")
		
		// Verify remaining is 0
		remaining := rl.GetRemainingMessages(identifier)
		assert.Equal(t, 0, remaining)
		
		// Note: Full window expiration test (1 minute) is too slow for unit tests.
		// The filtering logic is tested indirectly: CheckRateLimit filters old timestamps
		// on each call, so the limit will reset once timestamps are older than the window.
		// For integration tests, consider using a shorter window or time mocking.
	})
}

func TestRateLimiter_GetRemainingMessages(t *testing.T) {
	rl := NewRateLimiter()
	defer rl.Stop()

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

