package models

import "time"

type UserFavorite struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	FavoriteType string    `json:"favorite_type" db:"favorite_type"`
	FavoriteID   string    `json:"favorite_id" db:"favorite_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type AddFavoriteRequest struct {
	FavoriteType string `json:"favorite_type"`
	FavoriteID   string `json:"favorite_id"`
}

type WatchHistoryEntry struct {
	UserID          string     `json:"user_id" db:"user_id"`
	RaceID          string     `json:"race_id" db:"race_id"`
	RaceName        string     `json:"race_name" db:"race_name"`
	RaceCategory    *string    `json:"race_category,omitempty" db:"race_category"`
	RaceStartDate   *time.Time `json:"race_start_date,omitempty" db:"race_start_date"`
	SessionCount    int        `json:"session_count" db:"session_count"`
	TotalSeconds    int        `json:"total_seconds" db:"total_seconds"`
	TotalMinutes    float64    `json:"total_minutes" db:"total_minutes"`
	FirstWatched    time.Time  `json:"first_watched" db:"first_watched"`
	LastWatched     time.Time  `json:"last_watched" db:"last_watched"`
	LikelyCompleted bool       `json:"likely_completed" db:"likely_completed"`
}

