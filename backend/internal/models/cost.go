package models

import "time"

// CostType represents the type of cost
type CostType string

const (
	CostTypeCDN       CostType = "cdn"
	CostTypeServer    CostType = "server"
	CostTypeStorage   CostType = "storage"
	CostTypeBandwidth CostType = "bandwidth"
	CostTypeOther     CostType = "other"
)

// Cost represents a cost entry in the database
type Cost struct {
	ID          string    `json:"id" db:"id"`
	RaceID      *string   `json:"race_id,omitempty" db:"race_id"`
	CostType    CostType  `json:"cost_type" db:"cost_type"`
	AmountCents int       `json:"amount_cents" db:"amount_cents"`
	Year        int       `json:"year" db:"year"`
	Month       int       `json:"month" db:"month"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CostDetails represents a cost with race information (from view)
type CostDetails struct {
	ID            string    `json:"id" db:"id"`
	RaceID        *string   `json:"race_id,omitempty" db:"race_id"`
	RaceName      *string   `json:"race_name,omitempty" db:"race_name"`
	CostType      CostType  `json:"cost_type" db:"cost_type"`
	AmountCents   int       `json:"amount_cents" db:"amount_cents"`
	AmountDollars float64   `json:"amount_dollars" db:"amount_dollars"`
	Year          int       `json:"year" db:"year"`
	Month         int       `json:"month" db:"month"`
	Description   *string   `json:"description,omitempty" db:"description"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// CostSummaryMonthly represents monthly cost aggregation
type CostSummaryMonthly struct {
	RaceID       *string  `json:"race_id,omitempty" db:"race_id"`
	Year         int      `json:"year" db:"year"`
	Month        int      `json:"month" db:"month"`
	CDNCents     int      `json:"cdn_cents" db:"cdn_cents"`
	ServerCents  int      `json:"server_cents" db:"server_cents"`
	StorageCents int      `json:"storage_cents" db:"storage_cents"`
	BandwidthCents int    `json:"bandwidth_cents" db:"bandwidth_cents"`
	OtherCents   int      `json:"other_cents" db:"other_cents"`
	TotalCents   int      `json:"total_cents" db:"total_cents"`
	TotalDollars float64  `json:"total_dollars" db:"total_dollars"`
}

// CreateCostRequest represents a request to create a cost
type CreateCostRequest struct {
	RaceID      *string  `json:"race_id,omitempty"`
	CostType    CostType `json:"cost_type"`
	AmountCents int      `json:"amount_cents"`
	Year        int      `json:"year"`
	Month       int      `json:"month"`
	Description *string  `json:"description,omitempty"`
}

