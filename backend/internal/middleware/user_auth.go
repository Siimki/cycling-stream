package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func UserAuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		tokenString := parts[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("is_admin", claims.IsAdmin)

		return c.Next()
	}
}

// OptionalUserAuthMiddleware extracts user info from JWT if present, but doesn't fail if missing
// This allows endpoints to work for both authenticated and anonymous users
func OptionalUserAuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// No auth header, continue as anonymous user
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Invalid format, continue as anonymous user
			return c.Next()
		}

		tokenString := parts[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			// Invalid token, continue as anonymous user
			return c.Next()
		}

		// Valid token, set user info
		c.Locals("user_id", claims.UserID)
		c.Locals("is_admin", claims.IsAdmin)

		return c.Next()
	}
}
