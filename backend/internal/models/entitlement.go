package models

import "time"

type Entitlement struct {
	ID        string     `json:"id" db:"id"`
	UserID    string     `json:"user_id" db:"user_id"`
	RaceID    string     `json:"race_id" db:"race_id"`
	Type      string     `json:"type" db:"type"` // ticket, subscription
	ExpiresAt *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

