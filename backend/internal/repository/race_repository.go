package repository

import (
	"database/sql"
	"fmt"

	"github.com/cyclingstream/backend/internal/models"
)

type RaceRepository struct {
	db *sql.DB
}

func NewRaceRepository(db *sql.DB) *RaceRepository {
	return &RaceRepository{db: db}
}

func (r *RaceRepository) GetAll() ([]models.Race, error) {
	query := `
		SELECT id, name, description, start_date, end_date, location, category, 
		       is_free, price_cents, stage_name, stage_type, elevation_meters, 
		       estimated_finish_time, stage_length_km, created_at, updated_at
		FROM races
		ORDER BY start_date DESC NULLS LAST, created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query races: %w", err)
	}
	defer rows.Close()

	races := make([]models.Race, 0)
	for rows.Next() {
		var race models.Race
		err := rows.Scan(
			&race.ID,
			&race.Name,
			&race.Description,
			&race.StartDate,
			&race.EndDate,
			&race.Location,
			&race.Category,
			&race.IsFree,
			&race.PriceCents,
			&race.StageName,
			&race.StageType,
			&race.ElevationMeters,
			&race.EstimatedFinishTime,
			&race.StageLengthKm,
			&race.CreatedAt,
			&race.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan race: %w", err)
		}
		races = append(races, race)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating races: %w", err)
	}

	// Ensure we always return an empty slice, not nil, so JSON encodes as [] not null
	if races == nil {
		races = []models.Race{}
	}

	return races, nil
}

func (r *RaceRepository) GetByID(id string) (*models.Race, error) {
	query := `
		SELECT id, name, description, start_date, end_date, location, category, 
		       is_free, price_cents, stage_name, stage_type, elevation_meters, 
		       estimated_finish_time, stage_length_km, created_at, updated_at
		FROM races
		WHERE id = $1
	`

	var race models.Race
	err := r.db.QueryRow(query, id).Scan(
		&race.ID,
		&race.Name,
		&race.Description,
		&race.StartDate,
		&race.EndDate,
		&race.Location,
		&race.Category,
		&race.IsFree,
		&race.PriceCents,
		&race.StageName,
		&race.StageType,
		&race.ElevationMeters,
		&race.EstimatedFinishTime,
		&race.StageLengthKm,
		&race.CreatedAt,
		&race.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get race: %w", err)
	}

	return &race, nil
}

func (r *RaceRepository) Create(race *models.Race) error {
	query := `
		INSERT INTO races (name, description, start_date, end_date, location, category, is_free, price_cents,
		                   stage_name, stage_type, elevation_meters, estimated_finish_time, stage_length_km)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		race.Name,
		race.Description,
		race.StartDate,
		race.EndDate,
		race.Location,
		race.Category,
		race.IsFree,
		race.PriceCents,
		race.StageName,
		race.StageType,
		race.ElevationMeters,
		race.EstimatedFinishTime,
		race.StageLengthKm,
	).Scan(&race.ID, &race.CreatedAt, &race.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create race: %w", err)
	}

	return nil
}

func (r *RaceRepository) Update(race *models.Race) error {
	query := `
		UPDATE races
		SET name = $2, description = $3, start_date = $4, end_date = $5, 
		    location = $6, category = $7, is_free = $8, price_cents = $9,
		    stage_name = $10, stage_type = $11, elevation_meters = $12,
		    estimated_finish_time = $13, stage_length_km = $14, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		race.ID,
		race.Name,
		race.Description,
		race.StartDate,
		race.EndDate,
		race.Location,
		race.Category,
		race.IsFree,
		race.PriceCents,
		race.StageName,
		race.StageType,
		race.ElevationMeters,
		race.EstimatedFinishTime,
		race.StageLengthKm,
	).Scan(&race.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("race not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update race: %w", err)
	}

	return nil
}

func (r *RaceRepository) Delete(id string) error {
	query := `DELETE FROM races WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete race: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("race not found")
	}

	return nil
}
