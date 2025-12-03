package middleware

import (
	"strings"
	"time"

	"github.com/cyclingstream/backend/internal/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
)

// ChatAuthMiddleware extracts JWT token from WebSocket connection for optional authentication
// Similar to OptionalUserAuthMiddleware but works with WebSocket upgrade
func ChatAuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if this is a WebSocket upgrade request
		isWebSocketUpgrade := websocket.IsWebSocketUpgrade(c)
		logger.Info("ChatAuthMiddleware: Request received", map[string]interface{}{
			"path":               c.Path(),
			"method":             c.Method(),
			"isWebSocketUpgrade": isWebSocketUpgrade,
		})

		if !isWebSocketUpgrade {
			logger.Info("ChatAuthMiddleware: Not a WebSocket upgrade, passing through")
			return c.Next()
		}

		logger.Info("ChatAuthMiddleware: WebSocket upgrade detected")

		// Try to get token from query parameter first (common for WebSocket)
		tokenString := c.Query("token")

		// If not in query, try Authorization header
		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		if tokenString == "" {
			logger.Info("ChatAuthMiddleware: No token found, continuing as anonymous")
			return c.Next()
		}

		logger.Info("ChatAuthMiddleware: Token found, validating")

		// Parse and validate token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			logger.Info("ChatAuthMiddleware: Invalid or expired token", map[string]interface{}{
				"error": err,
			})
			return fiber.ErrUnauthorized
		}

		// Explicitly validate registered claims (expires/nbf) to be safe
		now := time.Now()
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(now) {
			logger.Info("ChatAuthMiddleware: Claims validation failed", map[string]interface{}{
				"error": "expired",
			})
			return fiber.ErrUnauthorized
		}

		if claims.NotBefore != nil && claims.NotBefore.Time.After(now) {
			logger.Info("ChatAuthMiddleware: Claims validation failed", map[string]interface{}{
				"error": "not before in future",
			})
			return fiber.ErrUnauthorized
		}

		logger.Info("ChatAuthMiddleware: Valid token, setting user info", map[string]interface{}{
			"user_id":  claims.UserID,
			"is_admin": claims.IsAdmin,
		})

		// Valid token, set user info in locals
		c.Locals("user_id", claims.UserID)
		c.Locals("is_admin", claims.IsAdmin)

		return c.Next()
	}
}
