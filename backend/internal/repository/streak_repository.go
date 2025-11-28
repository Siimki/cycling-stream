package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cyclingstream/backend/internal/models"
)

type StreakRepository struct {
	db *sql.DB
}

func NewStreakRepository(db *sql.DB) *StreakRepository {
	return &StreakRepository{db: db}
}

// GetOrCreateStreak gets or creates a streak record for a user
func (r *StreakRepository) GetOrCreateStreak(userID string) (*models.UserStreak, error) {
	query := `
		SELECT user_id, current_streak_weeks, last_completed_week_number, updated_at
		FROM user_streaks
		WHERE user_id = $1
	`

	var streak models.UserStreak
	var lastCompletedWeek sql.NullString
	err := r.db.QueryRow(query, userID).Scan(
		&streak.UserID,
		&streak.CurrentStreakWeeks,
		&lastCompletedWeek,
		&streak.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Create new streak
		streak.UserID = userID
		streak.CurrentStreakWeeks = 0
		streak.LastCompletedWeekNumber = nil
		streak.UpdatedAt = time.Now()

		insertQuery := `
			INSERT INTO user_streaks (user_id, current_streak_weeks, last_completed_week_number)
			VALUES ($1, $2, $3)
			RETURNING updated_at
		`

		err = r.db.QueryRow(
			insertQuery,
			streak.UserID,
			streak.CurrentStreakWeeks,
			streak.LastCompletedWeekNumber,
		).Scan(&streak.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to create streak: %w", err)
		}

		return &streak, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get streak: %w", err)
	}

	if lastCompletedWeek.Valid {
		streak.LastCompletedWeekNumber = &lastCompletedWeek.String
	}

	return &streak, nil
}

// UpdateStreak updates the streak based on whether the week was completed
func (r *StreakRepository) UpdateStreak(userID, weekNumber string, completed bool) error {
	streak, err := r.GetOrCreateStreak(userID)
	if err != nil {
		return fmt.Errorf("failed to get streak: %w", err)
	}

	if completed {
		// Check if this is a consecutive week
		if streak.LastCompletedWeekNumber == nil || *streak.LastCompletedWeekNumber != weekNumber {
			// Increment streak
			newStreak := streak.CurrentStreakWeeks + 1
			query := `
				UPDATE user_streaks
				SET current_streak_weeks = $1, last_completed_week_number = $2, updated_at = CURRENT_TIMESTAMP
				WHERE user_id = $3
			`

			result, err := r.db.Exec(query, newStreak, weekNumber, userID)
			if err != nil {
				return fmt.Errorf("failed to update streak: %w", err)
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return fmt.Errorf("failed to check rows affected: %w", err)
			}

			if rowsAffected == 0 {
				return fmt.Errorf("streak not found")
			}
		}
	} else {
		// Reset streak to 0
		query := `
			UPDATE user_streaks
			SET current_streak_weeks = 0, updated_at = CURRENT_TIMESTAMP
			WHERE user_id = $1
		`

		result, err := r.db.Exec(query, userID)
		if err != nil {
			return fmt.Errorf("failed to reset streak: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to check rows affected: %w", err)
		}

		if rowsAffected == 0 {
			// Create if doesn't exist (shouldn't happen, but handle it)
			_, err := r.GetOrCreateStreak(userID)
			if err != nil {
				return fmt.Errorf("failed to create streak: %w", err)
			}
		}
	}

	return nil
}

// GetStreak gets the current streak for a user
func (r *StreakRepository) GetStreak(userID string) (*models.UserStreak, error) {
	return r.GetOrCreateStreak(userID)
}


