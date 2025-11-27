package handlers

import (
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type StreamHandler struct {
	streamRepo *repository.StreamRepository
}

func NewStreamHandler(streamRepo *repository.StreamRepository) *StreamHandler {
	return &StreamHandler{
		streamRepo: streamRepo,
	}
}

func (h *StreamHandler) GetStreamStatus(c *fiber.Ctx) error {
	raceID, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}

	stream, ok := loadStreamOr404(c, h.streamRepo, raceID, "Stream not found")
	if !ok {
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": stream.Status,
	})
}

