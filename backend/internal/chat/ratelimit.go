package chat

import (
	"sync"
	"time"
)

const (
	// ChatRateLimitMaxMessages is the maximum number of messages allowed per window
	ChatRateLimitMaxMessages = 100 // Increased from 10 to 100 for testing
	// ChatRateLimitWindow is the time window for rate limiting
	ChatRateLimitWindow = 1 * time.Minute
)

// RateLimiter tracks message rates per user/IP
type RateLimiter struct {
	// Map of identifier (user_id or IP) to timestamps of messages
	entries map[string][]time.Time
	mu      sync.RWMutex
	// Cleanup ticker
	cleanupTicker *time.Ticker
	stopCleanup   chan bool
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		entries:       make(map[string][]time.Time),
		cleanupTicker: time.NewTicker(5 * time.Minute),
		stopCleanup:   make(chan bool),
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// cleanup removes old entries periodically
func (rl *RateLimiter) cleanup() {
	for {
		select {
		case <-rl.cleanupTicker.C:
			rl.mu.Lock()
			now := time.Now()
			for key, timestamps := range rl.entries {
				// Remove timestamps older than the window
				var valid []time.Time
				for _, ts := range timestamps {
					if now.Sub(ts) < ChatRateLimitWindow {
						valid = append(valid, ts)
					}
				}
				if len(valid) == 0 {
					delete(rl.entries, key)
				} else {
					rl.entries[key] = valid
				}
			}
			rl.mu.Unlock()
		case <-rl.stopCleanup:
			return
		}
	}
}

// Stop stops the cleanup goroutine
func (rl *RateLimiter) Stop() {
	rl.cleanupTicker.Stop()
	close(rl.stopCleanup)
}

// CheckRateLimit checks if a user/IP can send a message
// Returns true if allowed, false if rate limited
func (rl *RateLimiter) CheckRateLimit(identifier string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	timestamps, exists := rl.entries[identifier]

	if !exists {
		// First message from this identifier
		rl.entries[identifier] = []time.Time{now}
		return true
	}

	// Remove timestamps outside the window
	var valid []time.Time
	for _, ts := range timestamps {
		if now.Sub(ts) < ChatRateLimitWindow {
			valid = append(valid, ts)
		}
	}

	// Check if we're at the limit
	if len(valid) >= ChatRateLimitMaxMessages {
		return false
	}

	// Add current timestamp
	valid = append(valid, now)
	rl.entries[identifier] = valid
	return true
}

// GetRemainingMessages returns the number of messages remaining in the current window
func (rl *RateLimiter) GetRemainingMessages(identifier string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	now := time.Now()
	timestamps, exists := rl.entries[identifier]

	if !exists {
		return ChatRateLimitMaxMessages
	}

	// Count valid timestamps
	count := 0
	for _, ts := range timestamps {
		if now.Sub(ts) < ChatRateLimitWindow {
			count++
		}
	}

	remaining := ChatRateLimitMaxMessages - count
	if remaining < 0 {
		return 0
	}
	return remaining
}

