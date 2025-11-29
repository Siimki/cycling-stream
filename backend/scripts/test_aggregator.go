package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cyclingstream/backend/internal/config"
	"github.com/cyclingstream/backend/internal/database"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/cyclingstream/backend/internal/services/analytics"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.New(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	streamID := "6f447475-d86f-42de-9be7-bf304dbb4d78"
	if len(os.Args) > 1 {
		streamID = os.Args[1]
	}

	playbackRepo := repository.NewPlaybackEventRepository(db.DB)
	statsRepo := repository.NewStreamStatsRepository(db.DB)
	streamRepo := repository.NewStreamRepository(db.DB)

	aggregator := analytics.NewAggregator(playbackRepo, statsRepo, streamRepo)

	fmt.Printf("Aggregating stats for stream: %s\n", streamID)
	stats, err := aggregator.AggregateStream(context.Background(), streamID, nil)
	if err != nil {
		log.Fatalf("Failed to aggregate: %v", err)
	}

	fmt.Printf("\n=== Stream Stats ===\n")
	fmt.Printf("Stream ID: %s\n", stats.StreamID)
	fmt.Printf("Unique Viewers: %d\n", stats.UniqueViewers)
	fmt.Printf("Total Watch Seconds: %d\n", stats.TotalWatchSeconds)
	fmt.Printf("Avg Watch Seconds: %d\n", stats.AvgWatchSeconds)
	fmt.Printf("Peak Concurrent Viewers: %d\n", stats.PeakConcurrentViewers)
	fmt.Printf("Buffer Seconds: %d\n", stats.BufferSeconds)
	fmt.Printf("Buffer Ratio: %.4f\n", stats.BufferRatio)
	fmt.Printf("Error Rate: %.4f\n", stats.ErrorRate)
	fmt.Printf("Top Countries: %v\n", stats.TopCountries)
	fmt.Printf("Device Breakdown: %v\n", stats.DeviceBreakdown)
	fmt.Printf("Last Calculated At: %s\n", stats.LastCalculatedAt)
}

