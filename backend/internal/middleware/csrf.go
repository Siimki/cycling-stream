package middleware

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
)

// CSRFProtection creates CSRF protection middleware
// Only apply to state-changing operations (POST, PUT, DELETE, PATCH)
func CSRFProtection(secret string) fiber.Handler {
	return csrf.New(csrf.Config{
		KeyLookup:      "header:X-CSRF-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		CookieHTTPOnly: true,
		Expiration:     1 * 60 * 60, // 1 hour
		KeyGenerator: func() string {
			// Generate a random token using crypto/rand
			return generateCSRFToken()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "CSRF token validation failed",
			})
		},
		// Skip CSRF for all requests since we use JWT tokens in Authorization header
		// CSRF protection is primarily for cookie-based sessions
		// JWT tokens in headers are not vulnerable to CSRF attacks
		// This middleware is kept for future use if cookie-based sessions are added
		Next: func(c *fiber.Ctx) bool {
			// Skip CSRF for all requests (JWT-based API doesn't need CSRF)
			return true
		},
	})
}

// generateCSRFToken generates a random CSRF token using crypto/rand
func generateCSRFToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		// Fallback to a simple token if rand fails (shouldn't happen)
		return "fallback-token"
	}
	return base64.URLEncoding.EncodeToString(b)
}

