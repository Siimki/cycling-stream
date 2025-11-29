package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cyclingstream/backend/internal/models"
)

type StreamStatsRepository struct {
	db *sql.DB
}

func NewStreamStatsRepository(db *sql.DB) *StreamStatsRepository {
	return &StreamStatsRepository{db: db}
}

func (r *StreamStatsRepository) Upsert(ctx context.Context, stats *models.StreamStats) error {
	topCountriesRaw, err := json.Marshal(stats.TopCountries)
	if err != nil {
		return fmt.Errorf("marshal top countries: %w", err)
	}
	deviceBreakdownRaw, err := json.Marshal(stats.DeviceBreakdown)
	if err != nil {
		return fmt.Errorf("marshal device breakdown: %w", err)
	}

	query := `
		INSERT INTO stream_stats (
			stream_id,
			unique_viewers,
			total_watch_seconds,
			avg_watch_seconds,
			peak_concurrent_viewers,
			top_countries,
			device_breakdown,
			buffer_seconds,
			buffer_ratio,
			error_rate,
			last_calculated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (stream_id) DO UPDATE
		SET unique_viewers = EXCLUDED.unique_viewers,
			total_watch_seconds = EXCLUDED.total_watch_seconds,
			avg_watch_seconds = EXCLUDED.avg_watch_seconds,
			peak_concurrent_viewers = EXCLUDED.peak_concurrent_viewers,
			top_countries = EXCLUDED.top_countries,
			device_breakdown = EXCLUDED.device_breakdown,
			buffer_seconds = EXCLUDED.buffer_seconds,
			buffer_ratio = EXCLUDED.buffer_ratio,
			error_rate = EXCLUDED.error_rate,
			last_calculated_at = EXCLUDED.last_calculated_at,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		stats.StreamID,
		stats.UniqueViewers,
		stats.TotalWatchSeconds,
		stats.AvgWatchSeconds,
		stats.PeakConcurrentViewers,
		topCountriesRaw,
		deviceBreakdownRaw,
		stats.BufferSeconds,
		stats.BufferRatio,
		stats.ErrorRate,
		stats.LastCalculatedAt,
	)
	if err != nil {
		return fmt.Errorf("upsert stream stats: %w", err)
	}

	return nil
}

func (r *StreamStatsRepository) GetByStreamID(ctx context.Context, streamID string) (*models.StreamStats, error) {
	query := `
		SELECT stream_id, unique_viewers, total_watch_seconds, avg_watch_seconds, peak_concurrent_viewers, top_countries, device_breakdown, buffer_seconds, buffer_ratio, error_rate, last_calculated_at, created_at, updated_at
		FROM stream_stats
		WHERE stream_id = $1
	`

	var stats models.StreamStats
	var topCountriesRaw, deviceBreakdownRaw []byte
	err := r.db.QueryRowContext(ctx, query, streamID).Scan(
		&stats.StreamID,
		&stats.UniqueViewers,
		&stats.TotalWatchSeconds,
		&stats.AvgWatchSeconds,
		&stats.PeakConcurrentViewers,
		&topCountriesRaw,
		&deviceBreakdownRaw,
		&stats.BufferSeconds,
		&stats.BufferRatio,
		&stats.ErrorRate,
		&stats.LastCalculatedAt,
		&stats.CreatedAt,
		&stats.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get stream stats: %w", err)
	}

	if len(topCountriesRaw) > 0 {
		if err := json.Unmarshal(topCountriesRaw, &stats.TopCountries); err != nil {
			return nil, fmt.Errorf("unmarshal top countries: %w", err)
		}
	}
	if len(deviceBreakdownRaw) > 0 {
		if err := json.Unmarshal(deviceBreakdownRaw, &stats.DeviceBreakdown); err != nil {
			return nil, fmt.Errorf("unmarshal device breakdown: %w", err)
		}
	}

	return &stats, nil
}

// Summary returns aggregated stats across streams for admin summary boxes.
func (r *StreamStatsRepository) Summary(ctx context.Context) (models.StreamStatsSummary, error) {
	query := `
		SELECT 
			COUNT(*) as stream_count,
			COALESCE(SUM(unique_viewers), 0) as total_unique_viewers,
			COALESCE(SUM(total_watch_seconds), 0) as total_watch_seconds,
			COALESCE(AVG(peak_concurrent_viewers), 0) as avg_peak_concurrent
		FROM stream_stats
	`

	var summary models.StreamStatsSummary
	if err := r.db.QueryRowContext(ctx, query).Scan(
		&summary.StreamCount,
		&summary.TotalUniqueViewers,
		&summary.TotalWatchSeconds,
		&summary.AvgPeakConcurrent,
	); err != nil {
		return summary, fmt.Errorf("stream stats summary: %w", err)
	}

	return summary, nil
}

// StaleStreams returns stream IDs whose stats are older than the provided cutoff.
func (r *StreamStatsRepository) StaleStreams(ctx context.Context, cutoff time.Time) ([]string, error) {
	query := `
		SELECT stream_id
		FROM stream_stats
		WHERE last_calculated_at < $1
	`
	rows, err := r.db.QueryContext(ctx, query, cutoff)
	if err != nil {
		return nil, fmt.Errorf("stale streams: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan stale stream id: %w", err)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stale streams: %w", err)
	}

	return ids, nil
}
