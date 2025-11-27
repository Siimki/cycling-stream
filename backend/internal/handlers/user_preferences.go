package handlers

import (
	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type UserPreferencesHandler struct {
	prefsRepo *repository.UserPreferencesRepository
}

func NewUserPreferencesHandler(prefsRepo *repository.UserPreferencesRepository) *UserPreferencesHandler {
	return &UserPreferencesHandler{
		prefsRepo: prefsRepo,
	}
}

func (h *UserPreferencesHandler) GetPreferences(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	prefs, err := h.prefsRepo.GetByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to fetch preferences",
		})
	}

	return c.Status(fiber.StatusOK).JSON(prefs)
}

func (h *UserPreferencesHandler) UpdatePreferences(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	var req models.UpdatePreferencesRequest
	if !parseBody(c, &req) {
		return nil
	}

	// Validate enum values if provided
	if req.DataMode != nil {
		if *req.DataMode != "casual" && *req.DataMode != "standard" && *req.DataMode != "pro" {
			return c.Status(fiber.StatusBadRequest).JSON(APIError{
				Error: "Invalid data_mode. Must be 'casual', 'standard', or 'pro'",
			})
		}
	}

	if req.PreferredUnits != nil {
		if *req.PreferredUnits != "metric" && *req.PreferredUnits != "imperial" {
			return c.Status(fiber.StatusBadRequest).JSON(APIError{
				Error: "Invalid preferred_units. Must be 'metric' or 'imperial'",
			})
		}
	}

	if req.Theme != nil {
		if *req.Theme != "light" && *req.Theme != "dark" && *req.Theme != "auto" {
			return c.Status(fiber.StatusBadRequest).JSON(APIError{
				Error: "Invalid theme. Must be 'light', 'dark', or 'auto'",
			})
		}
	}

	if req.DeviceType != nil {
		validDeviceTypes := map[string]bool{
			"tv":      true,
			"desktop": true,
			"mobile":  true,
			"tablet":  true,
		}
		if !validDeviceTypes[*req.DeviceType] {
			return c.Status(fiber.StatusBadRequest).JSON(APIError{
				Error: "Invalid device_type. Must be 'tv', 'desktop', 'mobile', or 'tablet'",
			})
		}
	}

	prefs, err := h.prefsRepo.Update(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to update preferences",
		})
	}

	return c.Status(fiber.StatusOK).JSON(prefs)
}

func (h *UserPreferencesHandler) CompleteOnboarding(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	req := models.UpdatePreferencesRequest{
		OnboardingCompleted: func() *bool { b := true; return &b }(),
	}

	prefs, err := h.prefsRepo.Update(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to mark onboarding as complete",
		})
	}

	return c.Status(fiber.StatusOK).JSON(prefs)
}

