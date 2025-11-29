package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cyclingstream/backend/internal/models"
)

type PlaybackEventRepository struct {
	db *sql.DB
}

func NewPlaybackEventRepository(db *sql.DB) *PlaybackEventRepository {
	return &PlaybackEventRepository{db: db}
}

// InsertBatch persists a batch of playback events in a single transaction.
func (r *PlaybackEventRepository) InsertBatch(ctx context.Context, events []models.PlaybackEvent) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO playback_events (
			stream_id,
			viewer_session_id,
			client_id,
			event_type,
			video_time_seconds,
			country,
			device_type,
			extra
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("prepare insert: %w", err)
	}
	defer stmt.Close()

	for _, event := range events {
		var extraBytes []byte
		if event.Extra != nil {
			extraBytes, err = json.Marshal(event.Extra)
			if err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("marshal extra: %w", err)
			}
		}

		if _, err := stmt.ExecContext(
			ctx,
			event.StreamID,
			event.ViewerSessionID,
			event.ClientID,
			event.EventType,
			event.VideoTimeSeconds,
			event.Country,
			event.DeviceType,
			extraBytes,
		); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("insert playback event: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit playback events: %w", err)
	}

	return nil
}

// ListByStreamSince returns playback events for a stream optionally after a given time.
func (r *PlaybackEventRepository) ListByStreamSince(ctx context.Context, streamID string, since *time.Time) ([]models.PlaybackEvent, error) {
	query := `
		SELECT id, stream_id, viewer_session_id, client_id, event_type, video_time_seconds, country, device_type, extra, created_at
		FROM playback_events
		WHERE stream_id = $1
	`
	args := []interface{}{streamID}
	if since != nil {
		query += " AND created_at >= $2"
		args = append(args, *since)
	}
	query += " ORDER BY client_id, created_at"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list playback events: %w", err)
	}
	defer rows.Close()

	var events []models.PlaybackEvent
	for rows.Next() {
		var evt models.PlaybackEvent
		var extraRaw []byte
		if err := rows.Scan(
			&evt.ID,
			&evt.StreamID,
			&evt.ViewerSessionID,
			&evt.ClientID,
			&evt.EventType,
			&evt.VideoTimeSeconds,
			&evt.Country,
			&evt.DeviceType,
			&extraRaw,
			&evt.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan playback event: %w", err)
		}
		if len(extraRaw) > 0 {
			if err := json.Unmarshal(extraRaw, &evt.Extra); err != nil {
				return nil, fmt.Errorf("unmarshal playback extra: %w", err)
			}
		}
		events = append(events, evt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate playback events: %w", err)
	}

	return events, nil
}

// DeleteOlderThan removes events older than cutoff. Use with care (run after aggregation).
func (r *PlaybackEventRepository) DeleteOlderThan(ctx context.Context, cutoff time.Time) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM playback_events
		WHERE created_at < $1
	`, cutoff)
	if err != nil {
		return 0, fmt.Errorf("delete old playback events: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected: %w", err)
	}
	return rows, nil
}
