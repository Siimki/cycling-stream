package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/google/uuid"
)

type WatchSessionRepository struct {
	db *sql.DB
}

func NewWatchSessionRepository(db *sql.DB) *WatchSessionRepository {
	return &WatchSessionRepository{db: db}
}

func (r *WatchSessionRepository) Create(session *models.WatchSession) error {
	session.ID = uuid.New().String()
	query := `
		INSERT INTO watch_sessions (id, user_id, race_id, started_at)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`

	err := r.db.QueryRow(
		query,
		session.ID,
		session.UserID,
		session.RaceID,
		session.StartedAt,
	).Scan(&session.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create watch session: %w", err)
	}

	return nil
}

func (r *WatchSessionRepository) GetByID(sessionID string) (*models.WatchSession, error) {
	query := `
		SELECT id, user_id, race_id, started_at, ended_at, duration_seconds, created_at
		FROM watch_sessions
		WHERE id = $1
	`

	var session models.WatchSession
	err := r.db.QueryRow(query, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.RaceID,
		&session.StartedAt,
		&session.EndedAt,
		&session.DurationSeconds,
		&session.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get watch session: %w", err)
	}

	return &session, nil
}

func (r *WatchSessionRepository) EndSession(sessionID string, userID string) error {
	// Get session first to calculate duration
	session, err := r.GetByID(sessionID)
	if err != nil {
		return err
	}

	if session == nil {
		return fmt.Errorf("session not found")
	}

	// Verify session belongs to user
	if session.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	// Check if already ended
	if session.EndedAt != nil {
		return fmt.Errorf("session already ended")
	}

	// Calculate duration
	endedAt := time.Now()
	duration := int(endedAt.Sub(session.StartedAt).Seconds())

	query := `
		UPDATE watch_sessions
		SET ended_at = $1, duration_seconds = $2
		WHERE id = $3
	`

	result, err := r.db.Exec(query, endedAt, duration, sessionID)
	if err != nil {
		return fmt.Errorf("failed to end watch session: %w", err)
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

func (r *WatchSessionRepository) GetStatsByUserAndRace(userID, raceID string) (*models.WatchTimeStats, error) {
	query := `
		SELECT user_id, race_id, session_count, total_seconds, total_minutes, 
		       first_watched, last_watched
		FROM watch_time_aggregated
		WHERE user_id = $1 AND race_id = $2
	`

	var stats models.WatchTimeStats
	err := r.db.QueryRow(query, userID, raceID).Scan(
		&stats.UserID,
		&stats.RaceID,
		&stats.SessionCount,
		&stats.TotalSeconds,
		&stats.TotalMinutes,
		&stats.FirstWatched,
		&stats.LastWatched,
	)

	if err == sql.ErrNoRows {
		// Return empty stats if no sessions found
		return &models.WatchTimeStats{
			UserID:       userID,
			RaceID:       raceID,
			SessionCount: 0,
			TotalSeconds: 0,
			TotalMinutes: 0,
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get watch time stats: %w", err)
	}

	return &stats, nil
}

func (r *WatchSessionRepository) GetActiveSession(userID, raceID string) (*models.WatchSession, error) {
	query := `
		SELECT id, user_id, race_id, started_at, ended_at, duration_seconds, created_at
		FROM watch_sessions
		WHERE user_id = $1 AND race_id = $2 AND ended_at IS NULL
		ORDER BY started_at DESC
		LIMIT 1
	`

	var session models.WatchSession
	err := r.db.QueryRow(query, userID, raceID).Scan(
		&session.ID,
		&session.UserID,
		&session.RaceID,
		&session.StartedAt,
		&session.EndedAt,
		&session.DurationSeconds,
		&session.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active session: %w", err)
	}

	return &session, nil
}

// WatchTimeByRace represents aggregated watch time data per race
type WatchTimeByRace struct {
	RaceID       string  `json:"race_id" db:"race_id"`
	RaceName     string  `json:"race_name" db:"race_name"`
	TotalSeconds int     `json:"total_seconds" db:"total_seconds"`
	TotalMinutes float64 `json:"total_minutes" db:"total_minutes"`
	SessionCount int     `json:"session_count" db:"session_count"`
	UserCount    int     `json:"user_count" db:"user_count"`
	Year         *int    `json:"year,omitempty" db:"year"`
	Month        *int    `json:"month,omitempty" db:"month"`
}

// GetWatchTimeByRace returns aggregated watch time data grouped by race
// Optional year and month filters can be provided
func (r *WatchSessionRepository) GetWatchTimeByRace(year, month *int) ([]WatchTimeByRace, error) {
	var query string
	var args []interface{}

	if year != nil && month != nil {
		// Filter by specific year and month
		query = `
			SELECT 
				r.id as race_id,
				r.name as race_name,
				$1::int as year,
				$2::int as month,
				COALESCE(SUM(ws.duration_seconds), 0)::int as total_seconds,
				COALESCE(SUM(ws.duration_seconds) / 60.0, 0) as total_minutes,
				COUNT(ws.id)::int as session_count,
				COUNT(DISTINCT ws.user_id)::int as user_count
			FROM races r
			LEFT JOIN watch_sessions ws ON r.id = ws.race_id 
				AND ws.duration_seconds IS NOT NULL
				AND EXTRACT(YEAR FROM ws.started_at) = $1
				AND EXTRACT(MONTH FROM ws.started_at) = $2
			GROUP BY r.id, r.name
			ORDER BY total_seconds DESC
		`
		args = []interface{}{*year, *month}
	} else if year != nil {
		// Filter by year only
		query = `
			SELECT 
				r.id as race_id,
				r.name as race_name,
				$1::int as year,
				NULL::int as month,
				COALESCE(SUM(ws.duration_seconds), 0)::int as total_seconds,
				COALESCE(SUM(ws.duration_seconds) / 60.0, 0) as total_minutes,
				COUNT(ws.id)::int as session_count,
				COUNT(DISTINCT ws.user_id)::int as user_count
			FROM races r
			LEFT JOIN watch_sessions ws ON r.id = ws.race_id 
				AND ws.duration_seconds IS NOT NULL
				AND EXTRACT(YEAR FROM ws.started_at) = $1
			GROUP BY r.id, r.name
			ORDER BY total_seconds DESC
		`
		args = []interface{}{*year}
	} else {
		// No filters - aggregate all time
		query = `
			SELECT 
				r.id as race_id,
				r.name as race_name,
				NULL::int as year,
				NULL::int as month,
				COALESCE(SUM(ws.duration_seconds), 0)::int as total_seconds,
				COALESCE(SUM(ws.duration_seconds) / 60.0, 0) as total_minutes,
				COUNT(ws.id)::int as session_count,
				COUNT(DISTINCT ws.user_id)::int as user_count
			FROM races r
			LEFT JOIN watch_sessions ws ON r.id = ws.race_id 
				AND ws.duration_seconds IS NOT NULL
			GROUP BY r.id, r.name
			ORDER BY total_seconds DESC
		`
		args = []interface{}{}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get watch time by race: %w", err)
	}
	defer rows.Close()

	var results []WatchTimeByRace
	for rows.Next() {
		var wtr WatchTimeByRace
		var yearVal sql.NullInt64
		var monthVal sql.NullInt64

		err := rows.Scan(
			&wtr.RaceID,
			&wtr.RaceName,
			&yearVal,
			&monthVal,
			&wtr.TotalSeconds,
			&wtr.TotalMinutes,
			&wtr.SessionCount,
			&wtr.UserCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan watch time by race: %w", err)
		}

		if yearVal.Valid {
			yearInt := int(yearVal.Int64)
			wtr.Year = &yearInt
		}
		if monthVal.Valid {
			monthInt := int(monthVal.Int64)
			wtr.Month = &monthInt
		}

		results = append(results, wtr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating watch time by race: %w", err)
	}

	return results, nil
}

// GetTotalWatchMinutesByUser returns the total watch time in minutes for a user across all races
func (r *WatchSessionRepository) GetTotalWatchMinutesByUser(userID string) (int, error) {
	query := `
		SELECT COALESCE(SUM(duration_seconds), 0) / 60
		FROM watch_sessions
		WHERE user_id = $1 AND duration_seconds IS NOT NULL
	`

	var totalMinutes int
	err := r.db.QueryRow(query, userID).Scan(&totalMinutes)
	if err != nil {
		return 0, fmt.Errorf("failed to get total watch minutes: %w", err)
	}

	return totalMinutes, nil
}

