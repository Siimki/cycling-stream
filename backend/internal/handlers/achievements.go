package handlers

import (
	"github.com/cyclingstream/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

type AchievementsHandler struct {
	service *services.AchievementService
}

func NewAchievementsHandler(service *services.AchievementService) *AchievementsHandler {
	return &AchievementsHandler{service: service}
}

func (h *AchievementsHandler) GetUserAchievements(c *fiber.Ctx) error {
	if h.service == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(APIError{
			Error: "Achievement service unavailable",
		})
	}

	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	achievements, err := h.service.GetUserAchievements(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to fetch achievements",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"achievements": achievements,
	})
}
