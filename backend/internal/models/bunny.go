package models

import "time"

type BunnyVideoStats struct {
	ID               string                 `json:"id" db:"id"`
	BunnyVideoID     string                 `json:"bunny_video_id" db:"bunny_video_id"`
	StreamID         *string                `json:"stream_id,omitempty" db:"stream_id"`
	Date             time.Time              `json:"date" db:"date"`
	Views            int                    `json:"views" db:"views"`
	WatchTimeSeconds int64                  `json:"watch_time_seconds" db:"watch_time_seconds"`
	GeoBreakdown     map[string]int         `json:"geo_breakdown,omitempty" db:"geo_breakdown"`
	RawPayload       map[string]interface{} `json:"raw_payload,omitempty" db:"raw_payload"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at" db:"updated_at"`
}

type BunnyAnalyticsResponse struct {
	Views            int                    `json:"views"`
	WatchTimeSeconds int64                  `json:"watchTime"`
	Geo              map[string]int         `json:"geo"`
	Raw              map[string]interface{} `json:"raw"`
}
