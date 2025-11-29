package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cyclingstream/backend/internal/models"
)

type StreamProviderRepository struct {
	db *sql.DB
}

func NewStreamProviderRepository(db *sql.DB) *StreamProviderRepository {
	return &StreamProviderRepository{db: db}
}

// GetPrimaryByStreamID returns the latest provider row for a stream.
func (r *StreamProviderRepository) GetPrimaryByStreamID(streamID string) (*models.StreamProvider, error) {
	query := `
		SELECT id, stream_id, provider, provider_video_id, provider_url, metadata, created_at, updated_at
		FROM stream_providers
		WHERE stream_id = $1
		ORDER BY updated_at DESC, created_at DESC
		LIMIT 1
	`

	var sp models.StreamProvider
	var metadataRaw []byte
	err := r.db.QueryRow(query, streamID).Scan(
		&sp.ID,
		&sp.StreamID,
		&sp.Provider,
		&sp.ProviderVideoID,
		&sp.ProviderURL,
		&metadataRaw,
		&sp.CreatedAt,
		&sp.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load stream provider: %w", err)
	}

	if len(metadataRaw) > 0 {
		if err := json.Unmarshal(metadataRaw, &sp.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal provider metadata: %w", err)
		}
	}

	return &sp, nil
}

// Upsert creates or updates a provider row keyed by stream_id + provider.
func (r *StreamProviderRepository) Upsert(sp *models.StreamProvider) error {
	var metadataRaw []byte
	var err error
	if sp.Metadata != nil {
		metadataRaw, err = json.Marshal(sp.Metadata)
		if err != nil {
			return fmt.Errorf("marshal metadata: %w", err)
		}
	}

	query := `
		INSERT INTO stream_providers (stream_id, provider, provider_video_id, provider_url, metadata)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (stream_id, provider) DO UPDATE
		SET provider_video_id = EXCLUDED.provider_video_id,
			provider_url = EXCLUDED.provider_url,
			metadata = EXCLUDED.metadata,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at
	`

	err = r.db.QueryRow(
		query,
		sp.StreamID,
		sp.Provider,
		sp.ProviderVideoID,
		sp.ProviderURL,
		metadataRaw,
	).Scan(&sp.ID, &sp.CreatedAt, &sp.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert stream provider: %w", err)
	}

	return nil
}

// ListByProvider returns providers filtered by provider name.
func (r *StreamProviderRepository) ListByProvider(provider string) ([]models.StreamProvider, error) {
	query := `
		SELECT id, stream_id, provider, provider_video_id, provider_url, metadata, created_at, updated_at
		FROM stream_providers
		WHERE provider = $1
	`

	rows, err := r.db.Query(query, provider)
	if err != nil {
		return nil, fmt.Errorf("list stream providers: %w", err)
	}
	defer rows.Close()

	var providers []models.StreamProvider
	for rows.Next() {
		var sp models.StreamProvider
		var metadataRaw []byte
		if err := rows.Scan(
			&sp.ID,
			&sp.StreamID,
			&sp.Provider,
			&sp.ProviderVideoID,
			&sp.ProviderURL,
			&metadataRaw,
			&sp.CreatedAt,
			&sp.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan stream provider: %w", err)
		}
		if len(metadataRaw) > 0 {
			if err := json.Unmarshal(metadataRaw, &sp.Metadata); err != nil {
				return nil, fmt.Errorf("unmarshal stream provider metadata: %w", err)
			}
		}
		providers = append(providers, sp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stream providers: %w", err)
	}

	return providers, nil
}
