package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cyclingstream/backend/internal/config"
	"github.com/cyclingstream/backend/internal/models"
)

type BunnyClient struct {
	httpClient *http.Client
	cfg        *config.BunnyConfig
}

func NewBunnyClient(cfg *config.BunnyConfig) *BunnyClient {
	return &BunnyClient{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		cfg:        cfg,
	}
}

// FetchVideoAnalytics pulls daily summary for a Bunny video ID.
func (c *BunnyClient) FetchVideoAnalytics(ctx context.Context, videoID string) (*models.BunnyAnalyticsResponse, error) {
	if c.cfg == nil || c.cfg.APIKey == "" || c.cfg.LibraryID == "" {
		return nil, fmt.Errorf("bunny config missing")
	}

	url := fmt.Sprintf("%s/%s/videos/%s/analytics", c.cfg.BaseURL, c.cfg.LibraryID, videoID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build bunny request: %w", err)
	}
	req.Header.Set("AccessKey", c.cfg.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call bunny analytics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bunny analytics status %d: %s", resp.StatusCode, string(body))
	}

	var payload map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode bunny response: %w", err)
	}

	result := &models.BunnyAnalyticsResponse{
		Raw: payload,
	}

	// Best-effort extraction from common Bunny fields.
	if views, ok := payload["views"].(float64); ok {
		result.Views = int(views)
	}
	if wt, ok := payload["watchTime"].(float64); ok {
		result.WatchTimeSeconds = int64(wt)
	}
	if geoRaw, ok := payload["geo"].(map[string]interface{}); ok {
		geo := make(map[string]int)
		for k, v := range geoRaw {
			if f, ok := v.(float64); ok {
				geo[k] = int(f)
			}
		}
		result.Geo = geo
	}

	return result, nil
}
