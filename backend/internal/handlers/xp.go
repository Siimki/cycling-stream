package handlers

import (
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/cyclingstream/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

type XPHandler struct {
	userRepo  *repository.UserRepository
	xpService *services.XPService
}

func NewXPHandler(userRepo *repository.UserRepository, xpService *services.XPService) *XPHandler {
	return &XPHandler{
		userRepo:  userRepo,
		xpService: xpService,
	}
}

// GetUserXPProgress returns the current XP, level, and progress information for the authenticated user
func (h *XPHandler) GetUserXPProgress(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to get user XP progress",
		})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(APIError{
			Error: "User not found",
		})
	}

	// Get XP progress using the service method which uses correct formulas
	currentXPInLevel, xpNeededForNext := h.xpService.GetLevelProgress(user.XPTotal, user.Level)
	
	// XP threshold for current level start
	xpForCurrentLevelStart := h.xpService.GetXPForLevel(user.Level)
	
	// XP threshold for next level start
	xpForNextLevelStart := h.xpService.GetXPForNextLevel(user.Level)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":                  userID,
		"xp_total":                 user.XPTotal,
		"level":                    user.Level,
		"xp_for_current_level_start": xpForCurrentLevelStart,
		"xp_for_next_level":        xpForNextLevelStart,
		"xp_to_next_level":         xpNeededForNext,
		"progress_in_current_level": currentXPInLevel,
	})
}


