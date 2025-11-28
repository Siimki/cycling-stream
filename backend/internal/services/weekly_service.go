package services

import (
	"fmt"
	"time"

	"github.com/cyclingstream/backend/internal/config"
	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
)

type WeeklyService struct {
	weeklyRepo  *repository.WeeklyRepository
	streakRepo  *repository.StreakRepository
	userRepo    *repository.UserRepository
	xpService   *XPService
	xpConfig    *config.XPConfig
}

func NewWeeklyService(
	weeklyRepo *repository.WeeklyRepository,
	streakRepo *repository.StreakRepository,
	userRepo *repository.UserRepository,
	xpService *XPService,
	xpConfig *config.XPConfig,
) *WeeklyService {
	return &WeeklyService{
		weeklyRepo: weeklyRepo,
		streakRepo: streakRepo,
		userRepo:   userRepo,
		xpService:  xpService,
		xpConfig:   xpConfig,
	}
}

// GetCurrentWeekNumber returns the current ISO week number in format YYYY-WW
func (s *WeeklyService) GetCurrentWeekNumber() string {
	now := time.Now()
	year, week := now.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

// GetWeeklyProgress returns the current week's progress for a user
func (s *WeeklyService) GetWeeklyProgress(userID string) (*models.WeeklyGoalProgress, error) {
	weekNumber := s.GetCurrentWeekNumber()

	// Get weekly stats
	stats, err := s.weeklyRepo.GetCurrentWeekStats(userID, weekNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly stats: %w", err)
	}

	// Get streak
	streak, err := s.streakRepo.GetStreak(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get streak: %w", err)
	}

	// Get user for best streak
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	watchMinutesGoal := 30
	chatMessagesGoal := 3
	if s.xpConfig != nil {
		watchMinutesGoal = s.xpConfig.Weekly.WatchMinutesGoal
		chatMessagesGoal = s.xpConfig.Weekly.ChatMessagesGoal
	}

	xpReward := 150
	pointsReward := 200
	if s.xpConfig != nil {
		xpReward = s.xpConfig.Awards.WeeklyGoal.XP
		pointsReward = s.xpConfig.Awards.WeeklyGoal.Points
	}

	progress := &models.WeeklyGoalProgress{
		WeekNumber:         weekNumber,
		WatchMinutes:       stats.WatchMinutes,
		WatchMinutesGoal:   watchMinutesGoal,
		ChatMessages:       stats.ChatMessages,
		ChatMessagesGoal:   chatMessagesGoal,
		GoalCompleted:      stats.WeeklyGoalCompleted,
		RewardClaimed:      stats.WeeklyRewardClaimedAt != nil,
		RewardXP:           xpReward,
		RewardPoints:       pointsReward,
		CurrentStreakWeeks: streak.CurrentStreakWeeks,
		BestStreakWeeks:    user.BestStreakWeeks,
	}

	return progress, nil
}

// CheckAndCompleteWeeklyGoal checks if the weekly goal is completed and processes it if so
func (s *WeeklyService) CheckAndCompleteWeeklyGoal(userID string) error {
	weekNumber := s.GetCurrentWeekNumber()

	// Get weekly stats
	stats, err := s.weeklyRepo.GetCurrentWeekStats(userID, weekNumber)
	if err != nil {
		return fmt.Errorf("failed to get weekly stats: %w", err)
	}

	// Check if already completed
	if stats.WeeklyGoalCompleted {
		return nil
	}

	// Check if both thresholds are met (using config values)
	watchMinutesGoal := 30
	chatMessagesGoal := 3
	if s.xpConfig != nil {
		watchMinutesGoal = s.xpConfig.Weekly.WatchMinutesGoal
		chatMessagesGoal = s.xpConfig.Weekly.ChatMessagesGoal
	}

	if stats.WatchMinutes >= watchMinutesGoal && stats.ChatMessages >= chatMessagesGoal {
		// Complete the goal
		return s.ProcessWeeklyGoalCompletion(userID, weekNumber)
	}

	return nil
}

// ProcessWeeklyGoalCompletion awards rewards and updates streak when weekly goal is completed
func (s *WeeklyService) ProcessWeeklyGoalCompletion(userID, weekNumber string) error {
	// Mark goal as completed
	if err := s.weeklyRepo.CompleteWeeklyGoal(userID, weekNumber); err != nil {
		return fmt.Errorf("failed to complete weekly goal: %w", err)
	}

	// Award rewards based on config
	xpReward := 150
	pointsReward := 200
	if s.xpConfig != nil {
		xpReward = s.xpConfig.Awards.WeeklyGoal.XP
		pointsReward = s.xpConfig.Awards.WeeklyGoal.Points
	}

	if s.xpService != nil && xpReward > 0 {
		if err := s.xpService.AwardXP(userID, xpReward, "weekly_goal"); err != nil {
			// Log error but don't fail
			fmt.Printf("Failed to award XP for weekly goal: %v\n", err)
		}
	}

	if pointsReward > 0 {
		if err := s.userRepo.AddPoints(userID, pointsReward); err != nil {
			// Log error but don't fail
			fmt.Printf("Failed to award points for weekly goal: %v\n", err)
		}
	}

	// Update streak
	if err := s.streakRepo.UpdateStreak(userID, weekNumber, true); err != nil {
		return fmt.Errorf("failed to update streak: %w", err)
	}

	// Update best streak if current streak is higher
	streak, err := s.streakRepo.GetStreak(userID)
	if err != nil {
		return fmt.Errorf("failed to get streak: %w", err)
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if streak.CurrentStreakWeeks > user.BestStreakWeeks {
		// Update best streak in users table
		if err := s.userRepo.UpdateBestStreak(userID, streak.CurrentStreakWeeks); err != nil {
			// Log error but don't fail
			fmt.Printf("Failed to update best streak: %v\n", err)
		}
	}

	return nil
}

// UpdateStreakOnWeekEnd is called at week boundary to update streak based on previous week completion
// This should be called by a background job or cron
func (s *WeeklyService) UpdateStreakOnWeekEnd(userID, previousWeekNumber string) error {
	// Get stats for previous week
	stats, err := s.weeklyRepo.GetCurrentWeekStats(userID, previousWeekNumber)
	if err != nil {
		return fmt.Errorf("failed to get weekly stats: %w", err)
	}

	// Update streak based on whether goal was completed
	if err := s.streakRepo.UpdateStreak(userID, previousWeekNumber, stats.WeeklyGoalCompleted); err != nil {
		return fmt.Errorf("failed to update streak: %w", err)
	}

	return nil
}

