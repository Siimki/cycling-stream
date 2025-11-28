package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/cyclingstream/backend/internal/config"
	"github.com/cyclingstream/backend/internal/models"
)

type MissionTriggers struct {
	missionService     *MissionService
	xpService          *XPService
	weeklyService      *WeeklyService
	achievementService *AchievementService
	xpConfig           *config.XPConfig
	// Track XP awarded per user per race (in-memory, resets on restart)
	// Key: userID:raceID, Value: XP awarded
	xpPerRace map[string]int
	xpMutex   sync.RWMutex
	// Track last chat message XP timestamp per user per race for spam guard
	// Key: userID:raceID, Value: timestamp
	lastChatXP map[string]time.Time
	chatMutex  sync.RWMutex
}

func NewMissionTriggers(missionService *MissionService, xpService *XPService, weeklyService *WeeklyService, achievementService *AchievementService, xpConfig *config.XPConfig) *MissionTriggers {
	return &MissionTriggers{
		missionService:     missionService,
		xpService:          xpService,
		weeklyService:      weeklyService,
		achievementService: achievementService,
		xpConfig:           xpConfig,
		xpPerRace:          make(map[string]int),
		lastChatXP:         make(map[string]time.Time),
	}
}

// OnWatchTime tracks watch time for missions and awards XP
// durationSeconds is the duration in seconds that was watched
// raceID is required for XP cap tracking
// isLive indicates if the race stream is currently live (for weekly stats)
func (t *MissionTriggers) OnWatchTime(userID, raceID string, durationSeconds int, isLive bool) error {
	if t.achievementService != nil {
		go t.achievementService.HandleWatchTime(userID)
	}
	// Convert seconds to minutes for watch_time missions
	// We increment by minutes watched
	minutesWatched := durationSeconds / 60
	if minutesWatched > 0 {
		if err := t.missionService.UpdateMissionProgress(userID, models.MissionTypeWatchTime, minutesWatched); err != nil {
			return fmt.Errorf("failed to update watch time mission progress: %w", err)
		}
	}

	// Award XP based on config: XP per minute watched, cap per race
	if t.xpService != nil && raceID != "" && t.xpConfig != nil {
		// Calculate XP to award based on config
		// If XPPerMinute is 0, it means 1 XP per 2 minutes (legacy behavior)
		var xpToAward int
		if t.xpConfig.Awards.WatchTime.XPPerMinute == 0 {
			// Legacy: 1 XP per 2 minutes (round down)
			xpToAward = minutesWatched / 2
		} else {
			// New: XP per minute
			xpToAward = minutesWatched * t.xpConfig.Awards.WatchTime.XPPerMinute
		}

		if xpToAward > 0 {
			// Check current XP for this race
			key := fmt.Sprintf("%s:%s", userID, raceID)
			t.xpMutex.Lock()
			currentXP := t.xpPerRace[key]
			remainingCap := t.xpConfig.Awards.WatchTime.CapPerRace - currentXP
			if remainingCap > 0 {
				// Award XP up to the cap
				xpToAwardActual := xpToAward
				if xpToAwardActual > remainingCap {
					xpToAwardActual = remainingCap
				}
				t.xpPerRace[key] = currentXP + xpToAwardActual
				t.xpMutex.Unlock()

				// Award XP
				if err := t.xpService.AwardXP(userID, xpToAwardActual, "watch_time"); err != nil {
					// Log error but don't fail
					return fmt.Errorf("failed to award XP: %w", err)
				}
			} else {
				t.xpMutex.Unlock()
			}
		}
	}

	// Update weekly stats (only for live races)
	if t.weeklyService != nil && isLive && raceID != "" {
		minutesWatched := durationSeconds / 60
		if minutesWatched > 0 {
			weekNumber := t.weeklyService.GetCurrentWeekNumber()
			if err := t.weeklyService.weeklyRepo.UpdateWatchMinutes(userID, weekNumber, minutesWatched); err != nil {
				// Log error but don't fail
				fmt.Printf("Failed to update weekly watch minutes: %v\n", err)
			} else {
				// Check if weekly goal is completed
				if err := t.weeklyService.CheckAndCompleteWeeklyGoal(userID); err != nil {
					// Log error but don't fail
					fmt.Printf("Failed to check weekly goal: %v\n", err)
				}
			}
		}
	}

	return nil
}

// OnChatMessage tracks chat messages for missions and awards XP
// raceID is required for XP cap tracking and spam guard
// isLive indicates if the race stream is currently live (for weekly stats)
func (t *MissionTriggers) OnChatMessage(userID, raceID string, isLive bool) error {
	if t.achievementService != nil {
		go t.achievementService.HandleChatMessage(userID)
	}
	if err := t.missionService.UpdateMissionProgress(userID, models.MissionTypeChatMessage, 1); err != nil {
		return fmt.Errorf("failed to update chat message mission progress: %w", err)
	}

	// Award XP based on config: XP per message, with spam guard and cap per race
	if t.xpService != nil && raceID != "" && t.xpConfig != nil {
		key := fmt.Sprintf("%s:%s", userID, raceID)
		now := time.Now()

		// Check spam guard: max 1 message XP per configured seconds
		spamGuardDuration := time.Duration(t.xpConfig.Awards.ChatMessage.SpamGuardSeconds) * time.Second
		t.chatMutex.Lock()
		lastXP, exists := t.lastChatXP[key]
		if exists && now.Sub(lastXP) < spamGuardDuration {
			t.chatMutex.Unlock()
			// Spam guard: don't award XP
			return nil
		}
		t.lastChatXP[key] = now
		t.chatMutex.Unlock()

		// Check XP cap for chat
		// We'll use a separate tracking key for chat XP
		chatXPKey := fmt.Sprintf("chat:%s:%s", userID, raceID)
		t.xpMutex.Lock()
		currentChatXP := t.xpPerRace[chatXPKey]
		remainingChatCap := t.xpConfig.Awards.ChatMessage.CapPerRace - currentChatXP
		if remainingChatCap > 0 {
			// Award XP per message (or remaining cap if less)
			xpToAward := t.xpConfig.Awards.ChatMessage.XPPerMessage
			if xpToAward > remainingChatCap {
				xpToAward = remainingChatCap
			}
			t.xpPerRace[chatXPKey] = currentChatXP + xpToAward
			t.xpMutex.Unlock()

			// Award XP
			if err := t.xpService.AwardXP(userID, xpToAward, "chat_message"); err != nil {
				// Log error but don't fail
				return fmt.Errorf("failed to award XP: %w", err)
			}
		} else {
			t.xpMutex.Unlock()
		}
	}

	// Update weekly stats (only for live races)
	if t.weeklyService != nil && isLive && raceID != "" {
		weekNumber := t.weeklyService.GetCurrentWeekNumber()
		if err := t.weeklyService.weeklyRepo.IncrementChatMessages(userID, weekNumber); err != nil {
			// Log error but don't fail
			fmt.Printf("Failed to update weekly chat messages: %v\n", err)
		} else {
			// Check if weekly goal is completed
			if err := t.weeklyService.CheckAndCompleteWeeklyGoal(userID); err != nil {
				// Log error but don't fail
				fmt.Printf("Failed to check weekly goal: %v\n", err)
			}
		}
	}

	return nil
}

// OnRaceWatched tracks when a user watches a race
func (t *MissionTriggers) OnRaceWatched(userID string) error {
	if err := t.missionService.UpdateMissionProgress(userID, models.MissionTypeWatchRace, 1); err != nil {
		return fmt.Errorf("failed to update watch race mission progress: %w", err)
	}
	return nil
}

// OnSeriesFollowed tracks when a user follows a series
func (t *MissionTriggers) OnSeriesFollowed(userID string) error {
	if err := t.missionService.UpdateMissionProgress(userID, models.MissionTypeFollowSeries, 1); err != nil {
		return fmt.Errorf("failed to update follow series mission progress: %w", err)
	}
	return nil
}

// OnStreakDay tracks daily login/watch streak
func (t *MissionTriggers) OnStreakDay(userID string) error {
	if err := t.missionService.UpdateMissionProgress(userID, models.MissionTypeStreak, 1); err != nil {
		return fmt.Errorf("failed to update streak mission progress: %w", err)
	}
	if t.achievementService != nil {
		go t.achievementService.HandleStreak(userID)
	}
	return nil
}

// CheckAndCompleteAll checks all missions for a user and completes any that meet their targets
// This is useful to call after any action that might complete a mission
func (t *MissionTriggers) CheckAndCompleteAll(userID string) error {
	return t.missionService.CheckAndCompleteMissions(userID)
}

// OnPredictionPlaced tracks when a user places a prediction
func (t *MissionTriggers) OnPredictionPlaced(userID string) error {
	if err := t.missionService.UpdateMissionProgress(userID, models.MissionTypePredictWinner, 1); err != nil {
		return fmt.Errorf("failed to update prediction mission progress: %w", err)
	}
	return nil
}

// OnPredictionWon tracks when a user wins a prediction
// Updates progress for missions that track wins (missions with "Win" in title)
func (t *MissionTriggers) OnPredictionWon(userID string) error {
	// Get all active predict_winner missions
	missions, err := t.missionService.GetActiveMissions()
	if err != nil {
		return fmt.Errorf("failed to get missions: %w", err)
	}

	// Find win missions and update their progress
	for _, mission := range missions {
		if mission.MissionType != models.MissionTypePredictWinner {
			continue
		}

		// Check if this is a "win" mission (has "Win" in title, case-insensitive)
		isWinMission := false
		title := mission.Title
		if len(title) >= 3 {
			// Simple case-insensitive check for "win"
			for i := 0; i <= len(title)-3; i++ {
				if (title[i] == 'W' || title[i] == 'w') &&
					(title[i+1] == 'i' || title[i+1] == 'I') &&
					(title[i+2] == 'n' || title[i+2] == 'N') {
					isWinMission = true
					break
				}
			}
		}

		if !isWinMission {
			continue
		}

		// Update progress for this specific mission
		if err := t.missionService.UpdateMissionProgressForSpecificMission(userID, mission.ID, 1); err != nil {
			return fmt.Errorf("failed to update win mission progress: %w", err)
		}
	}

	return nil
}
