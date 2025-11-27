package models

import "time"

type Stream struct {
	ID         string    `json:"id" db:"id"`
	RaceID     string    `json:"race_id" db:"race_id"`
	Status     string    `json:"status" db:"status"` // planned, live, ended
	StreamType string    `json:"stream_type" db:"stream_type"` // hls, youtube
	SourceID   *string   `json:"source_id,omitempty" db:"source_id"`
	OriginURL  *string   `json:"origin_url,omitempty" db:"origin_url"`
	CDNURL     *string   `json:"cdn_url,omitempty" db:"cdn_url"`
	StreamKey  *string   `json:"stream_key,omitempty" db:"stream_key"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type StreamResponse struct {
	Status     string  `json:"status"`
	StreamType string  `json:"stream_type"`
	SourceID   *string `json:"source_id,omitempty"`
	OriginURL  *string `json:"origin_url,omitempty"`
	CDNURL     *string `json:"cdn_url,omitempty"`
}
