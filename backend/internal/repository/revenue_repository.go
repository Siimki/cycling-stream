package repository

import (
	"database/sql"
	"fmt"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/google/uuid"
)

type RevenueRepository struct {
	db *sql.DB
}

func NewRevenueRepository(db *sql.DB) *RevenueRepository {
	return &RevenueRepository{db: db}
}

// CalculateMonthlyRevenue calculates and stores monthly revenue share for a specific race and month
// Revenue split is 50/50 between platform and organizer
func (r *RevenueRepository) CalculateMonthlyRevenue(raceID string, year, month int) error {
	// Calculate total revenue from payments for this race in this month
	revenueQuery := `
		SELECT COALESCE(SUM(amount_cents), 0)
		FROM payments
		WHERE race_id = $1
		  AND status = 'succeeded'
		  AND EXTRACT(YEAR FROM created_at) = $2
		  AND EXTRACT(MONTH FROM created_at) = $3
	`

	var totalRevenueCents int
	err := r.db.QueryRow(revenueQuery, raceID, year, month).Scan(&totalRevenueCents)
	if err != nil {
		return fmt.Errorf("failed to calculate total revenue: %w", err)
	}

	// Calculate total watch minutes for this race in this month
	watchMinutesQuery := `
		SELECT COALESCE(SUM(duration_seconds) / 60.0, 0)
		FROM watch_sessions
		WHERE race_id = $1
		  AND duration_seconds IS NOT NULL
		  AND EXTRACT(YEAR FROM started_at) = $2
		  AND EXTRACT(MONTH FROM started_at) = $3
	`

	var totalWatchMinutes float64
	err = r.db.QueryRow(watchMinutesQuery, raceID, year, month).Scan(&totalWatchMinutes)
	if err != nil {
		return fmt.Errorf("failed to calculate total watch minutes: %w", err)
	}

	// Calculate 50/50 split
	platformShareCents := totalRevenueCents / 2
	organizerShareCents := totalRevenueCents - platformShareCents // Handle odd cents

	// Insert or update the monthly revenue record
	upsertQuery := `
		INSERT INTO revenue_share_monthly (
			id, race_id, year, month, total_revenue_cents, total_watch_minutes,
			platform_share_cents, organizer_share_cents, calculated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP)
		ON CONFLICT (race_id, year, month)
		DO UPDATE SET
			total_revenue_cents = EXCLUDED.total_revenue_cents,
			total_watch_minutes = EXCLUDED.total_watch_minutes,
			platform_share_cents = EXCLUDED.platform_share_cents,
			organizer_share_cents = EXCLUDED.organizer_share_cents,
			calculated_at = EXCLUDED.calculated_at,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err = r.db.Exec(
		upsertQuery,
		uuid.New().String(),
		raceID,
		year,
		month,
		totalRevenueCents,
		totalWatchMinutes,
		platformShareCents,
		organizerShareCents,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert monthly revenue: %w", err)
	}

	return nil
}

// GetMonthlyRevenueByRace gets monthly revenue data for a specific race
func (r *RevenueRepository) GetMonthlyRevenueByRace(raceID string) ([]models.RevenueShareDetails, error) {
	query := `
		SELECT id, race_id, race_name, year, month, total_revenue_cents,
		       total_revenue_dollars, total_watch_minutes, platform_share_cents,
		       platform_share_dollars, organizer_share_cents, organizer_share_dollars,
		       calculated_at, created_at, updated_at
		FROM revenue_share_details
		WHERE race_id = $1
		ORDER BY year DESC, month DESC
	`

	rows, err := r.db.Query(query, raceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query monthly revenue: %w", err)
	}
	defer rows.Close()

	var revenues []models.RevenueShareDetails
	for rows.Next() {
		var revenue models.RevenueShareDetails
		err := rows.Scan(
			&revenue.ID,
			&revenue.RaceID,
			&revenue.RaceName,
			&revenue.Year,
			&revenue.Month,
			&revenue.TotalRevenueCents,
			&revenue.TotalRevenueDollars,
			&revenue.TotalWatchMinutes,
			&revenue.PlatformShareCents,
			&revenue.PlatformShareDollars,
			&revenue.OrganizerShareCents,
			&revenue.OrganizerShareDollars,
			&revenue.CalculatedAt,
			&revenue.CreatedAt,
			&revenue.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan revenue: %w", err)
		}
		revenues = append(revenues, revenue)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating revenues: %w", err)
	}

	return revenues, nil
}

// GetAllMonthlyRevenue gets all monthly revenue data, optionally filtered by year and month
func (r *RevenueRepository) GetAllMonthlyRevenue(year, month *int) ([]models.RevenueShareDetails, error) {
	var query string
	var args []interface{}

	if year != nil && month != nil {
		query = `
			SELECT id, race_id, race_name, year, month, total_revenue_cents,
			       total_revenue_dollars, total_watch_minutes, platform_share_cents,
			       platform_share_dollars, organizer_share_cents, organizer_share_dollars,
			       calculated_at, created_at, updated_at
			FROM revenue_share_details
			WHERE year = $1 AND month = $2
			ORDER BY race_name, year DESC, month DESC
		`
		args = []interface{}{*year, *month}
	} else if year != nil {
		query = `
			SELECT id, race_id, race_name, year, month, total_revenue_cents,
			       total_revenue_dollars, total_watch_minutes, platform_share_cents,
			       platform_share_dollars, organizer_share_cents, organizer_share_dollars,
			       calculated_at, created_at, updated_at
			FROM revenue_share_details
			WHERE year = $1
			ORDER BY race_name, year DESC, month DESC
		`
		args = []interface{}{*year}
	} else {
		query = `
			SELECT id, race_id, race_name, year, month, total_revenue_cents,
			       total_revenue_dollars, total_watch_minutes, platform_share_cents,
			       platform_share_dollars, organizer_share_cents, organizer_share_dollars,
			       calculated_at, created_at, updated_at
			FROM revenue_share_details
			ORDER BY race_name, year DESC, month DESC
		`
		args = []interface{}{}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query monthly revenue: %w", err)
	}
	defer rows.Close()

	var revenues []models.RevenueShareDetails
	for rows.Next() {
		var revenue models.RevenueShareDetails
		err := rows.Scan(
			&revenue.ID,
			&revenue.RaceID,
			&revenue.RaceName,
			&revenue.Year,
			&revenue.Month,
			&revenue.TotalRevenueCents,
			&revenue.TotalRevenueDollars,
			&revenue.TotalWatchMinutes,
			&revenue.PlatformShareCents,
			&revenue.PlatformShareDollars,
			&revenue.OrganizerShareCents,
			&revenue.OrganizerShareDollars,
			&revenue.CalculatedAt,
			&revenue.CreatedAt,
			&revenue.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan revenue: %w", err)
		}
		revenues = append(revenues, revenue)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating revenues: %w", err)
	}

	return revenues, nil
}

// GetRevenueSummaryByRace gets aggregated revenue summary for a specific race
func (r *RevenueRepository) GetRevenueSummaryByRace(raceID string) (*models.RevenueSummary, error) {
	query := `
		SELECT 
			race_id,
			race_name,
			SUM(total_revenue_cents) as total_revenue_cents,
			SUM(total_revenue_cents) / 100.0 as total_revenue_dollars,
			SUM(total_watch_minutes) as total_watch_minutes,
			SUM(platform_share_cents) as platform_share_cents,
			SUM(platform_share_cents) / 100.0 as platform_share_dollars,
			SUM(organizer_share_cents) as organizer_share_cents,
			SUM(organizer_share_cents) / 100.0 as organizer_share_dollars,
			COUNT(*) as month_count
		FROM revenue_share_details
		WHERE race_id = $1
		GROUP BY race_id, race_name
	`

	var summary models.RevenueSummary
	err := r.db.QueryRow(query, raceID).Scan(
		&summary.RaceID,
		&summary.RaceName,
		&summary.TotalRevenueCents,
		&summary.TotalRevenueDollars,
		&summary.TotalWatchMinutes,
		&summary.PlatformShareCents,
		&summary.PlatformShareDollars,
		&summary.OrganizerShareCents,
		&summary.OrganizerShareDollars,
		&summary.MonthCount,
	)

	if err == sql.ErrNoRows {
		// Return empty summary if no revenue data found
		// Try to get race name
		raceQuery := `SELECT name FROM races WHERE id = $1`
		var raceName string
		if err := r.db.QueryRow(raceQuery, raceID).Scan(&raceName); err != nil {
			// If race not found, use empty string
			raceName = ""
		}
		
		return &models.RevenueSummary{
			RaceID:               raceID,
			RaceName:             raceName,
			TotalRevenueCents:    0,
			TotalRevenueDollars:  0,
			TotalWatchMinutes:    0,
			PlatformShareCents:   0,
			PlatformShareDollars: 0,
			OrganizerShareCents:  0,
			OrganizerShareDollars: 0,
			MonthCount:           0,
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue summary: %w", err)
	}

	return &summary, nil
}

// RecalculateAllMonthlyRevenue recalculates monthly revenue for all races with payments
func (r *RevenueRepository) RecalculateAllMonthlyRevenue() error {
	// Get all unique race_id, year, month combinations from payments
	query := `
		SELECT DISTINCT 
			race_id,
			EXTRACT(YEAR FROM created_at)::INTEGER as year,
			EXTRACT(MONTH FROM created_at)::INTEGER as month
		FROM payments
		WHERE race_id IS NOT NULL
		  AND status = 'succeeded'
		ORDER BY race_id, year, month
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query payment periods: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var raceID string
		var year, month int
		err := rows.Scan(&raceID, &year, &month)
		if err != nil {
			return fmt.Errorf("failed to scan payment period: %w", err)
		}

		// Calculate revenue for this race/month
		err = r.CalculateMonthlyRevenue(raceID, year, month)
		if err != nil {
			return fmt.Errorf("failed to calculate revenue for race %s, %d-%02d: %w", raceID, year, month, err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating payment periods: %w", err)
	}

	return nil
}

// RecalculateMonthlyRevenueForPeriod recalculates revenue for a specific year and month
func (r *RevenueRepository) RecalculateMonthlyRevenueForPeriod(year, month int) error {
	// Get all races with payments in this period
	query := `
		SELECT DISTINCT race_id
		FROM payments
		WHERE race_id IS NOT NULL
		  AND status = 'succeeded'
		  AND EXTRACT(YEAR FROM created_at) = $1
		  AND EXTRACT(MONTH FROM created_at) = $2
	`

	rows, err := r.db.Query(query, year, month)
	if err != nil {
		return fmt.Errorf("failed to query races for period: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var raceID string
		err := rows.Scan(&raceID)
		if err != nil {
			return fmt.Errorf("failed to scan race ID: %w", err)
		}

		// Calculate revenue for this race/month
		err = r.CalculateMonthlyRevenue(raceID, year, month)
		if err != nil {
			return fmt.Errorf("failed to calculate revenue for race %s: %w", raceID, err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating races: %w", err)
	}

	return nil
}

