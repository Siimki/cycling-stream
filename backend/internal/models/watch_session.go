package models

import "time"

type WatchSession struct {
	ID              string     `json:"id" db:"id"`
	UserID          string     `json:"user_id" db:"user_id"`
	RaceID          string     `json:"race_id" db:"race_id"`
	StartedAt       time.Time  `json:"started_at" db:"started_at"`
	EndedAt         *time.Time `json:"ended_at,omitempty" db:"ended_at"`
	DurationSeconds *int       `json:"duration_seconds,omitempty" db:"duration_seconds"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}

type WatchTimeStats struct {
	UserID       string     `json:"user_id" db:"user_id"`
	RaceID       string     `json:"race_id" db:"race_id"`
	SessionCount int        `json:"session_count" db:"session_count"`
	TotalSeconds int        `json:"total_seconds" db:"total_seconds"`
	TotalMinutes float64    `json:"total_minutes" db:"total_minutes"`
	FirstWatched time.Time  `json:"first_watched" db:"first_watched"`
	LastWatched  *time.Time `json:"last_watched,omitempty" db:"last_watched"`
}

type StartWatchSessionRequest struct {
	RaceID string `json:"race_id"`
}

type EndWatchSessionRequest struct {
	SessionID string `json:"session_id"`
}
