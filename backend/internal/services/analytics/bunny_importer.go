package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
)

type BunnyImporter struct {
	client          *BunnyClient
	providerRepo    *repository.StreamProviderRepository
	bunnyStatsRepo  *repository.BunnyStatsRepository
	streamStatsRepo *repository.StreamStatsRepository
}

func NewBunnyImporter(
	client *BunnyClient,
	providerRepo *repository.StreamProviderRepository,
	bunnyStatsRepo *repository.BunnyStatsRepository,
	streamStatsRepo *repository.StreamStatsRepository,
) *BunnyImporter {
	return &BunnyImporter{
		client:          client,
		providerRepo:    providerRepo,
		bunnyStatsRepo:  bunnyStatsRepo,
		streamStatsRepo: streamStatsRepo,
	}
}

// Sync fetches analytics for all bunny_stream providers and stores daily stats.
func (i *BunnyImporter) Sync(ctx context.Context) (int, error) {
	providers, err := i.providerRepo.ListByProvider("bunny_stream")
	if err != nil {
		return 0, fmt.Errorf("list bunny providers: %w", err)
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	count := 0
	for _, sp := range providers {
		stats, err := i.client.FetchVideoAnalytics(ctx, sp.ProviderVideoID)
		if err != nil {
			// Skip failures but continue processing
			continue
		}

		row := &models.BunnyVideoStats{
			BunnyVideoID:     sp.ProviderVideoID,
			StreamID:         &sp.StreamID,
			Date:             today,
			Views:            stats.Views,
			WatchTimeSeconds: stats.WatchTimeSeconds,
			GeoBreakdown:     stats.Geo,
			RawPayload:       stats.Raw,
		}

		if err := i.bunnyStatsRepo.Upsert(ctx, row); err != nil {
			continue
		}

		// Best-effort: update stream_stats with Bunny watch time if we have an existing row.
		existing, err := i.streamStatsRepo.GetByStreamID(ctx, sp.StreamID)
		if err == nil && existing != nil {
			existing.TotalWatchSeconds = maxInt64(existing.TotalWatchSeconds, stats.WatchTimeSeconds)
			if err := i.streamStatsRepo.Upsert(ctx, existing); err != nil {
				continue
			}
		}

		count++
	}

	return count, nil
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
