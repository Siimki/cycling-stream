package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/google/uuid"
)

type WeeklyRepository struct {
	db *sql.DB
}

func NewWeeklyRepository(db *sql.DB) *WeeklyRepository {
	return &WeeklyRepository{db: db}
}

// GetOrCreateWeeklyStats gets or creates weekly stats for a user and week
func (r *WeeklyRepository) GetOrCreateWeeklyStats(userID, weekNumber string) (*models.UserWeeklyStats, error) {
	query := `
		SELECT id, user_id, week_number, watch_minutes, chat_messages, 
		       weekly_goal_completed, weekly_reward_claimed_at, created_at, updated_at
		FROM user_weekly_stats
		WHERE user_id = $1 AND week_number = $2
	`

	var stats models.UserWeeklyStats
	err := r.db.QueryRow(query, userID, weekNumber).Scan(
		&stats.ID,
		&stats.UserID,
		&stats.WeekNumber,
		&stats.WatchMinutes,
		&stats.ChatMessages,
		&stats.WeeklyGoalCompleted,
		&stats.WeeklyRewardClaimedAt,
		&stats.CreatedAt,
		&stats.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Create new stats
		stats.ID = uuid.New().String()
		stats.UserID = userID
		stats.WeekNumber = weekNumber
		stats.WatchMinutes = 0
		stats.ChatMessages = 0
		stats.WeeklyGoalCompleted = false
		stats.CreatedAt = time.Now()
		stats.UpdatedAt = time.Now()

		insertQuery := `
			INSERT INTO user_weekly_stats (id, user_id, week_number, watch_minutes, chat_messages, weekly_goal_completed)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING created_at, updated_at
		`

		err = r.db.QueryRow(
			insertQuery,
			stats.ID,
			stats.UserID,
			stats.WeekNumber,
			stats.WatchMinutes,
			stats.ChatMessages,
			stats.WeeklyGoalCompleted,
		).Scan(&stats.CreatedAt, &stats.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to create weekly stats: %w", err)
		}

		return &stats, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get weekly stats: %w", err)
	}

	return &stats, nil
}

// UpdateWatchMinutes increments watch minutes for a user's weekly stats
func (r *WeeklyRepository) UpdateWatchMinutes(userID, weekNumber string, minutes int) error {
	query := `
		UPDATE user_weekly_stats
		SET watch_minutes = watch_minutes + $1, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $2 AND week_number = $3
	`

	result, err := r.db.Exec(query, minutes, userID, weekNumber)
	if err != nil {
		return fmt.Errorf("failed to update watch minutes: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		// Create if doesn't exist
		_, err := r.GetOrCreateWeeklyStats(userID, weekNumber)
		if err != nil {
			return fmt.Errorf("failed to create weekly stats: %w", err)
		}
		// Retry update
		result, err := r.db.Exec(query, minutes, userID, weekNumber)
		if err != nil {
			return fmt.Errorf("failed to update watch minutes after create: %w", err)
		}
		rowsAffected, _ = result.RowsAffected()
		if rowsAffected == 0 {
			return fmt.Errorf("failed to update watch minutes")
		}
	}

	return nil
}

// IncrementChatMessages increments chat messages for a user's weekly stats
func (r *WeeklyRepository) IncrementChatMessages(userID, weekNumber string) error {
	query := `
		UPDATE user_weekly_stats
		SET chat_messages = chat_messages + 1, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND week_number = $2
	`

	result, err := r.db.Exec(query, userID, weekNumber)
	if err != nil {
		return fmt.Errorf("failed to increment chat messages: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		// Create if doesn't exist
		_, err := r.GetOrCreateWeeklyStats(userID, weekNumber)
		if err != nil {
			return fmt.Errorf("failed to create weekly stats: %w", err)
		}
		// Retry increment
		result, err := r.db.Exec(query, userID, weekNumber)
		if err != nil {
			return fmt.Errorf("failed to increment chat messages after create: %w", err)
		}
		rowsAffected, _ = result.RowsAffected()
		if rowsAffected == 0 {
			return fmt.Errorf("failed to increment chat messages")
		}
	}

	return nil
}

// CompleteWeeklyGoal marks the weekly goal as completed
func (r *WeeklyRepository) CompleteWeeklyGoal(userID, weekNumber string) error {
	query := `
		UPDATE user_weekly_stats
		SET weekly_goal_completed = true, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND week_number = $2
	`

	result, err := r.db.Exec(query, userID, weekNumber)
	if err != nil {
		return fmt.Errorf("failed to complete weekly goal: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("weekly stats not found")
	}

	return nil
}

// ClaimWeeklyReward marks the weekly reward as claimed
func (r *WeeklyRepository) ClaimWeeklyReward(userID, weekNumber string) error {
	query := `
		UPDATE user_weekly_stats
		SET weekly_reward_claimed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND week_number = $2
	`

	result, err := r.db.Exec(query, userID, weekNumber)
	if err != nil {
		return fmt.Errorf("failed to claim weekly reward: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("weekly stats not found")
	}

	return nil
}

// GetCurrentWeekStats gets the current week's stats for a user
func (r *WeeklyRepository) GetCurrentWeekStats(userID, weekNumber string) (*models.UserWeeklyStats, error) {
	return r.GetOrCreateWeeklyStats(userID, weekNumber)
}


