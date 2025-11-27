package models

import "time"

type Race struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	StartDate   *time.Time `json:"start_date,omitempty" db:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
	Location    *string   `json:"location,omitempty" db:"location"`
	Category    *string   `json:"category,omitempty" db:"category"`
	IsFree      bool      `json:"is_free" db:"is_free"`
	PriceCents  int       `json:"price_cents" db:"price_cents"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

