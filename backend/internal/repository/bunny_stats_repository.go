package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cyclingstream/backend/internal/models"
)

type BunnyStatsRepository struct {
	db *sql.DB
}

func NewBunnyStatsRepository(db *sql.DB) *BunnyStatsRepository {
	return &BunnyStatsRepository{db: db}
}

func (r *BunnyStatsRepository) Upsert(ctx context.Context, stats *models.BunnyVideoStats) error {
	var geoRaw, rawPayload []byte
	var err error
	if stats.GeoBreakdown != nil {
		geoRaw, err = json.Marshal(stats.GeoBreakdown)
		if err != nil {
			return fmt.Errorf("marshal geo: %w", err)
		}
	}
	if stats.RawPayload != nil {
		rawPayload, err = json.Marshal(stats.RawPayload)
		if err != nil {
			return fmt.Errorf("marshal raw payload: %w", err)
		}
	}

	query := `
		INSERT INTO bunny_video_stats (
			bunny_video_id,
			stream_id,
			date,
			views,
			watch_time_seconds,
			geo_breakdown,
			raw_payload
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (bunny_video_id, date) DO UPDATE
		SET stream_id = EXCLUDED.stream_id,
			views = EXCLUDED.views,
			watch_time_seconds = EXCLUDED.watch_time_seconds,
			geo_breakdown = EXCLUDED.geo_breakdown,
			raw_payload = EXCLUDED.raw_payload,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		stats.BunnyVideoID,
		stats.StreamID,
		stats.Date,
		stats.Views,
		stats.WatchTimeSeconds,
		geoRaw,
		rawPayload,
	).Scan(&stats.ID, &stats.CreatedAt, &stats.UpdatedAt)
}
