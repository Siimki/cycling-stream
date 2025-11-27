package handlers

import (
	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userRepo         *repository.UserRepository
	watchSessionRepo *repository.WatchSessionRepository
}

func NewUserHandler(userRepo *repository.UserRepository, watchSessionRepo *repository.WatchSessionRepository) *UserHandler {
	return &UserHandler{
		userRepo:         userRepo,
		watchSessionRepo: watchSessionRepo,
	}
}

// GetPublicProfile returns publicly visible user information
func (h *UserHandler) GetPublicProfile(c *fiber.Ctx) error {
	userID, ok := requireParam(c, "id", "User ID is required")
	if !ok {
		return nil
	}

	// Get public user data (no email)
	user, err := h.userRepo.GetPublicByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get total watch minutes
	totalMinutes, err := h.watchSessionRepo.GetTotalWatchMinutesByUser(userID)
	if err != nil {
		// Log error but don't fail - return 0 as default
		totalMinutes = 0
	}

	// Create response with watch time included
	response := models.PublicUser{
		ID:                user.ID,
		Name:              user.Name,
		Bio:               user.Bio,
		Points:            user.Points,
		TotalWatchMinutes: totalMinutes,
		CreatedAt:         user.CreatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetLeaderboard returns all users with their points and watch time for the leaderboard
func (h *UserHandler) GetLeaderboard(c *fiber.Ctx) error {
	entries, err := h.userRepo.GetLeaderboard()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch leaderboard",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entries)
}

