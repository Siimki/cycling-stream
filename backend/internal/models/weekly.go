package models

import "time"

// UserWeeklyStats represents weekly goal progress for a user
type UserWeeklyStats struct {
	ID                    string     `json:"id" db:"id"`
	UserID                string     `json:"user_id" db:"user_id"`
	WeekNumber            string     `json:"week_number" db:"week_number"` // ISO format: YYYY-WW
	WatchMinutes          int        `json:"watch_minutes" db:"watch_minutes"`
	ChatMessages          int        `json:"chat_messages" db:"chat_messages"`
	WeeklyGoalCompleted   bool       `json:"weekly_goal_completed" db:"weekly_goal_completed"`
	WeeklyRewardClaimedAt *time.Time `json:"weekly_reward_claimed_at,omitempty" db:"weekly_reward_claimed_at"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
}

// UserStreak represents the current streak state for a user
type UserStreak struct {
	UserID                string     `json:"user_id" db:"user_id"`
	CurrentStreakWeeks    int        `json:"current_streak_weeks" db:"current_streak_weeks"`
	LastCompletedWeekNumber *string  `json:"last_completed_week_number,omitempty" db:"last_completed_week_number"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
}

// WeeklyGoalProgress represents the current week's progress for API responses
type WeeklyGoalProgress struct {
	WeekNumber          string `json:"week_number"`
	WatchMinutes        int    `json:"watch_minutes"`
	WatchMinutesGoal    int    `json:"watch_minutes_goal"` // 30
	ChatMessages        int    `json:"chat_messages"`
	ChatMessagesGoal    int    `json:"chat_messages_goal"` // 3
	GoalCompleted       bool   `json:"goal_completed"`
	RewardClaimed       bool   `json:"reward_claimed"`
	RewardXP            int    `json:"reward_xp,omitempty"`     // XP reward for completing weekly goal
	RewardPoints        int    `json:"reward_points,omitempty"` // Points reward for completing weekly goal
	CurrentStreakWeeks  int    `json:"current_streak_weeks"`
	BestStreakWeeks     int    `json:"best_streak_weeks"`
}


