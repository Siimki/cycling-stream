package models

import "time"

// PredictionOption represents a single option in a prediction market
type PredictionOption struct {
	ID   string  `json:"id"`   // UUID
	Text string  `json:"text"` // e.g., "Yes", "No", "Breakaway wins"
	Odds float64 `json:"odds"` // e.g., 3.0, 1.5
}

// PredictionMarket represents a prediction market for a race
type PredictionMarket struct {
	ID              string            `json:"id" db:"id"`
	RaceID          string            `json:"race_id" db:"race_id"`
	Question        string            `json:"question" db:"question"`
	Options         []PredictionOption `json:"options" db:"options"` // Stored as JSONB
	Status          string            `json:"status" db:"status"` // open, settled, cancelled
	SettledOptionID *string           `json:"settled_option_id,omitempty" db:"settled_option_id"`
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
	SettledAt       *time.Time        `json:"settled_at,omitempty" db:"settled_at"`
}

// PredictionBet represents a user's bet on a prediction market
type PredictionBet struct {
	ID             string     `json:"id" db:"id"`
	UserID         string     `json:"user_id" db:"user_id"`
	MarketID       string     `json:"market_id" db:"market_id"`
	OptionID       string     `json:"option_id" db:"option_id"`
	StakePoints    int        `json:"stake_points" db:"stake_points"`
	PotentialPayout int       `json:"potential_payout" db:"potential_payout"`
	Result         string     `json:"result" db:"result"` // pending, won, lost
	PayoutPoints   *int       `json:"payout_points,omitempty" db:"payout_points"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	SettledAt      *time.Time `json:"settled_at,omitempty" db:"settled_at"`
}

// PlaceBetRequest represents a request to place a bet
type PlaceBetRequest struct {
	OptionID    string `json:"option_id"`
	StakePoints int    `json:"stake_points"`
}

// PredictionBetWithMarket includes market details with the bet
type PredictionBetWithMarket struct {
	PredictionBet
	Market PredictionMarket `json:"market"`
}


