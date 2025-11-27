package models

import "time"

type ViewerSession struct {
	ID          string     `json:"id" db:"id"`
	UserID      *string    `json:"user_id,omitempty" db:"user_id"`
	RaceID      string     `json:"race_id" db:"race_id"`
	SessionToken string    `json:"session_token" db:"session_token"`
	StartedAt   time.Time  `json:"started_at" db:"started_at"`
	LastSeenAt  time.Time  `json:"last_seen_at" db:"last_seen_at"`
	EndedAt     *time.Time `json:"ended_at,omitempty" db:"ended_at"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

type ConcurrentViewers struct {
	RaceID             string `json:"race_id" db:"race_id"`
	ConcurrentCount    int    `json:"concurrent_count" db:"concurrent_count"`
	AuthenticatedCount int    `json:"authenticated_count" db:"authenticated_count"`
	AnonymousCount     int    `json:"anonymous_count" db:"anonymous_count"`
}

type UniqueViewers struct {
	RaceID                   string `json:"race_id" db:"race_id"`
	UniqueViewerCount        int    `json:"unique_viewer_count" db:"unique_viewer_count"`
	UniqueAuthenticatedCount int    `json:"unique_authenticated_count" db:"unique_authenticated_count"`
	UniqueAnonymousCount     int    `json:"unique_anonymous_count" db:"unique_anonymous_count"`
}

type StartViewerSessionRequest struct {
	RaceID string `json:"race_id"`
}

type EndViewerSessionRequest struct {
	SessionID string `json:"session_id"`
}

type HeartbeatViewerSessionRequest struct {
	SessionID string `json:"session_id"`
}

