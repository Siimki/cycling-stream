package models

import "time"

// MissionType represents the type of mission
type MissionType string

const (
	MissionTypeWatchTime    MissionType = "watch_time"
	MissionTypeChatMessage  MissionType = "chat_message"
	MissionTypeWatchRace    MissionType = "watch_race"
	MissionTypeFollowSeries MissionType = "follow_series"
	MissionTypeStreak       MissionType = "streak"
	MissionTypePredictWinner MissionType = "predict_winner"
)

// Mission represents a mission that users can complete
type Mission struct {
	ID             string     `json:"id" db:"id"`
	MissionType    MissionType `json:"mission_type" db:"mission_type"`
	Title          string     `json:"title" db:"title"`
	Description    *string    `json:"description,omitempty" db:"description"`
	PointsReward   int        `json:"points_reward" db:"points_reward"`
	XPReward       int        `json:"xp_reward" db:"xp_reward"`
	TargetValue    int        `json:"target_value" db:"target_value"`
	TierNumber     int        `json:"tier_number" db:"tier_number"`
	Category       string     `json:"category" db:"category"` // 'career' or 'weekly'
	RequirementJSON *string   `json:"requirement_json,omitempty" db:"requirement_json"` // JSONB stored as string
	ValidFrom      time.Time  `json:"valid_from" db:"valid_from"`
	ValidUntil     *time.Time `json:"valid_until,omitempty" db:"valid_until"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// UserMission represents a user's progress on a mission
type UserMission struct {
	ID          string     `json:"id" db:"id"`
	UserID      string     `json:"user_id" db:"user_id"`
	MissionID   string     `json:"mission_id" db:"mission_id"`
	Progress    int        `json:"progress" db:"progress"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	ClaimedAt   *time.Time `json:"claimed_at,omitempty" db:"claimed_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// UserMissionWithDetails includes mission details along with user progress
type UserMissionWithDetails struct {
	UserMission
	Mission
}

// IsCompleted returns true if the mission is completed
func (um *UserMission) IsCompleted() bool {
	return um.CompletedAt != nil
}

// IsClaimed returns true if the mission reward has been claimed
func (um *UserMission) IsClaimed() bool {
	return um.ClaimedAt != nil
}

// CanClaim returns true if the mission is completed but not yet claimed
func (um *UserMission) CanClaim() bool {
	return um.IsCompleted() && !um.IsClaimed()
}

