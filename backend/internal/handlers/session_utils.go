package handlers

import (
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

// verifyViewerSessionOwnership ensures that a viewer session belongs to the
// currently authenticated user (when a user is authenticated).
//
// If the user is not authenticated, this is a no-op and returns true.
// If the session exists and belongs to a different user, it sends a 403
// response and returns false.
// On other repository errors it sends a 500 response and returns false.
func verifyViewerSessionOwnership(c *fiber.Ctx, repo *repository.ViewerSessionRepository, sessionID string) bool {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		// Anonymous or unauthenticated: nothing to verify here.
		return true
	}

	session, err := repo.GetByID(sessionID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to get session",
		})
		return false
	}
	if session != nil && session.UserID != nil && *session.UserID != userID {
		c.Status(fiber.StatusForbidden).JSON(APIError{
			Error: "Unauthorized",
		})
		return false
	}

	return true
}


