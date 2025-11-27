package handlers

import (
	"time"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ViewerHandler struct {
	viewerSessionRepo *repository.ViewerSessionRepository
}

func NewViewerHandler(viewerSessionRepo *repository.ViewerSessionRepository) *ViewerHandler {
	return &ViewerHandler{
		viewerSessionRepo: viewerSessionRepo,
	}
}

// StartSession starts a viewer session for a race
// Supports both authenticated and anonymous users
func (h *ViewerHandler) StartSession(c *fiber.Ctx) error {
	var req models.StartViewerSessionRequest
	if !parseBody(c, &req) {
		return nil
	}

	if req.RaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Race ID is required",
		})
	}

	// Get user ID if authenticated (optional)
	var userID *string
	if userIDStr, ok := c.Locals("user_id").(string); ok && userIDStr != "" {
		userID = &userIDStr
	}

	// Get or generate session token
	sessionToken := c.Cookies("viewer_session_token")
	if sessionToken == "" {
		sessionToken = uuid.New().String()
		// Set cookie for anonymous users (expires in 24 hours)
		c.Cookie(&fiber.Cookie{
			Name:     "viewer_session_token",
			Value:    sessionToken,
			Expires:  time.Now().Add(24 * time.Hour),
			HTTPOnly: true,
			SameSite: "Lax",
		})
	}

	// Check if there's an active session for this race and token
	activeSession, err := h.viewerSessionRepo.GetActiveSessionByToken(req.RaceID, sessionToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check active session",
		})
	}

	if activeSession != nil {
		// Update heartbeat and return existing session
		if err := h.viewerSessionRepo.UpdateHeartbeat(activeSession.ID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update session",
			})
		}
		activeSession.LastSeenAt = time.Now()
		return c.Status(fiber.StatusOK).JSON(activeSession)
	}

	// Create new session
	now := time.Now()
	session := &models.ViewerSession{
		UserID:       userID,
		RaceID:       req.RaceID,
		SessionToken: sessionToken,
		StartedAt:    now,
		LastSeenAt:   now,
		IsActive:     true,
	}

	if err := h.viewerSessionRepo.Create(session); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create viewer session",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(session)
}

// EndSession ends a viewer session
func (h *ViewerHandler) EndSession(c *fiber.Ctx) error {
	var req models.EndViewerSessionRequest
	if !parseBody(c, &req) {
		return nil
	}

	if req.SessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Session ID is required",
		})
	}

	// Verify session belongs to user if authenticated
	if !verifyViewerSessionOwnership(c, h.viewerSessionRepo, req.SessionID) {
		return nil
	}

	if err := h.viewerSessionRepo.EndSession(req.SessionID); err != nil {
		if err.Error() == "session not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Session not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to end viewer session",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Viewer session ended successfully",
	})
}

// Heartbeat updates the last_seen_at timestamp for an active session
func (h *ViewerHandler) Heartbeat(c *fiber.Ctx) error {
	var req models.HeartbeatViewerSessionRequest
	if !parseBody(c, &req) {
		return nil
	}

	if req.SessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Session ID is required",
		})
	}

	// Verify session belongs to user if authenticated
	if !verifyViewerSessionOwnership(c, h.viewerSessionRepo, req.SessionID) {
		return nil
	}

	if err := h.viewerSessionRepo.UpdateHeartbeat(req.SessionID); err != nil {
		if err.Error() == "session not found or not active" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Session not found or not active",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update heartbeat",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Heartbeat updated successfully",
	})
}

// GetConcurrentViewers returns the current concurrent viewer count for a race
func (h *ViewerHandler) GetConcurrentViewers(c *fiber.Ctx) error {
	raceID, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}

	viewers, err := h.viewerSessionRepo.GetConcurrentViewers(raceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get concurrent viewers",
		})
	}

	return c.Status(fiber.StatusOK).JSON(viewers)
}

// GetUniqueViewers returns the total unique viewer count for a race
func (h *ViewerHandler) GetUniqueViewers(c *fiber.Ctx) error {
	raceID, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}

	viewers, err := h.viewerSessionRepo.GetUniqueViewers(raceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get unique viewers",
		})
	}

	return c.Status(fiber.StatusOK).JSON(viewers)
}

