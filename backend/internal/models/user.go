package models

import "time"

type User struct {
	ID              string    `json:"id" db:"id"`
	Email           string    `json:"email" db:"email"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	Name            *string   `json:"name,omitempty" db:"name"`
	Bio             string    `json:"bio" db:"bio"`
	Points          int       `json:"points" db:"points"`
	XPTotal         int       `json:"xp_total" db:"xp_total"`
	Level           int       `json:"level" db:"level"`
	BestStreakWeeks int       `json:"best_streak_weeks" db:"best_streak_weeks"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// PublicUser represents user data safe to expose publicly (no email)
type PublicUser struct {
	ID                string    `json:"id"`
	Name              *string   `json:"name,omitempty"`
	Bio               string    `json:"bio"`
	Points            int       `json:"points"`
	XPTotal           int       `json:"xp_total"`
	Level             int       `json:"level"`
	BestStreakWeeks   int       `json:"best_streak_weeks"`
	TotalWatchMinutes int       `json:"total_watch_minutes"`
	CreatedAt         time.Time `json:"created_at"`
}

// LeaderboardEntry represents a user entry in the leaderboard
type LeaderboardEntry struct {
	ID                string  `json:"id"`
	Name              *string `json:"name,omitempty"`
	Points            int     `json:"points"`
	TotalWatchMinutes int     `json:"total_watch_minutes"`
}

// XPProgress represents XP and level progress information
type XPProgress struct {
	XPTotal          int `json:"xp_total"`
	Level            int `json:"level"`
	CurrentXPInLevel int `json:"current_xp_in_level"`
	XPNeededForNext  int `json:"xp_needed_for_next"`
}

