package services

import (
	"fmt"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
)

type AchievementService struct {
	repo      *repository.AchievementRepository
	chatRepo  *repository.ChatRepository
	watchRepo *repository.WatchSessionRepository
	userRepo  *repository.UserRepository
}

func NewAchievementService(
	repo *repository.AchievementRepository,
	chatRepo *repository.ChatRepository,
	watchRepo *repository.WatchSessionRepository,
	userRepo *repository.UserRepository,
) *AchievementService {
	return &AchievementService{
		repo:      repo,
		chatRepo:  chatRepo,
		watchRepo: watchRepo,
		userRepo:  userRepo,
	}
}

func (s *AchievementService) SeedDefaults() error {
	defaults := []models.AchievementSeed{
		{Slug: "first_chat", Title: "First Words", Description: "Send your first chat message.", Icon: "💬", Points: 50},
		{Slug: "chatty_50", Title: "Chatty Cyclist", Description: "Send 50 chat messages.", Icon: "🗯️", Points: 150},
		{Slug: "watch_30", Title: "Warm Up", Description: "Watch 30 minutes of streams.", Icon: "⏱️", Points: 100},
		{Slug: "watch_120", Title: "Century Club", Description: "Watch 2 hours of streams.", Icon: "🎥", Points: 200},
		{Slug: "level_5", Title: "Level 5 Achieved", Description: "Reach level 5.", Icon: "🚀", Points: 150},
		{Slug: "level_10", Title: "Level 10 Achieved", Description: "Reach level 10.", Icon: "🏆", Points: 250},
		{Slug: "streak_4_weeks", Title: "Consistency Champ", Description: "Maintain a 4-week streak.", Icon: "🔥", Points: 200},
	}
	return s.repo.EnsureDefaults(defaults)
}

func (s *AchievementService) HandleChatMessage(userID string) {
	if s == nil || s.repo == nil || s.chatRepo == nil {
		return
	}
	count, err := s.chatRepo.CountByUser(userID)
	if err != nil {
		fmt.Printf("AchievementService HandleChatMessage: %v\n", err)
		return
	}
	if count >= 1 {
		s.repo.UnlockBySlug(userID, "first_chat", map[string]interface{}{"messages": count})
	}
	if count >= 50 {
		s.repo.UnlockBySlug(userID, "chatty_50", map[string]interface{}{"messages": count})
	}
}

func (s *AchievementService) HandleWatchTime(userID string) {
	if s == nil || s.repo == nil || s.watchRepo == nil {
		return
	}
	minutes, err := s.watchRepo.GetTotalWatchMinutesByUser(userID)
	if err != nil {
		fmt.Printf("AchievementService HandleWatchTime: %v\n", err)
		return
	}
	if minutes >= 30 {
		s.repo.UnlockBySlug(userID, "watch_30", map[string]interface{}{"minutes": minutes})
	}
	if minutes >= 120 {
		s.repo.UnlockBySlug(userID, "watch_120", map[string]interface{}{"minutes": minutes})
	}
}

func (s *AchievementService) HandleLevelUp(userID string, newLevel int) {
	if s == nil || s.repo == nil {
		return
	}
	if newLevel >= 5 {
		s.repo.UnlockBySlug(userID, "level_5", map[string]interface{}{"level": newLevel})
	}
	if newLevel >= 10 {
		s.repo.UnlockBySlug(userID, "level_10", map[string]interface{}{"level": newLevel})
	}
}

func (s *AchievementService) HandleStreak(userID string) {
	if s == nil || s.repo == nil || s.userRepo == nil {
		return
	}
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		if err != nil {
			fmt.Printf("AchievementService HandleStreak: %v\n", err)
		}
		return
	}
	if user.BestStreakWeeks >= 4 {
		s.repo.UnlockBySlug(userID, "streak_4_weeks", map[string]interface{}{"weeks": user.BestStreakWeeks})
	}
}

func (s *AchievementService) GetUserAchievements(userID string) ([]models.UserAchievement, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("achievement repository not configured")
	}
	return s.repo.GetUserAchievements(userID)
}
