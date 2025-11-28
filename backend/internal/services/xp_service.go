package services

import (
	"fmt"

	"github.com/cyclingstream/backend/internal/config"
	"github.com/cyclingstream/backend/internal/repository"
)

type XPService struct {
	userRepo *repository.UserRepository
	cfg      *config.LevelingConfig
}

func NewXPService(userRepo *repository.UserRepository, cfg *config.LevelingConfig) *XPService {
	return &XPService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// AwardXP awards XP to a user and recalculates their level if needed.
// It returns the new level if it changed, or 0 if it didn't.
func (s *XPService) AwardXP(userID string, xp int, source string) error {
	if xp <= 0 {
		return nil
	}

	// Add XP
	if err := s.userRepo.AddXP(userID, xp); err != nil {
		return fmt.Errorf("failed to award XP: %w", err)
	}

	// Get updated user to check level
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user after XP award: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Calculate new level from updated XP
	newLevel := repository.GetLevelFromXP(user.XPTotal, s.cfg)

	// Update level if it changed
	if newLevel > user.Level {
		if err := s.userRepo.UpdateLevel(userID, newLevel); err != nil {
			return fmt.Errorf("failed to update level: %w", err)
		}
	}

	return nil
}

// CalculateLevel calculates the level from total XP using the configured formula.
func (s *XPService) CalculateLevel(xp int) int {
	return repository.GetLevelFromXP(xp, s.cfg)
}

// GetXPForLevel returns the XP needed to REACH level N (the minimum XP for that level).
func (s *XPService) GetXPForLevel(level int) int {
	return repository.GetXPForLevel(level, s.cfg)
}

// GetXPForNextLevel returns the XP needed to reach the NEXT level (level N+1) from the current level N.
func (s *XPService) GetXPForNextLevel(level int) int {
	return repository.GetXPForNextLevel(level, s.cfg)
}

// GetLevelProgress returns the current XP progress within the current level.
// Returns: (currentXPInLevel, xpNeededForNextLevel)
// currentXPInLevel: XP earned within the current level (from level start to current XP)
// xpNeededForNextLevel: XP still needed to reach the next level
func (s *XPService) GetLevelProgress(xp int, level int) (currentXP, neededXP int) {
	if level < 1 {
		level = 1
	}

	// XP threshold for current level start
	xpForCurrentLevelStart := repository.GetXPForLevel(level, s.cfg)

	// XP threshold for next level start
	xpForNextLevelStart := repository.GetXPForNextLevel(level, s.cfg)

	// Current XP within this level (from level start to current XP)
	currentXPInLevel := xp - xpForCurrentLevelStart
	if currentXPInLevel < 0 {
		currentXPInLevel = 0
	}

	// XP still needed to reach next level
	xpNeeded := xpForNextLevelStart - xp
	if xpNeeded < 0 {
		xpNeeded = 0
	}

	return currentXPInLevel, xpNeeded
}


