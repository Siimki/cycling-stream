package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/google/uuid"
)

type ViewerSessionRepository struct {
	db *sql.DB
}

func NewViewerSessionRepository(db *sql.DB) *ViewerSessionRepository {
	return &ViewerSessionRepository{db: db}
}

func (r *ViewerSessionRepository) Create(session *models.ViewerSession) error {
	session.ID = uuid.New().String()
	if session.SessionToken == "" {
		session.SessionToken = uuid.New().String()
	}
	
	query := `
		INSERT INTO viewer_sessions (id, user_id, race_id, session_token, started_at, last_seen_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`

	err := r.db.QueryRow(
		query,
		session.ID,
		session.UserID,
		session.RaceID,
		session.SessionToken,
		session.StartedAt,
		session.LastSeenAt,
		session.IsActive,
	).Scan(&session.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create viewer session: %w", err)
	}

	return nil
}

func (r *ViewerSessionRepository) GetByID(sessionID string) (*models.ViewerSession, error) {
	query := `
		SELECT id, user_id, race_id, session_token, started_at, last_seen_at, ended_at, is_active, created_at
		FROM viewer_sessions
		WHERE id = $1
	`

	var session models.ViewerSession
	var userID sql.NullString
	err := r.db.QueryRow(query, sessionID).Scan(
		&session.ID,
		&userID,
		&session.RaceID,
		&session.SessionToken,
		&session.StartedAt,
		&session.LastSeenAt,
		&session.EndedAt,
		&session.IsActive,
		&session.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get viewer session: %w", err)
	}

	if userID.Valid {
		session.UserID = &userID.String
	}

	return &session, nil
}

func (r *ViewerSessionRepository) GetActiveSessionByToken(raceID, sessionToken string) (*models.ViewerSession, error) {
	query := `
		SELECT id, user_id, race_id, session_token, started_at, last_seen_at, ended_at, is_active, created_at
		FROM viewer_sessions
		WHERE race_id = $1 AND session_token = $2 AND is_active = TRUE
		ORDER BY started_at DESC
		LIMIT 1
	`

	var session models.ViewerSession
	var userID sql.NullString
	err := r.db.QueryRow(query, raceID, sessionToken).Scan(
		&session.ID,
		&userID,
		&session.RaceID,
		&session.SessionToken,
		&session.StartedAt,
		&session.LastSeenAt,
		&session.EndedAt,
		&session.IsActive,
		&session.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active viewer session: %w", err)
	}

	if userID.Valid {
		session.UserID = &userID.String
	}

	return &session, nil
}

func (r *ViewerSessionRepository) UpdateHeartbeat(sessionID string) error {
	query := `
		UPDATE viewer_sessions
		SET last_seen_at = $1
		WHERE id = $2 AND is_active = TRUE
	`

	result, err := r.db.Exec(query, time.Now(), sessionID)
	if err != nil {
		return fmt.Errorf("failed to update heartbeat: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found or not active")
	}

	return nil
}

func (r *ViewerSessionRepository) EndSession(sessionID string) error {
	endedAt := time.Now()
	query := `
		UPDATE viewer_sessions
		SET ended_at = $1, is_active = FALSE
		WHERE id = $2
	`

	result, err := r.db.Exec(query, endedAt, sessionID)
	if err != nil {
		return fmt.Errorf("failed to end viewer session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

func (r *ViewerSessionRepository) GetConcurrentViewers(raceID string) (*models.ConcurrentViewers, error) {
	query := `
		SELECT race_id, concurrent_count, authenticated_count, anonymous_count
		FROM concurrent_viewers
		WHERE race_id = $1
	`

	var viewers models.ConcurrentViewers
	err := r.db.QueryRow(query, raceID).Scan(
		&viewers.RaceID,
		&viewers.ConcurrentCount,
		&viewers.AuthenticatedCount,
		&viewers.AnonymousCount,
	)

	if err == sql.ErrNoRows {
		// Return zero counts if no viewers
		return &models.ConcurrentViewers{
			RaceID:             raceID,
			ConcurrentCount:    0,
			AuthenticatedCount: 0,
			AnonymousCount:     0,
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get concurrent viewers: %w", err)
	}

	return &viewers, nil
}

func (r *ViewerSessionRepository) GetUniqueViewers(raceID string) (*models.UniqueViewers, error) {
	query := `
		SELECT race_id, unique_viewer_count, unique_authenticated_count, unique_anonymous_count
		FROM unique_viewers
		WHERE race_id = $1
	`

	var viewers models.UniqueViewers
	err := r.db.QueryRow(query, raceID).Scan(
		&viewers.RaceID,
		&viewers.UniqueViewerCount,
		&viewers.UniqueAuthenticatedCount,
		&viewers.UniqueAnonymousCount,
	)

	if err == sql.ErrNoRows {
		// Return zero counts if no viewers
		return &models.UniqueViewers{
			RaceID:                   raceID,
			UniqueViewerCount:        0,
			UniqueAuthenticatedCount: 0,
			UniqueAnonymousCount:     0,
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get unique viewers: %w", err)
	}

	return &viewers, nil
}

// CleanupStaleSessions marks sessions as inactive if they haven't been seen in the last 5 minutes
func (r *ViewerSessionRepository) CleanupStaleSessions(timeoutMinutes int) error {
	timeout := time.Now().Add(-time.Duration(timeoutMinutes) * time.Minute)
	query := `
		UPDATE viewer_sessions
		SET is_active = FALSE, ended_at = $1
		WHERE is_active = TRUE AND last_seen_at < $2
	`

	_, err := r.db.Exec(query, time.Now(), timeout)
	if err != nil {
		return fmt.Errorf("failed to cleanup stale sessions: %w", err)
	}

	return nil
}

// GetAllConcurrentViewers returns concurrent viewers for all races
func (r *ViewerSessionRepository) GetAllConcurrentViewers() ([]models.ConcurrentViewers, error) {
	query := `
		SELECT race_id, concurrent_count, authenticated_count, anonymous_count
		FROM concurrent_viewers
		ORDER BY concurrent_count DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all concurrent viewers: %w", err)
	}
	defer rows.Close()

	var viewers []models.ConcurrentViewers
	for rows.Next() {
		var v models.ConcurrentViewers
		if err := rows.Scan(
			&v.RaceID,
			&v.ConcurrentCount,
			&v.AuthenticatedCount,
			&v.AnonymousCount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan concurrent viewers: %w", err)
		}
		viewers = append(viewers, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating concurrent viewers: %w", err)
	}

	return viewers, nil
}

// GetAllUniqueViewers returns unique viewers for all races
func (r *ViewerSessionRepository) GetAllUniqueViewers() ([]models.UniqueViewers, error) {
	query := `
		SELECT race_id, unique_viewer_count, unique_authenticated_count, unique_anonymous_count
		FROM unique_viewers
		ORDER BY unique_viewer_count DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all unique viewers: %w", err)
	}
	defer rows.Close()

	var viewers []models.UniqueViewers
	for rows.Next() {
		var v models.UniqueViewers
		if err := rows.Scan(
			&v.RaceID,
			&v.UniqueViewerCount,
			&v.UniqueAuthenticatedCount,
			&v.UniqueAnonymousCount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan unique viewers: %w", err)
		}
		viewers = append(viewers, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating unique viewers: %w", err)
	}

	return viewers, nil
}

