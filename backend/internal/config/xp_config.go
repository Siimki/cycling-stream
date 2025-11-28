package config

// XPConfig holds all configuration for the XP, points, and leveling system
type XPConfig struct {
	Leveling LevelingConfig
	Awards   AwardsConfig
	Points   PointsConfig
	Weekly   WeeklyConfig
}

// LevelingConfig defines how levels are calculated from XP
type LevelingConfig struct {
	// BaseXP is the XP needed to reach level 2 (default: 100)
	// Level 1: 0 to (BaseXP - 1) XP
	BaseXP int
	// IncrementPerLevel is the additional XP needed per level beyond level 2 (default: 20)
	// Level N (N > 1) requires: BaseXP + (N-2) * IncrementPerLevel XP
	IncrementPerLevel int
}

// AwardsConfig defines XP award rates and caps for different activities
type AwardsConfig struct {
	WatchTime   WatchTimeAwardConfig
	ChatMessage ChatMessageAwardConfig
	WeeklyGoal  WeeklyGoalAwardConfig
	Mission     MissionAwardConfig
	Prediction  PredictionAwardConfig
}

// WatchTimeAwardConfig defines XP awards for watching races
type WatchTimeAwardConfig struct {
	// XPPerMinute is the XP awarded per minute watched (default: 0, meaning 1 XP per 2 minutes)
	// Note: We use integer division, so 0 means 1 XP per 2 minutes (minutes/2)
	// If you want 1 XP per minute, set to 1
	XPPerMinute int
	// CapPerRace is the maximum XP that can be earned per race from watching (default: 200)
	CapPerRace int
}

// ChatMessageAwardConfig defines XP awards for chat messages
type ChatMessageAwardConfig struct {
	// XPPerMessage is the XP awarded per chat message (default: 2)
	XPPerMessage int
	// CapPerRace is the maximum XP that can be earned per race from chat (default: 50)
	CapPerRace int
	// SpamGuardSeconds is the minimum seconds between XP awards for chat (default: 10)
	SpamGuardSeconds int
}

// WeeklyGoalAwardConfig defines rewards for completing weekly goals
type WeeklyGoalAwardConfig struct {
	// XP is the XP reward for completing weekly goal (default: 150)
	XP int
	// Points is the points reward for completing weekly goal (default: 200)
	Points int
}

// MissionAwardConfig defines default rewards for missions (can be overridden per mission)
type MissionAwardConfig struct {
	// DefaultXP is the default XP reward if mission doesn't specify (default: 0)
	DefaultXP int
	// DefaultPoints is the default points reward if mission doesn't specify (default: 0)
	DefaultPoints int
}

// PredictionAwardConfig defines XP awards for predictions
type PredictionAwardConfig struct {
	// XPForPlacing is the XP awarded when placing a prediction (default: 5)
	XPForPlacing int
	// XPForWinning is the XP awarded when winning a prediction (default: 15)
	XPForWinning int
}

// PointsConfig defines points award rates
type PointsConfig struct {
	// PerBlockSeconds is the time interval for points blocks (default: 10)
	PerBlockSeconds int
	// PerBlock is the points awarded per block (default: 10)
	PerBlock int
}

// WeeklyConfig defines weekly goal thresholds
type WeeklyConfig struct {
	// WatchMinutesGoal is the watch time goal in minutes (default: 30)
	WatchMinutesGoal int
	// ChatMessagesGoal is the chat messages goal (default: 3)
	ChatMessagesGoal int
}

// LoadXPConfig loads XP configuration from environment variables with defaults
func LoadXPConfig() *XPConfig {
	return &XPConfig{
		Leveling: LevelingConfig{
			BaseXP:           getEnvAsInt("XP_LEVELING_BASE_XP", 100),
			IncrementPerLevel: getEnvAsInt("XP_LEVELING_INCREMENT", 20),
		},
		Awards: AwardsConfig{
			WatchTime: WatchTimeAwardConfig{
				XPPerMinute: getEnvAsInt("XP_AWARD_WATCH_XP_PER_MINUTE", 0), // 0 means 1 per 2 min
				CapPerRace:  getEnvAsInt("XP_AWARD_WATCH_CAP_PER_RACE", 200),
			},
			ChatMessage: ChatMessageAwardConfig{
				XPPerMessage:      getEnvAsInt("XP_AWARD_CHAT_XP_PER_MESSAGE", 2),
				CapPerRace:        getEnvAsInt("XP_AWARD_CHAT_CAP_PER_RACE", 50),
				SpamGuardSeconds:  getEnvAsInt("XP_AWARD_CHAT_SPAM_GUARD_SECONDS", 10),
			},
			WeeklyGoal: WeeklyGoalAwardConfig{
				XP:     getEnvAsInt("XP_AWARD_WEEKLY_GOAL_XP", 150),
				Points: getEnvAsInt("XP_AWARD_WEEKLY_GOAL_POINTS", 200),
			},
			Mission: MissionAwardConfig{
				DefaultXP:     getEnvAsInt("XP_AWARD_MISSION_DEFAULT_XP", 0),
				DefaultPoints: getEnvAsInt("XP_AWARD_MISSION_DEFAULT_POINTS", 0),
			},
			Prediction: PredictionAwardConfig{
				XPForPlacing: getEnvAsInt("XP_AWARD_PREDICTION_PLACING", 5),
				XPForWinning: getEnvAsInt("XP_AWARD_PREDICTION_WINNING", 15),
			},
		},
		Points: PointsConfig{
			PerBlockSeconds: getEnvAsInt("POINTS_PER_BLOCK_SECONDS", 10),
			PerBlock:        getEnvAsInt("POINTS_PER_BLOCK", 10),
		},
		Weekly: WeeklyConfig{
			WatchMinutesGoal:  getEnvAsInt("WEEKLY_GOAL_WATCH_MINUTES", 30),
			ChatMessagesGoal:  getEnvAsInt("WEEKLY_GOAL_CHAT_MESSAGES", 3),
		},
	}
}

