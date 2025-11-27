package handlers

import (
	"github.com/cyclingstream/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

type RecommendationsHandler struct {
	recommendationService *services.RecommendationService
}

func NewRecommendationsHandler(recommendationService *services.RecommendationService) *RecommendationsHandler {
	return &RecommendationsHandler{
		recommendationService: recommendationService,
	}
}

func (h *RecommendationsHandler) GetAllRecommendations(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	recommendations, err := h.recommendationService.GetAllRecommendations(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to get recommendations",
		})
	}

	return c.Status(fiber.StatusOK).JSON(recommendations)
}

func (h *RecommendationsHandler) GetContinueWatching(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	races, err := h.recommendationService.GetContinueWatching(userID, 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to get continue watching",
		})
	}

	return c.Status(fiber.StatusOK).JSON(races)
}

func (h *RecommendationsHandler) GetUpcoming(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	races, err := h.recommendationService.GetUpcomingRacesForUser(userID, 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to get upcoming races",
		})
	}

	return c.Status(fiber.StatusOK).JSON(races)
}

func (h *RecommendationsHandler) GetReplays(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	races, err := h.recommendationService.GetRecommendedReplays(userID, 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to get recommended replays",
		})
	}

	return c.Status(fiber.StatusOK).JSON(races)
}

