package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/google/uuid"
)

type EntitlementRepository struct {
	db *sql.DB
}

func NewEntitlementRepository(db *sql.DB) *EntitlementRepository {
	return &EntitlementRepository{db: db}
}

func (r *EntitlementRepository) GetByUserAndRace(userID, raceID string) (*models.Entitlement, error) {
	query := `
		SELECT id, user_id, race_id, type, expires_at, created_at
		FROM entitlements
		WHERE user_id = $1 AND race_id = $2
	`

	var entitlement models.Entitlement
	err := r.db.QueryRow(query, userID, raceID).Scan(
		&entitlement.ID,
		&entitlement.UserID,
		&entitlement.RaceID,
		&entitlement.Type,
		&entitlement.ExpiresAt,
		&entitlement.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get entitlement: %w", err)
	}

	// Check if expired
	if entitlement.ExpiresAt != nil && entitlement.ExpiresAt.Before(time.Now()) {
		return nil, nil
	}

	return &entitlement, nil
}

func (r *EntitlementRepository) Create(entitlement *models.Entitlement) error {
	entitlement.ID = uuid.New().String()
	query := `
		INSERT INTO entitlements (id, user_id, race_id, type, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, race_id) DO UPDATE
		SET type = $4, expires_at = $5
		RETURNING created_at
	`

	err := r.db.QueryRow(
		query,
		entitlement.ID,
		entitlement.UserID,
		entitlement.RaceID,
		entitlement.Type,
		entitlement.ExpiresAt,
	).Scan(&entitlement.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create entitlement: %w", err)
	}

	return nil
}

func (r *EntitlementRepository) HasAccess(userID, raceID string) (bool, error) {
	// First check if race is free
	var isFree bool
	err := r.db.QueryRow("SELECT is_free FROM races WHERE id = $1", raceID).Scan(&isFree)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check race: %w", err)
	}

	if isFree {
		return true, nil
	}

	// Check entitlement
	entitlement, err := r.GetByUserAndRace(userID, raceID)
	if err != nil {
		return false, err
	}

	return entitlement != nil, nil
}

// HasActiveSubscription returns true if the user has an active recurring entitlement.
func (r *EntitlementRepository) HasActiveSubscription(userID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM entitlements
			WHERE user_id = $1
				AND type IN ('subscription', 'season_pass')
				AND (expires_at IS NULL OR expires_at > NOW())
		)
	`

	var exists bool
	if err := r.db.QueryRow(query, userID).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check subscription entitlements: %w", err)
	}

	return exists, nil
}
