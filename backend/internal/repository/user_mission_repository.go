package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/google/uuid"
)

type UserMissionRepository struct {
	db *sql.DB
}

func NewUserMissionRepository(db *sql.DB) *UserMissionRepository {
	return &UserMissionRepository{db: db}
}

func (r *UserMissionRepository) GetOrCreate(userID, missionID string) (*models.UserMission, error) {
	// Try to get existing user mission
	userMission, err := r.GetByUserAndMission(userID, missionID)
	if err != nil {
		return nil, err
	}
	if userMission != nil {
		return userMission, nil
	}

	// Create new user mission
	userMission = &models.UserMission{
		ID:        uuid.New().String(),
		UserID:    userID,
		MissionID: missionID,
		Progress:  0,
	}

	query := `
		INSERT INTO user_missions (id, user_id, mission_id, progress)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at
	`

	err = r.db.QueryRow(
		query,
		userMission.ID,
		userMission.UserID,
		userMission.MissionID,
		userMission.Progress,
	).Scan(&userMission.CreatedAt, &userMission.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user mission: %w", err)
	}

	return userMission, nil
}

func (r *UserMissionRepository) GetByUserAndMission(userID, missionID string) (*models.UserMission, error) {
	query := `
		SELECT id, user_id, mission_id, progress, completed_at, claimed_at, created_at, updated_at
		FROM user_missions
		WHERE user_id = $1 AND mission_id = $2
	`

	var userMission models.UserMission
	var completedAt sql.NullTime
	var claimedAt sql.NullTime

	err := r.db.QueryRow(query, userID, missionID).Scan(
		&userMission.ID,
		&userMission.UserID,
		&userMission.MissionID,
		&userMission.Progress,
		&completedAt,
		&claimedAt,
		&userMission.CreatedAt,
		&userMission.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user mission: %w", err)
	}

	if completedAt.Valid {
		userMission.CompletedAt = &completedAt.Time
	}
	if claimedAt.Valid {
		userMission.ClaimedAt = &claimedAt.Time
	}

	return &userMission, nil
}

func (r *UserMissionRepository) GetByUserID(userID string) ([]models.UserMissionWithDetails, error) {
	query := `
		SELECT 
			um.id, um.user_id, um.mission_id, um.progress, um.completed_at, um.claimed_at, 
			um.created_at, um.updated_at,
			m.id, m.mission_type, m.title, m.description, m.points_reward, m.xp_reward, m.target_value,
			m.tier_number, m.category, m.requirement_json, m.valid_from, m.valid_until, m.is_active, m.created_at, m.updated_at
		FROM user_missions um
		INNER JOIN missions m ON um.mission_id = m.id
		WHERE um.user_id = $1
		ORDER BY um.created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user missions: %w", err)
	}
	defer rows.Close()

	var userMissions []models.UserMissionWithDetails
	for rows.Next() {
		var um models.UserMissionWithDetails
		var completedAt sql.NullTime
		var claimedAt sql.NullTime
		var missionValidUntil sql.NullTime
		var requirementJSON sql.NullString

		err := rows.Scan(
			&um.UserMission.ID,
			&um.UserMission.UserID,
			&um.UserMission.MissionID,
			&um.UserMission.Progress,
			&completedAt,
			&claimedAt,
			&um.UserMission.CreatedAt,
			&um.UserMission.UpdatedAt,
			&um.Mission.ID,
			&um.Mission.MissionType,
			&um.Mission.Title,
			&um.Mission.Description,
			&um.Mission.PointsReward,
			&um.Mission.XPReward,
			&um.Mission.TargetValue,
			&um.Mission.TierNumber,
			&um.Mission.Category,
			&requirementJSON,
			&um.Mission.ValidFrom,
			&missionValidUntil,
			&um.Mission.IsActive,
			&um.Mission.CreatedAt,
			&um.Mission.UpdatedAt,
		)
		if requirementJSON.Valid {
			um.Mission.RequirementJSON = &requirementJSON.String
		}
		if err != nil {
			return nil, fmt.Errorf("failed to scan user mission: %w", err)
		}

		if completedAt.Valid {
			um.UserMission.CompletedAt = &completedAt.Time
		}
		if claimedAt.Valid {
			um.UserMission.ClaimedAt = &claimedAt.Time
		}
		if missionValidUntil.Valid {
			um.Mission.ValidUntil = &missionValidUntil.Time
		}

		userMissions = append(userMissions, um)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user missions: %w", err)
	}

	return userMissions, nil
}

func (r *UserMissionRepository) UpdateProgress(userID, missionID string, progress int) error {
	query := `
		UPDATE user_missions
		SET progress = $1, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $2 AND mission_id = $3
	`

	result, err := r.db.Exec(query, progress, userID, missionID)
	if err != nil {
		return fmt.Errorf("failed to update progress: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user mission not found")
	}

	return nil
}

func (r *UserMissionRepository) IncrementProgress(userID, missionID string, increment int) error {
	query := `
		UPDATE user_missions
		SET progress = progress + $1, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $2 AND mission_id = $3
	`

	result, err := r.db.Exec(query, increment, userID, missionID)
	if err != nil {
		return fmt.Errorf("failed to increment progress: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user mission not found")
	}

	return nil
}

func (r *UserMissionRepository) Complete(userID, missionID string) error {
	now := time.Now()
	query := `
		UPDATE user_missions
		SET completed_at = $1, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $2 AND mission_id = $3
		  AND completed_at IS NULL
	`

	result, err := r.db.Exec(query, now, userID, missionID)
	if err != nil {
		return fmt.Errorf("failed to complete mission: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user mission not found or already completed")
	}

	return nil
}

func (r *UserMissionRepository) Claim(userID, missionID string) error {
	now := time.Now()
	query := `
		UPDATE user_missions
		SET claimed_at = $1, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $2 AND mission_id = $3
		  AND completed_at IS NOT NULL
		  AND claimed_at IS NULL
	`

	result, err := r.db.Exec(query, now, userID, missionID)
	if err != nil {
		return fmt.Errorf("failed to claim mission: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user mission not found, not completed, or already claimed")
	}

	return nil
}

func (r *UserMissionRepository) GetByType(userID string, missionType models.MissionType) ([]models.UserMissionWithDetails, error) {
	query := `
		SELECT 
			um.id, um.user_id, um.mission_id, um.progress, um.completed_at, um.claimed_at, 
			um.created_at, um.updated_at,
			m.id, m.mission_type, m.title, m.description, m.points_reward, m.xp_reward, m.target_value,
			m.tier_number, m.category, m.requirement_json, m.valid_from, m.valid_until, m.is_active, m.created_at, m.updated_at
		FROM user_missions um
		INNER JOIN missions m ON um.mission_id = m.id
		WHERE um.user_id = $1 AND m.mission_type = $2
		ORDER BY m.tier_number ASC, um.created_at DESC
	`

	rows, err := r.db.Query(query, userID, missionType)
	if err != nil {
		return nil, fmt.Errorf("failed to get user missions by type: %w", err)
	}
	defer rows.Close()

	var userMissions []models.UserMissionWithDetails
	for rows.Next() {
		var um models.UserMissionWithDetails
		var completedAt sql.NullTime
		var claimedAt sql.NullTime
		var missionValidUntil sql.NullTime
		var requirementJSON sql.NullString

		err := rows.Scan(
			&um.UserMission.ID,
			&um.UserMission.UserID,
			&um.UserMission.MissionID,
			&um.UserMission.Progress,
			&completedAt,
			&claimedAt,
			&um.UserMission.CreatedAt,
			&um.UserMission.UpdatedAt,
			&um.Mission.ID,
			&um.Mission.MissionType,
			&um.Mission.Title,
			&um.Mission.Description,
			&um.Mission.PointsReward,
			&um.Mission.XPReward,
			&um.Mission.TargetValue,
			&um.Mission.TierNumber,
			&um.Mission.Category,
			&requirementJSON,
			&um.Mission.ValidFrom,
			&missionValidUntil,
			&um.Mission.IsActive,
			&um.Mission.CreatedAt,
			&um.Mission.UpdatedAt,
		)
		if requirementJSON.Valid {
			um.Mission.RequirementJSON = &requirementJSON.String
		}
		if err != nil {
			return nil, fmt.Errorf("failed to scan user mission: %w", err)
		}

		if completedAt.Valid {
			um.UserMission.CompletedAt = &completedAt.Time
		}
		if claimedAt.Valid {
			um.UserMission.ClaimedAt = &claimedAt.Time
		}
		if missionValidUntil.Valid {
			um.Mission.ValidUntil = &missionValidUntil.Time
		}

		userMissions = append(userMissions, um)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user missions: %w", err)
	}

	return userMissions, nil
}

