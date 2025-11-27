package repository

import (
	"database/sql"
	"fmt"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/google/uuid"
)

type CostRepository struct {
	db *sql.DB
}

func NewCostRepository(db *sql.DB) *CostRepository {
	return &CostRepository{db: db}
}

func (r *CostRepository) Create(cost *models.Cost) error {
	query := `
		INSERT INTO costs (id, race_id, cost_type, amount_cents, year, month, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`

	if cost.ID == "" {
		cost.ID = uuid.New().String()
	}

	err := r.db.QueryRow(
		query,
		cost.ID,
		cost.RaceID,
		cost.CostType,
		cost.AmountCents,
		cost.Year,
		cost.Month,
		cost.Description,
	).Scan(&cost.CreatedAt, &cost.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create cost: %w", err)
	}

	return nil
}

func (r *CostRepository) GetByID(id string) (*models.Cost, error) {
	query := `
		SELECT id, race_id, cost_type, amount_cents, year, month, description, created_at, updated_at
		FROM costs
		WHERE id = $1
	`

	var cost models.Cost
	err := r.db.QueryRow(query, id).Scan(
		&cost.ID,
		&cost.RaceID,
		&cost.CostType,
		&cost.AmountCents,
		&cost.Year,
		&cost.Month,
		&cost.Description,
		&cost.CreatedAt,
		&cost.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cost: %w", err)
	}

	return &cost, nil
}

func (r *CostRepository) GetAll(year *int, month *int) ([]models.CostDetails, error) {
	query := `
		SELECT id, race_id, race_name, cost_type, amount_cents, amount_dollars, 
		       year, month, description, created_at, updated_at
		FROM cost_details
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if year != nil {
		query += fmt.Sprintf(" AND year = $%d", argIndex)
		args = append(args, *year)
		argIndex++
	}

	if month != nil {
		query += fmt.Sprintf(" AND month = $%d", argIndex)
		args = append(args, *month)
		argIndex++
	}

	query += " ORDER BY year DESC, month DESC, created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query costs: %w", err)
	}
	defer rows.Close()

	var costs []models.CostDetails
	for rows.Next() {
		var cost models.CostDetails
		err := rows.Scan(
			&cost.ID,
			&cost.RaceID,
			&cost.RaceName,
			&cost.CostType,
			&cost.AmountCents,
			&cost.AmountDollars,
			&cost.Year,
			&cost.Month,
			&cost.Description,
			&cost.CreatedAt,
			&cost.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cost: %w", err)
		}
		costs = append(costs, cost)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating costs: %w", err)
	}

	return costs, nil
}

func (r *CostRepository) GetByRace(raceID string, year *int, month *int) ([]models.CostDetails, error) {
	query := `
		SELECT id, race_id, race_name, cost_type, amount_cents, amount_dollars, 
		       year, month, description, created_at, updated_at
		FROM cost_details
		WHERE race_id = $1
	`
	args := []interface{}{raceID}
	argIndex := 2

	if year != nil {
		query += fmt.Sprintf(" AND year = $%d", argIndex)
		args = append(args, *year)
		argIndex++
	}

	if month != nil {
		query += fmt.Sprintf(" AND month = $%d", argIndex)
		args = append(args, *month)
		argIndex++
	}

	query += " ORDER BY year DESC, month DESC, created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query costs by race: %w", err)
	}
	defer rows.Close()

	var costs []models.CostDetails
	for rows.Next() {
		var cost models.CostDetails
		err := rows.Scan(
			&cost.ID,
			&cost.RaceID,
			&cost.RaceName,
			&cost.CostType,
			&cost.AmountCents,
			&cost.AmountDollars,
			&cost.Year,
			&cost.Month,
			&cost.Description,
			&cost.CreatedAt,
			&cost.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cost: %w", err)
		}
		costs = append(costs, cost)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating costs: %w", err)
	}

	return costs, nil
}

func (r *CostRepository) GetMonthlySummary(raceID *string, year *int, month *int) ([]models.CostSummaryMonthly, error) {
	query := `
		SELECT race_id, year, month, cdn_cents, server_cents, storage_cents, 
		       bandwidth_cents, other_cents, total_cents, total_dollars
		FROM cost_summary_monthly
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if raceID != nil {
		query += fmt.Sprintf(" AND race_id = $%d", argIndex)
		args = append(args, *raceID)
		argIndex++
	}

	if year != nil {
		query += fmt.Sprintf(" AND year = $%d", argIndex)
		args = append(args, *year)
		argIndex++
	}

	if month != nil {
		query += fmt.Sprintf(" AND month = $%d", argIndex)
		args = append(args, *month)
		argIndex++
	}

	query += " ORDER BY year DESC, month DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query cost summary: %w", err)
	}
	defer rows.Close()

	var summaries []models.CostSummaryMonthly
	for rows.Next() {
		var summary models.CostSummaryMonthly
		err := rows.Scan(
			&summary.RaceID,
			&summary.Year,
			&summary.Month,
			&summary.CDNCents,
			&summary.ServerCents,
			&summary.StorageCents,
			&summary.BandwidthCents,
			&summary.OtherCents,
			&summary.TotalCents,
			&summary.TotalDollars,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cost summary: %w", err)
		}
		summaries = append(summaries, summary)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cost summaries: %w", err)
	}

	return summaries, nil
}

func (r *CostRepository) Update(cost *models.Cost) error {
	query := `
		UPDATE costs
		SET race_id = $2, cost_type = $3, amount_cents = $4, year = $5, 
		    month = $6, description = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		cost.ID,
		cost.RaceID,
		cost.CostType,
		cost.AmountCents,
		cost.Year,
		cost.Month,
		cost.Description,
	).Scan(&cost.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("cost not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update cost: %w", err)
	}

	return nil
}

func (r *CostRepository) Delete(id string) error {
	query := `DELETE FROM costs WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete cost: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cost not found")
	}

	return nil
}

