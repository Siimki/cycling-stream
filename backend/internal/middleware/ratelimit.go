package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimiter creates a rate limiting middleware
// Configurable limits for different route groups
func RateLimiter(maxRequests int, window time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxRequests,
		Expiration: window,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as the key for rate limiting
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests, please try again later",
			})
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	})
}

// StrictRateLimiter for sensitive endpoints (auth, payments)
func StrictRateLimiter() fiber.Handler {
	return RateLimiter(1000, 10*time.Second) // 1000 requests per 10 seconds (increased significantly for testing)
}

// StandardRateLimiter for general API endpoints
func StandardRateLimiter() fiber.Handler {
	return RateLimiter(500, 10*time.Second) // 500 requests per 10 seconds (increased from 100)
}

// LenientRateLimiter for public read-only endpoints
func LenientRateLimiter() fiber.Handler {
	return RateLimiter(1000, 10*time.Second) // 1000 requests per 10 seconds (increased from 200)
}
