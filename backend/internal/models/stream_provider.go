package models

import "time"

type StreamProvider struct {
	ID              string                 `json:"id" db:"id"`
	StreamID        string                 `json:"stream_id" db:"stream_id"`
	Provider        string                 `json:"provider" db:"provider"`
	ProviderVideoID string                 `json:"provider_video_id" db:"provider_video_id"`
	ProviderURL     *string                `json:"provider_url,omitempty" db:"provider_url"`
	Metadata        map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
}
