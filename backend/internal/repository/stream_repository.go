package repository

import (
	"database/sql"
	"fmt"

	"github.com/cyclingstream/backend/internal/models"
)

type StreamRepository struct {
	db *sql.DB
}

func NewStreamRepository(db *sql.DB) *StreamRepository {
	return &StreamRepository{db: db}
}

func (r *StreamRepository) GetByRaceID(raceID string) (*models.Stream, error) {
	query := `
		SELECT id, race_id, status, stream_type, source_id, origin_url, cdn_url, stream_key, created_at, updated_at
		FROM streams
		WHERE race_id = $1
		LIMIT 1
	`

	var stream models.Stream
	err := r.db.QueryRow(query, raceID).Scan(
		&stream.ID,
		&stream.RaceID,
		&stream.Status,
		&stream.StreamType,
		&stream.SourceID,
		&stream.OriginURL,
		&stream.CDNURL,
		&stream.StreamKey,
		&stream.CreatedAt,
		&stream.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get stream: %w", err)
	}

	return &stream, nil
}

func (r *StreamRepository) CreateOrUpdate(stream *models.Stream) error {
	query := `
		INSERT INTO streams (race_id, status, stream_type, source_id, origin_url, cdn_url, stream_key)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (race_id) DO UPDATE
		SET status = $2, stream_type = $3, source_id = $4, origin_url = $5, cdn_url = $6, stream_key = $7, updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at
	`

	if stream.StreamType == "" {
		stream.StreamType = "hls"
	}

	err := r.db.QueryRow(
		query,
		stream.RaceID,
		stream.Status,
		stream.StreamType,
		stream.SourceID,
		stream.OriginURL,
		stream.CDNURL,
		stream.StreamKey,
	).Scan(&stream.ID, &stream.CreatedAt, &stream.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create or update stream: %w", err)
	}

	return nil
}

func (r *StreamRepository) UpdateStatus(raceID string, status string) error {
	query := `
		UPDATE streams
		SET status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE race_id = $1
	`

	result, err := r.db.Exec(query, raceID, status)
	if err != nil {
		return fmt.Errorf("failed to update stream status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stream not found")
	}

	return nil
}
