package handlers

import (
	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type RaceHandler struct {
	raceRepo        *repository.RaceRepository
	streamRepo      *repository.StreamRepository
	entitlementRepo *repository.EntitlementRepository
}

func NewRaceHandler(raceRepo *repository.RaceRepository, streamRepo *repository.StreamRepository, entitlementRepo *repository.EntitlementRepository) *RaceHandler {
	return &RaceHandler{
		raceRepo:        raceRepo,
		streamRepo:      streamRepo,
		entitlementRepo: entitlementRepo,
	}
}

func (h *RaceHandler) GetRaces(c *fiber.Ctx) error {
	races, err := h.raceRepo.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to fetch races",
		})
	}

	// Ensure we always return an array, not null
	// Use make() to create a non-nil empty slice that JSON encodes as []
	if races == nil || len(races) == 0 {
		races = make([]models.Race, 0)
	}

	return c.Status(fiber.StatusOK).JSON(races)
}

func (h *RaceHandler) GetRaceByID(c *fiber.Ctx) error {
	id, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}

	race, ok := loadRaceOr404(c, h.raceRepo, id)
	if !ok {
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(race)
}

func (h *RaceHandler) GetRaceStream(c *fiber.Ctx) error {
	id, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}

	// Verify race exists
	race, ok := loadRaceOr404(c, h.raceRepo, id)
	if !ok {
		return nil
	}

	// Check if race requires login
	if race.RequiresLogin {
		userID, ok := requireUserID(c, "Authentication required to access this stream")
		if !ok {
			return nil
		}
		_ = userID // User ID is validated, continue
	}

	// Check if user has access (if race is paid)
	if !race.IsFree {
		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":            "Authentication required",
				"requires_payment": true,
			})
		}

		hasAccess, err := h.entitlementRepo.HasAccess(userID, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check access",
			})
		}

		if !hasAccess {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":            "Payment required to access this race",
				"requires_payment": true,
			})
		}
	}

	// Get stream for race
	stream, ok := loadStreamOr404(c, h.streamRepo, id, "Stream not found for this race")
	if !ok {
		return nil
	}

	// Return stream info (prefer CDN URL if available)
	response := fiber.Map{
		"stream_id":   stream.ID,
		"status":      stream.Status,
		"stream_type": stream.StreamType,
		"provider":    stream.StreamType,
		"source_id":   stream.SourceID,
	}

	if stream.CDNURL != nil && *stream.CDNURL != "" {
		response["cdn_url"] = *stream.CDNURL
	} else if stream.OriginURL != nil && *stream.OriginURL != "" {
		response["origin_url"] = *stream.OriginURL
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
