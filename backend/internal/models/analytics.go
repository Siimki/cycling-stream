package models

import "time"

// AnalyticsEvent represents a client-side player event.
type AnalyticsEvent struct {
	Type             string                 `json:"type"`
	VideoTimeSeconds *int                   `json:"videoTime,omitempty"`
	Extra            map[string]interface{} `json:"extra,omitempty"`
}

// AnalyticsEventBatch is the ingestion payload for /analytics/events.
type AnalyticsEventBatch struct {
	StreamID string           `json:"streamId"`
	ClientID string           `json:"clientId"`
	Events   []AnalyticsEvent `json:"events"`
}

// PlaybackEvent is the persisted representation of a player event.
type PlaybackEvent struct {
	ID               string                 `json:"id" db:"id"`
	StreamID         string                 `json:"stream_id" db:"stream_id"`
	ViewerSessionID  *string                `json:"viewer_session_id,omitempty" db:"viewer_session_id"`
	ClientID         string                 `json:"client_id" db:"client_id"`
	EventType        string                 `json:"event_type" db:"event_type"`
	VideoTimeSeconds *int                   `json:"video_time_seconds,omitempty" db:"video_time_seconds"`
	Country          string                 `json:"country" db:"country"`
	DeviceType       string                 `json:"device_type" db:"device_type"`
	Extra            map[string]interface{} `json:"extra,omitempty" db:"extra"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at"`
}

type StreamStats struct {
	StreamID              string                 `json:"stream_id" db:"stream_id"`
	UniqueViewers         int                    `json:"unique_viewers" db:"unique_viewers"`
	TotalWatchSeconds     int64                  `json:"total_watch_seconds" db:"total_watch_seconds"`
	AvgWatchSeconds       int                    `json:"avg_watch_seconds" db:"avg_watch_seconds"`
	PeakConcurrentViewers int                    `json:"peak_concurrent_viewers" db:"peak_concurrent_viewers"`
	TopCountries          map[string]int         `json:"top_countries,omitempty" db:"top_countries"`
	DeviceBreakdown       map[string]int         `json:"device_breakdown,omitempty" db:"device_breakdown"`
	BufferSeconds         int64                  `json:"buffer_seconds" db:"buffer_seconds"`
	BufferRatio           float64                `json:"buffer_ratio" db:"buffer_ratio"`
	ErrorRate             float64                `json:"error_rate" db:"error_rate"`
	LastCalculatedAt      time.Time              `json:"last_calculated_at" db:"last_calculated_at"`
	CreatedAt             time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at" db:"updated_at"`
}

type StreamStatsSummary struct {
	StreamCount        int     `json:"stream_count"`
	TotalUniqueViewers int64   `json:"total_unique_viewers"`
	TotalWatchSeconds  int64   `json:"total_watch_seconds"`
	AvgPeakConcurrent  float64 `json:"avg_peak_concurrent"`
}
