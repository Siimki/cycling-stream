package handlers

import (
	"strconv"
	"time"

	"github.com/cyclingstream/backend/internal/middleware"
	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	raceRepo    *repository.RaceRepository
	streamRepo  *repository.StreamRepository
	revenueRepo *repository.RevenueRepository
}

func NewAdminHandler(raceRepo *repository.RaceRepository, streamRepo *repository.StreamRepository, revenueRepo *repository.RevenueRepository) *AdminHandler {
	return &AdminHandler{
		raceRepo:    raceRepo,
		streamRepo:  streamRepo,
		revenueRepo: revenueRepo,
	}
}

type CreateRaceRequest struct {
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Location    *string    `json:"location"`
	Category    *string    `json:"category"`
	IsFree      bool       `json:"is_free"`
	PriceCents  int        `json:"price_cents"`
}

func (h *AdminHandler) CreateRace(c *fiber.Ctx) error {
	var req CreateRaceRequest
	if !parseBody(c, &req) {
		return nil
	}

	// Validate and sanitize name
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
		})
	}
	req.Name = middleware.SanitizeString(req.Name, 255)

	// Sanitize optional string fields
	if req.Description != nil {
		sanitized := middleware.SanitizeString(*req.Description, 1000)
		req.Description = &sanitized
	}
	if req.Location != nil {
		sanitized := middleware.SanitizeString(*req.Location, 255)
		req.Location = &sanitized
	}
	if req.Category != nil {
		sanitized := middleware.SanitizeString(*req.Category, 100)
		req.Category = &sanitized
	}

	// Validate price
	if req.PriceCents < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Price cannot be negative",
		})
	}
	if !req.IsFree && req.PriceCents == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Price must be greater than 0 for paid races",
		})
	}

	race := &models.Race{
		Name:        req.Name,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Location:    req.Location,
		Category:    req.Category,
		IsFree:      req.IsFree,
		PriceCents:  req.PriceCents,
	}

	if err := h.raceRepo.Create(race); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create race",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(race)
}

func (h *AdminHandler) UpdateRace(c *fiber.Ctx) error {
	id, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}
	if !middleware.ValidateUUID(id) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid race ID format",
		})
	}

	var req CreateRaceRequest
	if !parseBody(c, &req) {
		return nil
	}

	// Validate and sanitize name if provided
	if req.Name != "" {
		req.Name = middleware.SanitizeString(req.Name, 255)
	}

	// Sanitize optional string fields
	if req.Description != nil {
		sanitized := middleware.SanitizeString(*req.Description, 1000)
		req.Description = &sanitized
	}
	if req.Location != nil {
		sanitized := middleware.SanitizeString(*req.Location, 255)
		req.Location = &sanitized
	}
	if req.Category != nil {
		sanitized := middleware.SanitizeString(*req.Category, 100)
		req.Category = &sanitized
	}

	// Validate price
	if req.PriceCents < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Price cannot be negative",
		})
	}

	race := &models.Race{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Location:    req.Location,
		Category:    req.Category,
		IsFree:      req.IsFree,
		PriceCents:  req.PriceCents,
	}

	if err := h.raceRepo.Update(race); err != nil {
		if err.Error() == "race not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Race not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update race",
		})
	}

	return c.Status(fiber.StatusOK).JSON(race)
}

func (h *AdminHandler) DeleteRace(c *fiber.Ctx) error {
	id, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}
	if !middleware.ValidateUUID(id) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid race ID format",
		})
	}

	if err := h.raceRepo.Delete(id); err != nil {
		if err.Error() == "race not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Race not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete race",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Race deleted successfully",
	})
}

type UpdateStreamRequest struct {
	OriginURL  *string `json:"origin_url"`
	CDNURL     *string `json:"cdn_url"`
	StreamKey  *string `json:"stream_key"`
	Status     string  `json:"status"`
	StreamType string  `json:"stream_type"`
	SourceID   *string `json:"source_id"`
}

func (h *AdminHandler) UpdateStream(c *fiber.Ctx) error {
	raceID, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}
	if !middleware.ValidateUUID(raceID) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid race ID format",
		})
	}

	var req UpdateStreamRequest
	if !parseBody(c, &req) {
		return nil
	}

	// Validate status
	validStatuses := map[string]bool{
		"offline": true,
		"live":    true,
		"ended":   true,
	}
	if !validStatuses[req.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status. Must be one of: offline, live, ended",
		})
	}

	// Validate stream type
	if req.StreamType == "" {
		req.StreamType = "hls"
	}
	validTypes := map[string]bool{
		"hls":     true,
		"youtube": true,
	}
	if !validTypes[req.StreamType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid stream type. Must be one of: hls, youtube",
		})
	}

	// Validate source_id for youtube
	if req.StreamType == "youtube" && (req.SourceID == nil || *req.SourceID == "") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Source ID is required for YouTube streams",
		})
	}

	// Sanitize URLs if provided
	if req.OriginURL != nil {
		sanitized := middleware.SanitizeString(*req.OriginURL, 500)
		req.OriginURL = &sanitized
	}
	if req.CDNURL != nil {
		sanitized := middleware.SanitizeString(*req.CDNURL, 500)
		req.CDNURL = &sanitized
	}
	if req.StreamKey != nil {
		sanitized := middleware.SanitizeString(*req.StreamKey, 255)
		req.StreamKey = &sanitized
	}
	if req.SourceID != nil {
		sanitized := middleware.SanitizeString(*req.SourceID, 255)
		req.SourceID = &sanitized
	}

	stream := &models.Stream{
		RaceID:     raceID,
		Status:     req.Status,
		StreamType: req.StreamType,
		SourceID:   req.SourceID,
		OriginURL:  req.OriginURL,
		CDNURL:     req.CDNURL,
		StreamKey:  req.StreamKey,
	}

	if err := h.streamRepo.CreateOrUpdate(stream); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update stream",
		})
	}

	return c.Status(fiber.StatusOK).JSON(stream)
}

func (h *AdminHandler) UpdateStreamStatus(c *fiber.Ctx) error {
	raceID, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}
	if !middleware.ValidateUUID(raceID) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid race ID format",
		})
	}

	var req struct {
		Status string `json:"status"`
	}
	if !parseBody(c, &req) {
		return nil
	}

	// Validate status
	validStatuses := map[string]bool{
		"offline": true,
		"live":    true,
		"ended":   true,
	}
	if !validStatuses[req.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status. Must be one of: offline, live, ended",
		})
	}

	if err := h.streamRepo.UpdateStatus(raceID, req.Status); err != nil {
		if err.Error() == "stream not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Stream not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update stream status",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Stream status updated successfully",
	})
}

// GetRevenue gets all revenue data, optionally filtered by year and month
func (h *AdminHandler) GetRevenue(c *fiber.Ctx) error {
	var year, month *int

	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr != "" {
		y, err := strconv.Atoi(yearStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid year parameter",
			})
		}
		year = &y
	}

	if monthStr != "" {
		m, err := strconv.Atoi(monthStr)
		if err != nil || m < 1 || m > 12 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid month parameter (must be 1-12)",
			})
		}
		month = &m
	}

	revenues, err := h.revenueRepo.GetAllMonthlyRevenue(year, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get revenue data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": revenues,
	})
}

// GetRevenueByRace gets monthly revenue data for a specific race
func (h *AdminHandler) GetRevenueByRace(c *fiber.Ctx) error {
	raceID, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}
	if !middleware.ValidateUUID(raceID) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid race ID format",
		})
	}

	revenues, err := h.revenueRepo.GetMonthlyRevenueByRace(raceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get revenue data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": revenues,
	})
}

// GetRevenueSummaryByRace gets aggregated revenue summary for a specific race
func (h *AdminHandler) GetRevenueSummaryByRace(c *fiber.Ctx) error {
	raceID, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}
	if !middleware.ValidateUUID(raceID) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid race ID format",
		})
	}

	summary, err := h.revenueRepo.GetRevenueSummaryByRace(raceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get revenue summary",
		})
	}

	return c.Status(fiber.StatusOK).JSON(summary)
}

// RecalculateRevenue recalculates all monthly revenue data
func (h *AdminHandler) RecalculateRevenue(c *fiber.Ctx) error {
	err := h.revenueRepo.RecalculateAllMonthlyRevenue()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to recalculate revenue",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Revenue data recalculated successfully",
	})
}

// RecalculateRevenueForPeriod recalculates revenue for a specific year and month
func (h *AdminHandler) RecalculateRevenueForPeriod(c *fiber.Ctx) error {
	yearStr := c.Params("year")
	monthStr := c.Params("month")

	if yearStr == "" || monthStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Year and month are required",
		})
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid year parameter",
		})
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid month parameter (must be 1-12)",
		})
	}

	err = h.revenueRepo.RecalculateMonthlyRevenueForPeriod(year, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to recalculate revenue for period",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Revenue data recalculated successfully",
	})
}

