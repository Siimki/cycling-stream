package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/google/uuid"
)

type MissionRepository struct {
	db *sql.DB
}

func NewMissionRepository(db *sql.DB) *MissionRepository {
	return &MissionRepository{db: db}
}

func (r *MissionRepository) Create(mission *models.Mission) error {
	mission.ID = uuid.New().String()
	query := `
		INSERT INTO missions (id, mission_type, title, description, points_reward, xp_reward, target_value, tier_number, category, requirement_json, valid_from, valid_until, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		mission.ID,
		mission.MissionType,
		mission.Title,
		mission.Description,
		mission.PointsReward,
		mission.XPReward,
		mission.TargetValue,
		mission.TierNumber,
		mission.Category,
		mission.RequirementJSON,
		mission.ValidFrom,
		mission.ValidUntil,
		mission.IsActive,
	).Scan(&mission.CreatedAt, &mission.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create mission: %w", err)
	}

	return nil
}

func (r *MissionRepository) GetByID(id string) (*models.Mission, error) {
	query := `
		SELECT id, mission_type, title, description, points_reward, xp_reward, target_value, 
		       tier_number, category, requirement_json, valid_from, valid_until, is_active, created_at, updated_at
		FROM missions
		WHERE id = $1
	`

	var mission models.Mission
	var requirementJSON sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&mission.ID,
		&mission.MissionType,
		&mission.Title,
		&mission.Description,
		&mission.PointsReward,
		&mission.XPReward,
		&mission.TargetValue,
		&mission.TierNumber,
		&mission.Category,
		&requirementJSON,
		&mission.ValidFrom,
		&mission.ValidUntil,
		&mission.IsActive,
		&mission.CreatedAt,
		&mission.UpdatedAt,
	)
	
	if requirementJSON.Valid {
		mission.RequirementJSON = &requirementJSON.String
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get mission: %w", err)
	}

	return &mission, nil
}

func (r *MissionRepository) GetActiveMissions() ([]models.Mission, error) {
	now := time.Now()
	query := `
		SELECT id, mission_type, title, description, points_reward, xp_reward, target_value, 
		       tier_number, category, requirement_json, valid_from, valid_until, is_active, created_at, updated_at
		FROM missions
		WHERE is_active = true
		  AND valid_from <= $1
		  AND (valid_until IS NULL OR valid_until >= $1)
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get active missions: %w", err)
	}
	defer rows.Close()

	var missions []models.Mission
	for rows.Next() {
		var mission models.Mission
		var requirementJSON sql.NullString
		err := rows.Scan(
			&mission.ID,
			&mission.MissionType,
			&mission.Title,
			&mission.Description,
			&mission.PointsReward,
			&mission.XPReward,
			&mission.TargetValue,
			&mission.TierNumber,
			&mission.Category,
			&requirementJSON,
			&mission.ValidFrom,
			&mission.ValidUntil,
			&mission.IsActive,
			&mission.CreatedAt,
			&mission.UpdatedAt,
		)
		if requirementJSON.Valid {
			mission.RequirementJSON = &requirementJSON.String
		}
		if err != nil {
			return nil, fmt.Errorf("failed to scan mission: %w", err)
		}
		missions = append(missions, mission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating missions: %w", err)
	}

	return missions, nil
}

func (r *MissionRepository) GetByType(missionType models.MissionType) ([]models.Mission, error) {
	now := time.Now()
	query := `
		SELECT id, mission_type, title, description, points_reward, xp_reward, target_value, 
		       tier_number, category, requirement_json, valid_from, valid_until, is_active, created_at, updated_at
		FROM missions
		WHERE mission_type = $1
		  AND is_active = true
		  AND valid_from <= $2
		  AND (valid_until IS NULL OR valid_until >= $2)
		ORDER BY tier_number ASC, created_at DESC
	`

	rows, err := r.db.Query(query, missionType, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get missions by type: %w", err)
	}
	defer rows.Close()

	var missions []models.Mission
	for rows.Next() {
		var mission models.Mission
		var requirementJSON sql.NullString
		err := rows.Scan(
			&mission.ID,
			&mission.MissionType,
			&mission.Title,
			&mission.Description,
			&mission.PointsReward,
			&mission.XPReward,
			&mission.TargetValue,
			&mission.TierNumber,
			&mission.Category,
			&requirementJSON,
			&mission.ValidFrom,
			&mission.ValidUntil,
			&mission.IsActive,
			&mission.CreatedAt,
			&mission.UpdatedAt,
		)
		if requirementJSON.Valid {
			mission.RequirementJSON = &requirementJSON.String
		}
		if err != nil {
			return nil, fmt.Errorf("failed to scan mission: %w", err)
		}
		missions = append(missions, mission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating missions: %w", err)
	}

	return missions, nil
}


// GetActiveCareerMissionsByType returns active career missions of a specific type, ordered by tier
func (r *MissionRepository) GetActiveCareerMissionsByType(missionType models.MissionType) ([]models.Mission, error) {
	now := time.Now()
	query := `
		SELECT id, mission_type, title, description, points_reward, xp_reward, target_value, 
		       tier_number, category, requirement_json, valid_from, valid_until, is_active, created_at, updated_at
		FROM missions
		WHERE mission_type = $1
		  AND category = 'career'
		  AND is_active = true
		  AND valid_from <= $2
		  AND (valid_until IS NULL OR valid_until >= $2)
		ORDER BY tier_number ASC
	`

	rows, err := r.db.Query(query, missionType, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get career missions by type: %w", err)
	}
	defer rows.Close()

	var missions []models.Mission
	for rows.Next() {
		var mission models.Mission
		var requirementJSON sql.NullString
		err := rows.Scan(
			&mission.ID,
			&mission.MissionType,
			&mission.Title,
			&mission.Description,
			&mission.PointsReward,
			&mission.XPReward,
			&mission.TargetValue,
			&mission.TierNumber,
			&mission.Category,
			&requirementJSON,
			&mission.ValidFrom,
			&mission.ValidUntil,
			&mission.IsActive,
			&mission.CreatedAt,
			&mission.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mission: %w", err)
		}
		if requirementJSON.Valid {
			mission.RequirementJSON = &requirementJSON.String
		}
		missions = append(missions, mission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating missions: %w", err)
	}

	return missions, nil
}

// GetWeeklyMissions returns all active weekly missions
func (r *MissionRepository) GetWeeklyMissions() ([]models.Mission, error) {
	now := time.Now()
	query := `
		SELECT id, mission_type, title, description, points_reward, xp_reward, target_value, 
		       tier_number, category, requirement_json, valid_from, valid_until, is_active, created_at, updated_at
		FROM missions
		WHERE category = 'weekly'
		  AND is_active = true
		  AND valid_from <= $1
		  AND (valid_until IS NULL OR valid_until >= $1)
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly missions: %w", err)
	}
	defer rows.Close()

	var missions []models.Mission
	for rows.Next() {
		var mission models.Mission
		var requirementJSON sql.NullString
		err := rows.Scan(
			&mission.ID,
			&mission.MissionType,
			&mission.Title,
			&mission.Description,
			&mission.PointsReward,
			&mission.XPReward,
			&mission.TargetValue,
			&mission.TierNumber,
			&mission.Category,
			&requirementJSON,
			&mission.ValidFrom,
			&mission.ValidUntil,
			&mission.IsActive,
			&mission.CreatedAt,
			&mission.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mission: %w", err)
		}
		if requirementJSON.Valid {
			mission.RequirementJSON = &requirementJSON.String
		}
		missions = append(missions, mission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating missions: %w", err)
	}

	return missions, nil
}
