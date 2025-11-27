package handlers

import "github.com/gofiber/fiber/v2"

// APIError is a minimal, consistent error payload for all handlers.
type APIError struct {
	Error string `json:"error"`
}

// parseBody is a small helper to parse JSON request bodies and return a
// standardized 400 error response when parsing fails.
//
// Usage:
//   var req SomeRequest
//   if !parseBody(c, &req) {
//       return nil // response already sent
//   }
func parseBody(c *fiber.Ctx, dst interface{}) bool {
	if err := c.BodyParser(dst); err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(APIError{
			Error: "Invalid request body",
		})
		return false
	}
	return true
}

// requireParam ensures a required path parameter is present. If the parameter
// is missing, it sends a 400 response with the provided error message and
// returns an empty string and false.
//
// This keeps the exact error messages stable while removing duplication.
func requireParam(c *fiber.Ctx, name string, errorMessage string) (string, bool) {
	value := c.Params(name)
	if value == "" {
		_ = c.Status(fiber.StatusBadRequest).JSON(APIError{
			Error: errorMessage,
		})
		return "", false
	}
	return value, true
}

// requireUserID fetches the authenticated user_id from context. If it is
// missing or empty, it returns an empty string and sends a 401 response using
// the provided error message.
//
// This helper is intentionally strict (requires a non-empty string) and is
// meant for endpoints that require authentication. For optional authentication
// flows, handlers should continue to read c.Locals("user_id") directly.
func requireUserID(c *fiber.Ctx, errorMessage string) (string, bool) {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		_ = c.Status(fiber.StatusUnauthorized).JSON(APIError{
			Error: errorMessage,
		})
		return "", false
	}
	return userID, true
}


