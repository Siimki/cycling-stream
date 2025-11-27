package handlers

import (
	"strconv"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type CostHandler struct {
	costRepo *repository.CostRepository
	raceRepo *repository.RaceRepository
}

func NewCostHandler(costRepo *repository.CostRepository, raceRepo *repository.RaceRepository) *CostHandler {
	return &CostHandler{
		costRepo: costRepo,
		raceRepo: raceRepo,
	}
}

// CreateCost creates a new cost entry
// POST /admin/costs
func (h *CostHandler) CreateCost(c *fiber.Ctx) error {
	var req models.CreateCostRequest
	if !parseBody(c, &req) {
		return nil
	}

	// Validate cost type
	if req.CostType != models.CostTypeCDN &&
		req.CostType != models.CostTypeServer &&
		req.CostType != models.CostTypeStorage &&
		req.CostType != models.CostTypeBandwidth &&
		req.CostType != models.CostTypeOther {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid cost_type. Must be one of: cdn, server, storage, bandwidth, other",
		})
	}

	// Validate month
	if req.Month < 1 || req.Month > 12 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Month must be between 1 and 12",
		})
	}

	// Validate year
	if req.Year < 2000 || req.Year > 2100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Year must be between 2000 and 2100",
		})
	}

	// Validate amount
	if req.AmountCents < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount must be non-negative",
		})
	}

	// If race_id is provided, verify it exists
	if req.RaceID != nil && *req.RaceID != "" {
		race, err := h.raceRepo.GetByID(*req.RaceID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to verify race",
			})
		}
		if race == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Race not found",
			})
		}
	}

	cost := &models.Cost{
		RaceID:      req.RaceID,
		CostType:    req.CostType,
		AmountCents: req.AmountCents,
		Year:        req.Year,
		Month:       req.Month,
		Description: req.Description,
	}

	if err := h.costRepo.Create(cost); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create cost",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(cost)
}

// GetCosts gets all costs with optional filters
// GET /admin/costs?year=2024&month=1
func (h *CostHandler) GetCosts(c *fiber.Ctx) error {
	var year, month *int

	if yearStr := c.Query("year"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = &y
		}
	}

	if monthStr := c.Query("month"); monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil {
			month = &m
		}
	}

	costs, err := h.costRepo.GetAll(year, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get costs",
		})
	}

	return c.JSON(costs)
}

// GetCostByID gets a specific cost by ID
// GET /admin/costs/:id
func (h *CostHandler) GetCostByID(c *fiber.Ctx) error {
	id, ok := requireParam(c, "id", "Cost ID is required")
	if !ok {
		return nil
	}

	cost, err := h.costRepo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get cost",
		})
	}

	if cost == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cost not found",
		})
	}

	return c.JSON(cost)
}

// GetCostsByRace gets all costs for a specific race
// GET /admin/costs/races/:race_id?year=2024&month=1
func (h *CostHandler) GetCostsByRace(c *fiber.Ctx) error {
	raceID, ok := requireParam(c, "race_id", "Race ID is required")
	if !ok {
		return nil
	}

	var year, month *int

	if yearStr := c.Query("year"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = &y
		}
	}

	if monthStr := c.Query("month"); monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil {
			month = &m
		}
	}

	costs, err := h.costRepo.GetByRace(raceID, year, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get costs",
		})
	}

	return c.JSON(costs)
}

// GetCostSummary gets monthly cost summary
// GET /admin/costs/summary?race_id=xxx&year=2024&month=1
func (h *CostHandler) GetCostSummary(c *fiber.Ctx) error {
	var raceID *string
	var year, month *int

	if raceIDStr := c.Query("race_id"); raceIDStr != "" {
		raceID = &raceIDStr
	}

	if yearStr := c.Query("year"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = &y
		}
	}

	if monthStr := c.Query("month"); monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil {
			month = &m
		}
	}

	summary, err := h.costRepo.GetMonthlySummary(raceID, year, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get cost summary",
		})
	}

	return c.JSON(summary)
}

// UpdateCost updates an existing cost
// PUT /admin/costs/:id
func (h *CostHandler) UpdateCost(c *fiber.Ctx) error {
	id, ok := requireParam(c, "id", "Cost ID is required")
	if !ok {
		return nil
	}

	var req models.CreateCostRequest
	if !parseBody(c, &req) {
		return nil
	}

	// Validate cost type
	if req.CostType != models.CostTypeCDN &&
		req.CostType != models.CostTypeServer &&
		req.CostType != models.CostTypeStorage &&
		req.CostType != models.CostTypeBandwidth &&
		req.CostType != models.CostTypeOther {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid cost_type. Must be one of: cdn, server, storage, bandwidth, other",
		})
	}

	// Validate month
	if req.Month < 1 || req.Month > 12 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Month must be between 1 and 12",
		})
	}

	// Validate year
	if req.Year < 2000 || req.Year > 2100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Year must be between 2000 and 2100",
		})
	}

	// Validate amount
	if req.AmountCents < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount must be non-negative",
		})
	}

	cost := &models.Cost{
		ID:          id,
		RaceID:      req.RaceID,
		CostType:    req.CostType,
		AmountCents: req.AmountCents,
		Year:        req.Year,
		Month:       req.Month,
		Description: req.Description,
	}

	if err := h.costRepo.Update(cost); err != nil {
		if err.Error() == "cost not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Cost not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update cost",
		})
	}

	return c.JSON(cost)
}

// DeleteCost deletes a cost
// DELETE /admin/costs/:id
func (h *CostHandler) DeleteCost(c *fiber.Ctx) error {
	id, ok := requireParam(c, "id", "Cost ID is required")
	if !ok {
		return nil
	}

	if err := h.costRepo.Delete(id); err != nil {
		if err.Error() == "cost not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Cost not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete cost",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

