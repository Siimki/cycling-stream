package handlers

import (
	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

// loadRaceOr404 fetches a race by ID and sends a standardized JSON error
// response when something goes wrong.
//
// It returns (race, true) on success, or (nil, false) if a response has
// already been sent to the client.
func loadRaceOr404(c *fiber.Ctx, raceRepo *repository.RaceRepository, id string) (*models.Race, bool) {
	race, err := raceRepo.GetByID(id)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to fetch race",
		})
		return nil, false
	}

	if race == nil {
		c.Status(fiber.StatusNotFound).JSON(APIError{
			Error: "Race not found",
		})
		return nil, false
	}

	return race, true
}

// loadStreamOr404 fetches a stream by race ID and sends a standardized JSON
// error response when the stream cannot be fetched or does not exist.
//
// It returns (stream, true) on success, or (nil, false) if a response has
// already been sent to the client.
func loadStreamOr404(c *fiber.Ctx, streamRepo *repository.StreamRepository, raceID string, notFoundMessage string) (*models.Stream, bool) {
	stream, err := streamRepo.GetByRaceID(raceID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to fetch stream",
		})
		return nil, false
	}

	if stream == nil {
		c.Status(fiber.StatusNotFound).JSON(APIError{
			Error: notFoundMessage,
		})
		return nil, false
	}

	return stream, true
}


