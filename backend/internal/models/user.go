package models

import "time"

type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Name         *string   `json:"name,omitempty" db:"name"`
	Bio          string    `json:"bio" db:"bio"`
	Points       int       `json:"points" db:"points"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// PublicUser represents user data safe to expose publicly (no email)
type PublicUser struct {
	ID                string    `json:"id"`
	Name              *string   `json:"name,omitempty"`
	Bio               string    `json:"bio"`
	Points            int       `json:"points"`
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

