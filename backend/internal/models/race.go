package models

import "time"

type Race struct {
	ID                 string    `json:"id" db:"id"`
	Name               string    `json:"name" db:"name"`
	Description        *string   `json:"description,omitempty" db:"description"`
	StartDate          *time.Time `json:"start_date,omitempty" db:"start_date"`
	EndDate            *time.Time `json:"end_date,omitempty" db:"end_date"`
	Location           *string   `json:"location,omitempty" db:"location"`
	Category           *string   `json:"category,omitempty" db:"category"`
	IsFree             bool      `json:"is_free" db:"is_free"`
	PriceCents         int       `json:"price_cents" db:"price_cents"`
	StageName          *string   `json:"stage_name,omitempty" db:"stage_name"`
	StageType          *string   `json:"stage_type,omitempty" db:"stage_type"`
	ElevationMeters    *int      `json:"elevation_meters,omitempty" db:"elevation_meters"`
	EstimatedFinishTime *string  `json:"estimated_finish_time,omitempty" db:"estimated_finish_time"`
	StageLengthKm      *int      `json:"stage_length_km,omitempty" db:"stage_length_km"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

