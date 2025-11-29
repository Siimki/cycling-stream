package handlers

import (
	"strings"
	"time"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type AnalyticsIngestionHandler struct {
	streamRepo        *repository.StreamRepository
	playbackEventRepo *repository.PlaybackEventRepository
}

func NewAnalyticsIngestionHandler(
	streamRepo *repository.StreamRepository,
	playbackEventRepo *repository.PlaybackEventRepository,
) *AnalyticsIngestionHandler {
	return &AnalyticsIngestionHandler{
		streamRepo:        streamRepo,
		playbackEventRepo: playbackEventRepo,
	}
}

func (h *AnalyticsIngestionHandler) IngestEvents(c *fiber.Ctx) error {
	var req models.AnalyticsEventBatch
	if !parseBody(c, &req) {
		return nil
	}

	if req.StreamID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{Error: "streamId is required"})
	}
	if req.ClientID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{Error: "clientId is required"})
	}
	if len(req.ClientID) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{Error: "clientId too long"})
	}
	if len(req.Events) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{Error: "events array is required"})
	}
	if len(req.Events) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{Error: "events array too large"})
	}

	stream, err := h.streamRepo.GetByID(req.StreamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{Error: "failed to load stream"})
	}
	if stream == nil {
		return c.Status(fiber.StatusNotFound).JSON(APIError{Error: "stream not found"})
	}

	now := time.Now()
	deviceType := detectDeviceType(c.Get(fiber.HeaderUserAgent))
	country := detectCountry(c)

	allowed := map[string]bool{
		"play":         true,
		"pause":        true,
		"heartbeat":    true,
		"ended":        true,
		"error":        true,
		"buffer_start": true,
		"buffer_end":   true,
	}

	events := make([]models.PlaybackEvent, 0, len(req.Events))
	for _, evt := range req.Events {
		if !allowed[strings.ToLower(evt.Type)] {
			return c.Status(fiber.StatusBadRequest).JSON(APIError{Error: "invalid event type: " + evt.Type})
		}
		if evt.VideoTimeSeconds != nil && *evt.VideoTimeSeconds < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(APIError{Error: "videoTime must be >= 0"})
		}

		eventType := strings.ToLower(evt.Type)
		events = append(events, models.PlaybackEvent{
			StreamID:         req.StreamID,
			ClientID:         req.ClientID,
			EventType:        eventType,
			VideoTimeSeconds: evt.VideoTimeSeconds,
			Country:          country,
			DeviceType:       deviceType,
			Extra:            evt.Extra,
			CreatedAt:        now,
		})
	}

	if err := h.playbackEventRepo.InsertBatch(c.Context(), events); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{Error: "failed to store analytics events"})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"ingested": len(events),
	})
}

func detectCountry(c *fiber.Ctx) string {
	if country := c.Get("CF-IPCountry"); country != "" {
		return strings.ToLower(country)
	}
	if country := c.Get("X-Country-Code"); country != "" {
		return strings.ToLower(country)
	}
	if country := c.Get("X-Appengine-Country"); country != "" {
		return strings.ToLower(country)
	}
	return "unknown"
}

func detectDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)
	switch {
	case ua == "":
		return "unknown"
	case strings.Contains(ua, "ipad") || strings.Contains(ua, "tablet"):
		return "tablet"
	case strings.Contains(ua, "mobi") || strings.Contains(ua, "android"):
		return "mobile"
	case strings.Contains(ua, "smart-tv") || strings.Contains(ua, "hbbtv") || strings.Contains(ua, "tv"):
		return "tv"
	default:
		return "desktop"
	}
}
