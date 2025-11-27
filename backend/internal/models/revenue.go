package models

import "time"

// RevenueShareMonthly represents monthly revenue share data for a race
type RevenueShareMonthly struct {
	ID                 string    `json:"id" db:"id"`
	RaceID             string    `json:"race_id" db:"race_id"`
	Year               int       `json:"year" db:"year"`
	Month              int       `json:"month" db:"month"`
	TotalRevenueCents  int       `json:"total_revenue_cents" db:"total_revenue_cents"`
	TotalWatchMinutes  float64   `json:"total_watch_minutes" db:"total_watch_minutes"`
	PlatformShareCents int       `json:"platform_share_cents" db:"platform_share_cents"`
	OrganizerShareCents int      `json:"organizer_share_cents" db:"organizer_share_cents"`
	CalculatedAt       time.Time `json:"calculated_at" db:"calculated_at"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// RevenueShareDetails includes race information
type RevenueShareDetails struct {
	ID                   string    `json:"id" db:"id"`
	RaceID               string    `json:"race_id" db:"race_id"`
	RaceName             string    `json:"race_name" db:"race_name"`
	Year                 int       `json:"year" db:"year"`
	Month                int       `json:"month" db:"month"`
	TotalRevenueCents    int       `json:"total_revenue_cents" db:"total_revenue_cents"`
	TotalRevenueDollars  float64   `json:"total_revenue_dollars" db:"total_revenue_dollars"`
	TotalWatchMinutes    float64   `json:"total_watch_minutes" db:"total_watch_minutes"`
	PlatformShareCents   int       `json:"platform_share_cents" db:"platform_share_cents"`
	PlatformShareDollars float64   `json:"platform_share_dollars" db:"platform_share_dollars"`
	OrganizerShareCents  int       `json:"organizer_share_cents" db:"organizer_share_cents"`
	OrganizerShareDollars float64  `json:"organizer_share_dollars" db:"organizer_share_dollars"`
	CalculatedAt         time.Time `json:"calculated_at" db:"calculated_at"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// RevenueSummary represents aggregated revenue data
type RevenueSummary struct {
	RaceID               string  `json:"race_id"`
	RaceName             string  `json:"race_name"`
	TotalRevenueCents    int     `json:"total_revenue_cents"`
	TotalRevenueDollars  float64 `json:"total_revenue_dollars"`
	TotalWatchMinutes    float64 `json:"total_watch_minutes"`
	PlatformShareCents   int     `json:"platform_share_cents"`
	PlatformShareDollars float64 `json:"platform_share_dollars"`
	OrganizerShareCents  int     `json:"organizer_share_cents"`
	OrganizerShareDollars float64 `json:"organizer_share_dollars"`
	MonthCount           int     `json:"month_count"`
}

