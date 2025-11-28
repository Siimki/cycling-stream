package handlers

import (
	"github.com/cyclingstream/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

type MissionsHandler struct {
	missionService *services.MissionService
}

func NewMissionsHandler(missionService *services.MissionService) *MissionsHandler {
	return &MissionsHandler{
		missionService: missionService,
	}
}

// GetUserMissions returns all missions for the authenticated user
func (h *MissionsHandler) GetUserMissions(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	userMissions, err := h.missionService.GetActiveMissionsForUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to get user missions",
		})
	}

	return c.Status(fiber.StatusOK).JSON(userMissions)
}

// GetActiveMissions returns all active missions (public endpoint)
func (h *MissionsHandler) GetActiveMissions(c *fiber.Ctx) error {
	missions, err := h.missionService.GetActiveMissions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to get active missions",
		})
	}

	return c.Status(fiber.StatusOK).JSON(missions)
}

// ClaimMissionReward claims the reward for a completed mission
func (h *MissionsHandler) ClaimMissionReward(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	missionID, ok := requireParam(c, "missionId", "Mission ID is required")
	if !ok {
		return nil
	}

	if err := h.missionService.ClaimMissionReward(userID, missionID); err != nil {
		if err.Error() == "mission not completed" {
			return c.Status(fiber.StatusBadRequest).JSON(APIError{
				Error: "Mission not completed",
			})
		}
		if err.Error() == "mission reward already claimed" {
			return c.Status(fiber.StatusBadRequest).JSON(APIError{
				Error: "Mission reward already claimed",
			})
		}
		if err.Error() == "user mission not found" || err.Error() == "mission not found" {
			return c.Status(fiber.StatusNotFound).JSON(APIError{
				Error: "Mission not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to claim mission reward",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Mission reward claimed successfully",
	})
}

// GetCareerMissions returns career missions for the authenticated user
func (h *MissionsHandler) GetCareerMissions(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	userMissions, err := h.missionService.GetActiveCareerMissionsForUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to get career missions",
		})
	}

	return c.Status(fiber.StatusOK).JSON(userMissions)
}

// GetWeeklyMissions returns weekly missions for the authenticated user
func (h *MissionsHandler) GetWeeklyMissions(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	userMissions, err := h.missionService.GetWeeklyMissionsForUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to get weekly missions",
		})
	}

	return c.Status(fiber.StatusOK).JSON(userMissions)
}

