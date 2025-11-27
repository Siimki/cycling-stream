package repository

import (
	"database/sql"
	"fmt"

	"github.com/cyclingstream/backend/internal/models"
)

type WatchHistoryRepository struct {
	db *sql.DB
}

func NewWatchHistoryRepository(db *sql.DB) *WatchHistoryRepository {
	return &WatchHistoryRepository{db: db}
}

func (r *WatchHistoryRepository) GetByUserID(userID string, limit, offset int) ([]*models.WatchHistoryEntry, error) {
	query := `
		SELECT user_id, race_id, race_name, race_category, race_start_date,
		       session_count, total_seconds, total_minutes, first_watched, last_watched, likely_completed
		FROM user_watch_history
		WHERE user_id = $1
		ORDER BY last_watched DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get watch history: %w", err)
	}
	defer rows.Close()

	var entries []*models.WatchHistoryEntry
	for rows.Next() {
		var entry models.WatchHistoryEntry
		var raceCategory sql.NullString
		var raceStartDate sql.NullTime

		err := rows.Scan(
			&entry.UserID,
			&entry.RaceID,
			&entry.RaceName,
			&raceCategory,
			&raceStartDate,
			&entry.SessionCount,
			&entry.TotalSeconds,
			&entry.TotalMinutes,
			&entry.FirstWatched,
			&entry.LastWatched,
			&entry.LikelyCompleted,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan watch history entry: %w", err)
		}

		if raceCategory.Valid {
			entry.RaceCategory = &raceCategory.String
		}
		if raceStartDate.Valid {
			entry.RaceStartDate = &raceStartDate.Time
		}

		entries = append(entries, &entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating watch history entries: %w", err)
	}

	return entries, nil
}

func (r *WatchHistoryRepository) GetCountByUserID(userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM user_watch_history
		WHERE user_id = $1
	`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get watch history count: %w", err)
	}

	return count, nil
}

