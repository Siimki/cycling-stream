package models

import "time"

type Achievement struct {
	ID          string    `json:"id" db:"id"`
	Slug        string    `json:"slug" db:"slug"`
	Title       string    `json:"title" db:"title"`
	Description *string   `json:"description,omitempty" db:"description"`
	Icon        *string   `json:"icon,omitempty" db:"icon"`
	Points      int       `json:"points" db:"points"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type AchievementSeed struct {
	Slug        string
	Title       string
	Description string
	Icon        string
	Points      int
}

type UserAchievement struct {
	ID            string                 `json:"id" db:"id"`
	UserID        string                 `json:"user_id" db:"user_id"`
	AchievementID string                 `json:"achievement_id" db:"achievement_id"`
	Slug          string                 `json:"slug" db:"slug"`
	Title         string                 `json:"title" db:"title"`
	Description   *string                `json:"description,omitempty" db:"description"`
	Icon          *string                `json:"icon,omitempty" db:"icon"`
	Points        int                    `json:"points" db:"points"`
	UnlockedAt    time.Time              `json:"unlocked_at" db:"unlocked_at"`
	Metadata      map[string]interface{} `json:"metadata" db:"metadata"`
}
