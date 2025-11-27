package handlers

import (
	"time"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type WatchHandler struct {
	watchSessionRepo *repository.WatchSessionRepository
	userRepo         *repository.UserRepository
}

func NewWatchHandler(watchSessionRepo *repository.WatchSessionRepository, userRepo *repository.UserRepository) *WatchHandler {
	return &WatchHandler{
		watchSessionRepo: watchSessionRepo,
		userRepo:         userRepo,
	}
}

// Faster iteration: 10 points per 10 seconds (1 point/second).
// Points are now awarded from the frontend via the /users/me/points/bonus endpoint
// while watching, so this handler no longer manages point accrual directly.
const pointsPerBlockSeconds = 10
const pointsPerBlock = 10

func (h *WatchHandler) StartSession(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	var req models.StartWatchSessionRequest
	if !parseBody(c, &req) {
		return nil
	}

	if req.RaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Race ID is required",
		})
	}

	// Check if there's an active session
	activeSession, err := h.watchSessionRepo.GetActiveSession(userID, req.RaceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check active session",
		})
	}

	if activeSession != nil {
		// Return existing active session
		return c.Status(fiber.StatusOK).JSON(activeSession)
	}

	// Create new session
	session := &models.WatchSession{
		UserID:    userID,
		RaceID:    req.RaceID,
		StartedAt: time.Now(),
	}

	if err := h.watchSessionRepo.Create(session); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create watch session",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(session)
}

func (h *WatchHandler) EndSession(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	var req models.EndWatchSessionRequest
	if !parseBody(c, &req) {
		return nil
	}

	if req.SessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Session ID is required",
		})
	}

	if err := h.watchSessionRepo.EndSession(req.SessionID, userID); err != nil {
		if err.Error() == "session not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Session not found",
			})
		}
		if err.Error() == "unauthorized" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		if err.Error() == "session already ended" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Session already ended",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to end watch session",
		})
	}

	// Fetch the session again to get the final duration for analytics.
	// Points are now awarded incrementally via the tick endpoint while watching,
	// so we no longer grant additional points here to avoid double counting.
	session, err := h.watchSessionRepo.GetByID(req.SessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load watch session for points",
		})
	}

	if session == nil || session.DurationSeconds == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Watch session ended successfully",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":          "Watch session ended successfully",
		"awarded_points":   0,
		"duration_seconds": *session.DurationSeconds,
	})
}

func (h *WatchHandler) GetStats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	raceID := c.Params("race_id")
	if raceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Race ID is required",
		})
	}

	stats, err := h.watchSessionRepo.GetStatsByUserAndRace(userID, raceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get watch time stats",
		})
	}

	return c.Status(fiber.StatusOK).JSON(stats)
}

// Note: live point accrual is handled by AuthHandler.AwardBonusPoints.
