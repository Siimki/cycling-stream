package models

import "time"

type ChatMessage struct {
	ID        string     `json:"id" db:"id"`
	RaceID    string     `json:"race_id" db:"race_id"`
	UserID    *string    `json:"user_id,omitempty" db:"user_id"`
	Username  string     `json:"username" db:"username"`
	Message   string     `json:"message" db:"message"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

