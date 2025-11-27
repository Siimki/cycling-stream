package handlers

import (
	"strconv"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	raceRepo          *repository.RaceRepository
	viewerSessionRepo *repository.ViewerSessionRepository
	watchSessionRepo  *repository.WatchSessionRepository
	revenueRepo       *repository.RevenueRepository
}

func NewAnalyticsHandler(
	raceRepo *repository.RaceRepository,
	viewerSessionRepo *repository.ViewerSessionRepository,
	watchSessionRepo *repository.WatchSessionRepository,
	revenueRepo *repository.RevenueRepository,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		raceRepo:          raceRepo,
		viewerSessionRepo: viewerSessionRepo,
		watchSessionRepo:  watchSessionRepo,
		revenueRepo:       revenueRepo,
	}
}

// GetRaceAnalytics returns viewer analytics for all races
func (h *AnalyticsHandler) GetRaceAnalytics(c *fiber.Ctx) error {
	// Get all races
	races, err := h.raceRepo.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch races",
		})
	}

	// Get concurrent viewers for all races
	concurrentViewers, err := h.viewerSessionRepo.GetAllConcurrentViewers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch concurrent viewers",
		})
	}

	// Create a map for quick lookup
	concurrentMap := make(map[string]models.ConcurrentViewers)
	for _, cv := range concurrentViewers {
		concurrentMap[cv.RaceID] = cv
	}

	// Get unique viewers for all races
	uniqueViewers, err := h.viewerSessionRepo.GetAllUniqueViewers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch unique viewers",
		})
	}

	// Create a map for quick lookup
	uniqueMap := make(map[string]models.UniqueViewers)
	for _, uv := range uniqueViewers {
		uniqueMap[uv.RaceID] = uv
	}

	// Build response with race information
	type RaceAnalytics struct {
		RaceID                string `json:"race_id"`
		RaceName              string `json:"race_name"`
		ConcurrentViewers     int    `json:"concurrent_viewers"`
		AuthenticatedViewers  int    `json:"authenticated_viewers"`
		AnonymousViewers      int    `json:"anonymous_viewers"`
		UniqueViewers         int    `json:"unique_viewers"`
		UniqueAuthenticated   int    `json:"unique_authenticated"`
		UniqueAnonymous       int    `json:"unique_anonymous"`
	}

	var analytics []RaceAnalytics
	for _, race := range races {
		analyticsItem := RaceAnalytics{
			RaceID:   race.ID,
			RaceName: race.Name,
		}

		if cv, ok := concurrentMap[race.ID]; ok {
			analyticsItem.ConcurrentViewers = cv.ConcurrentCount
			analyticsItem.AuthenticatedViewers = cv.AuthenticatedCount
			analyticsItem.AnonymousViewers = cv.AnonymousCount
		}

		if uv, ok := uniqueMap[race.ID]; ok {
			analyticsItem.UniqueViewers = uv.UniqueViewerCount
			analyticsItem.UniqueAuthenticated = uv.UniqueAuthenticatedCount
			analyticsItem.UniqueAnonymous = uv.UniqueAnonymousCount
		}

		analytics = append(analytics, analyticsItem)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": analytics,
	})
}

// GetWatchTimeAnalytics returns watch time analytics aggregated by race
func (h *AnalyticsHandler) GetWatchTimeAnalytics(c *fiber.Ctx) error {
	// Parse optional filters
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

	// Get watch time analytics by race
	watchTimeData, err := h.watchSessionRepo.GetWatchTimeByRace(year, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch watch time analytics",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": watchTimeData,
	})
}

// GetRevenueAnalytics returns revenue analytics (wrapper around existing revenue endpoint)
func (h *AnalyticsHandler) GetRevenueAnalytics(c *fiber.Ctx) error {
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

