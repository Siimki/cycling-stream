package middleware

import (
	"time"

	"github.com/cyclingstream/backend/internal/logger"
	"github.com/gofiber/fiber/v2"
)

// StructuredLogger creates a middleware for structured request logging
func StructuredLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Build log entry with structured fields
		entry := logger.WithFields(map[string]interface{}{
			"method":     c.Method(),
			"path":       c.Path(),
			"status":     c.Response().StatusCode(),
			"duration_ms": duration.Milliseconds(),
			"ip":         c.IP(),
			"user_agent": c.Get("User-Agent"),
		})

		// Add request ID if available
		if requestID := c.Get("X-Request-ID"); requestID != "" {
			entry = entry.WithField("request_id", requestID)
		}

		// Log based on status code
		status := c.Response().StatusCode()
		switch {
		case status >= 500:
			entry.Error("Request completed with server error")
		case status >= 400:
			entry.Warn("Request completed with client error")
		case duration > 1*time.Second:
			entry.Warn("Slow request detected")
		default:
			entry.Info("Request completed")
		}

		return err
	}
}

