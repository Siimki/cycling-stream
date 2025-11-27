package handlers

import (
	"strconv"

	"github.com/cyclingstream/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type WatchHistoryHandler struct {
	historyRepo *repository.WatchHistoryRepository
}

func NewWatchHistoryHandler(historyRepo *repository.WatchHistoryRepository) *WatchHistoryHandler {
	return &WatchHistoryHandler{
		historyRepo: historyRepo,
	}
}

func (h *WatchHistoryHandler) GetWatchHistory(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	// Parse pagination parameters
	limit := 20 // default
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	entries, err := h.historyRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to fetch watch history",
		})
	}

	total, err := h.historyRepo.GetCountByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to fetch watch history count",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"entries": entries,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

