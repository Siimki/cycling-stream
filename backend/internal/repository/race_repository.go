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
		SELECT r.id, r.name, r.description, r.start_date, r.end_date, r.location, r.category, 
		       r.is_free, r.price_cents, r.requires_login, r.stage_name, r.stage_type, r.elevation_meters, 
		       r.estimated_finish_time, r.stage_length_km, r.created_at, r.updated_at, s.status AS stream_status
		FROM races r
		LEFT JOIN streams s ON r.id = s.race_id
		ORDER BY r.start_date DESC NULLS LAST, r.created_at DESC
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
			&race.RequiresLogin,
			&race.StageName,
			&race.StageType,
			&race.ElevationMeters,
			&race.EstimatedFinishTime,
			&race.StageLengthKm,
			&race.CreatedAt,
			&race.UpdatedAt,
			&race.StreamStatus,
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
		       is_free, price_cents, requires_login, stage_name, stage_type, elevation_meters, 
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
		&race.RequiresLogin,
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
		INSERT INTO races (name, description, start_date, end_date, location, category, is_free, price_cents, requires_login,
		                   stage_name, stage_type, elevation_meters, estimated_finish_time, stage_length_km)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
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
		race.RequiresLogin,
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
		    location = $6, category = $7, is_free = $8, price_cents = $9, requires_login = $10,
		    stage_name = $11, stage_type = $12, elevation_meters = $13,
		    estimated_finish_time = $14, stage_length_km = $15, updated_at = CURRENT_TIMESTAMP
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
		race.RequiresLogin,
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

// GetRacesByCategory returns races with the same category
func (r *RaceRepository) GetRacesByCategory(category string, limit int) ([]models.Race, error) {
	query := `
		SELECT id, name, description, start_date, end_date, location, category, 
		       is_free, price_cents, requires_login, stage_name, stage_type, elevation_meters, 
		       estimated_finish_time, stage_length_km, created_at, updated_at
		FROM races
		WHERE category = $1
		ORDER BY start_date DESC NULLS LAST
		LIMIT $2
	`

	return r.queryRaces(query, category, limit)
}

// GetUpcomingRaces returns races with start_date in the future
func (r *RaceRepository) GetUpcomingRaces(limit int) ([]models.Race, error) {
	query := `
		SELECT id, name, description, start_date, end_date, location, category, 
		       is_free, price_cents, requires_login, stage_name, stage_type, elevation_meters, 
		       estimated_finish_time, stage_length_km, created_at, updated_at
		FROM races
		WHERE start_date > CURRENT_TIMESTAMP
		ORDER BY start_date ASC
		LIMIT $1
	`

	return r.queryRaces(query, limit)
}

// GetLiveRaces returns races that are currently live (have a live stream)
func (r *RaceRepository) GetLiveRaces() ([]models.Race, error) {
	query := `
		SELECT DISTINCT r.id, r.name, r.description, r.start_date, r.end_date, r.location, r.category, 
		       r.is_free, r.price_cents, r.requires_login, r.stage_name, r.stage_type, r.elevation_meters, 
		       r.estimated_finish_time, r.stage_length_km, r.created_at, r.updated_at
		FROM races r
		INNER JOIN streams s ON r.id = s.race_id
		WHERE s.status = 'live'
		ORDER BY r.start_date DESC
	`

	return r.queryRaces(query)
}

// GetSimilarRaces returns races similar to the given race (same category, similar elevation/distance)
func (r *RaceRepository) GetSimilarRaces(raceID string, limit int) ([]models.Race, error) {
	// First get the reference race
	race, err := r.GetByID(raceID)
	if err != nil || race == nil {
		return []models.Race{}, nil
	}

	query := `
		SELECT id, name, description, start_date, end_date, location, category, 
		       is_free, price_cents, requires_login, stage_name, stage_type, elevation_meters, 
		       estimated_finish_time, stage_length_km, created_at, updated_at
		FROM races
		WHERE id != $1
		  AND (
		    category = $2
		    OR (elevation_meters IS NOT NULL AND $3 IS NOT NULL 
		        AND ABS(elevation_meters - $3) < 500)
		    OR (stage_length_km IS NOT NULL AND $4 IS NOT NULL 
		        AND ABS(stage_length_km - $4) < 50)
		  )
		ORDER BY 
		  CASE WHEN category = $2 THEN 1 ELSE 2 END,
		  start_date DESC
		LIMIT $5
	`

	var elevation *int
	if race.ElevationMeters != nil {
		elevation = race.ElevationMeters
	}
	var length *float64
	if race.StageLengthKm != nil {
		lengthVal := float64(*race.StageLengthKm)
		length = &lengthVal
	}

	return r.queryRaces(query, raceID, race.Category, elevation, length, limit)
}

// Helper method to query races with variable parameters
func (r *RaceRepository) queryRaces(query string, args ...interface{}) ([]models.Race, error) {
	rows, err := r.db.Query(query, args...)
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
			&race.RequiresLogin,
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

	if races == nil {
		races = []models.Race{}
	}

	return races, nil
}
