package handlers

import (
	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

type PredictionsHandler struct {
	predictionService *services.PredictionService
}

func NewPredictionsHandler(predictionService *services.PredictionService) *PredictionsHandler {
	return &PredictionsHandler{
		predictionService: predictionService,
	}
}

// GetRacePredictions returns all prediction markets for a race
func (h *PredictionsHandler) GetRacePredictions(c *fiber.Ctx) error {
	raceID, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}

	// We need access to the repository to get markets
	// For now, we'll add a method to the service
	markets, err := h.predictionService.GetMarketsByRace(raceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to get prediction markets",
		})
	}

	return c.Status(fiber.StatusOK).JSON(markets)
}

// PlaceBet places a bet on a prediction market
func (h *PredictionsHandler) PlaceBet(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	_, ok = requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}

	marketID, ok := requireParam(c, "marketId", "Market ID is required")
	if !ok {
		return nil
	}

	var req models.PlaceBetRequest
	if !parseBody(c, &req) {
		return nil
	}

	if req.OptionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{
			Error: "Option ID is required",
		})
	}

	if req.StakePoints <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{
			Error: "Stake must be greater than 0",
		})
	}

	if err := h.predictionService.PlaceBet(userID, marketID, req.OptionID, req.StakePoints); err != nil {
		if err.Error() == "insufficient points balance" {
			return c.Status(fiber.StatusBadRequest).JSON(APIError{
				Error: "Insufficient points balance",
			})
		}
		if err.Error() == "stake exceeds maximum allowed" {
			return c.Status(fiber.StatusBadRequest).JSON(APIError{
				Error: err.Error(),
			})
		}
		if err.Error() == "market not found" || err.Error() == "option not found in market" {
			return c.Status(fiber.StatusNotFound).JSON(APIError{
				Error: err.Error(),
			})
		}
		if err.Error() == "market is not open for betting" {
			return c.Status(fiber.StatusBadRequest).JSON(APIError{
				Error: "Market is not open for betting",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to place bet",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Bet placed successfully",
	})
}

// GetUserPredictions returns all bets for the authenticated user
func (h *PredictionsHandler) GetUserPredictions(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	bets, err := h.predictionService.GetBetsByUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to get user predictions",
		})
	}

	return c.Status(fiber.StatusOK).JSON(bets)
}

