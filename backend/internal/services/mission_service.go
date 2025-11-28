package services

import (
	"fmt"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
)

type MissionService struct {
	missionRepo     *repository.MissionRepository
	userMissionRepo *repository.UserMissionRepository
	userRepo        *repository.UserRepository
	xpService       *XPService
}

func NewMissionService(
	missionRepo *repository.MissionRepository,
	userMissionRepo *repository.UserMissionRepository,
	userRepo *repository.UserRepository,
	xpService *XPService,
) *MissionService {
	return &MissionService{
		missionRepo:     missionRepo,
		userMissionRepo: userMissionRepo,
		userRepo:        userRepo,
		xpService:       xpService,
	}
}

// GetActiveMissionsForUser returns all active missions with user progress
func (s *MissionService) GetActiveMissionsForUser(userID string) ([]models.UserMissionWithDetails, error) {
	// Get all active missions
	activeMissions, err := s.missionRepo.GetActiveMissions()
	if err != nil {
		return nil, fmt.Errorf("failed to get active missions: %w", err)
	}

	// Get or create user missions for each active mission
	var userMissions []models.UserMissionWithDetails
	for _, mission := range activeMissions {
		userMission, err := s.userMissionRepo.GetOrCreate(userID, mission.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get or create user mission: %w", err)
		}

		userMissions = append(userMissions, models.UserMissionWithDetails{
			UserMission: *userMission,
			Mission:     mission,
		})
	}

	return userMissions, nil
}

// GetUserMissions returns all missions for a user (including completed/claimed)
func (s *MissionService) GetUserMissions(userID string) ([]models.UserMissionWithDetails, error) {
	return s.userMissionRepo.GetByUserID(userID)
}

// UpdateMissionProgress updates progress for a mission type and checks for completion
// For predict_winner missions, it updates all missions (both "place" and "win" types)
// The caller should filter if needed (e.g., OnPredictionWon should only update "win" missions)
func (s *MissionService) UpdateMissionProgress(userID string, missionType models.MissionType, increment int) error {
	// Get all active missions of this type
	missions, err := s.missionRepo.GetByType(missionType)
	if err != nil {
		return fmt.Errorf("failed to get missions by type: %w", err)
	}

	// Update progress for each mission
	for _, mission := range missions {
		// Get or create user mission
		userMission, err := s.userMissionRepo.GetOrCreate(userID, mission.ID)
		if err != nil {
			return fmt.Errorf("failed to get or create user mission: %w", err)
		}

		// Skip if already completed
		if userMission.IsCompleted() {
			continue
		}

		// Increment progress
		newProgress := userMission.Progress + increment
		if err := s.userMissionRepo.UpdateProgress(userID, mission.ID, newProgress); err != nil {
			return fmt.Errorf("failed to update progress: %w", err)
		}

		// Check if mission is completed
		if newProgress >= mission.TargetValue {
			if err := s.userMissionRepo.Complete(userID, mission.ID); err != nil {
				return fmt.Errorf("failed to complete mission: %w", err)
			}
		}
	}

	return nil
}

// UpdateMissionProgressForSpecificMission updates progress for a specific mission by ID
func (s *MissionService) UpdateMissionProgressForSpecificMission(userID, missionID string, increment int) error {
	// Get mission
	mission, err := s.missionRepo.GetByID(missionID)
	if err != nil {
		return fmt.Errorf("failed to get mission: %w", err)
	}
	if mission == nil {
		return fmt.Errorf("mission not found")
	}

	// Get or create user mission
	userMission, err := s.userMissionRepo.GetOrCreate(userID, missionID)
	if err != nil {
		return fmt.Errorf("failed to get or create user mission: %w", err)
	}

	// Skip if already completed
	if userMission.IsCompleted() {
		return nil
	}

	// Increment progress
	newProgress := userMission.Progress + increment
	if err := s.userMissionRepo.UpdateProgress(userID, missionID, newProgress); err != nil {
		return fmt.Errorf("failed to update progress: %w", err)
	}

	// Check if mission is completed
	if newProgress >= mission.TargetValue {
		if err := s.userMissionRepo.Complete(userID, missionID); err != nil {
			return fmt.Errorf("failed to complete mission: %w", err)
		}
	}

	return nil
}

// CheckAndCompleteMissions checks all missions for a user and completes any that meet their targets
func (s *MissionService) CheckAndCompleteMissions(userID string) error {
	userMissions, err := s.userMissionRepo.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user missions: %w", err)
	}

	for _, um := range userMissions {
		// Skip if already completed
		if um.IsCompleted() {
			continue
		}

		// Check if progress meets target
		if um.Progress >= um.Mission.TargetValue {
			if err := s.userMissionRepo.Complete(userID, um.MissionID); err != nil {
				return fmt.Errorf("failed to complete mission: %w", err)
			}
		}
	}

	return nil
}

// ClaimMissionReward claims the reward for a completed mission and awards points
func (s *MissionService) ClaimMissionReward(userID, missionID string) error {
	// Get user mission
	userMission, err := s.userMissionRepo.GetByUserAndMission(userID, missionID)
	if err != nil {
		return fmt.Errorf("failed to get user mission: %w", err)
	}
	if userMission == nil {
		return fmt.Errorf("user mission not found")
	}

	// Check if mission is completed
	if !userMission.IsCompleted() {
		return fmt.Errorf("mission not completed")
	}

	// Check if already claimed
	if userMission.IsClaimed() {
		return fmt.Errorf("mission reward already claimed")
	}

	// Get mission details
	mission, err := s.missionRepo.GetByID(missionID)
	if err != nil {
		return fmt.Errorf("failed to get mission: %w", err)
	}
	if mission == nil {
		return fmt.Errorf("mission not found")
	}

	// Claim the reward
	if err := s.userMissionRepo.Claim(userID, missionID); err != nil {
		return fmt.Errorf("failed to claim mission: %w", err)
	}

	// Award points
	if err := s.userRepo.AddPoints(userID, mission.PointsReward); err != nil {
		// If points award fails, we should still mark as claimed to avoid double claiming
		// Log the error but don't fail the claim
		return fmt.Errorf("failed to award points (mission already claimed): %w", err)
	}

	// Award XP if mission has XP reward
	if mission.XPReward > 0 && s.xpService != nil {
		if err := s.xpService.AwardXP(userID, mission.XPReward, "mission"); err != nil {
			// Log error but don't fail the claim
			fmt.Printf("Failed to award XP for mission: %v\n", err)
		}
	}

	return nil
}

// InitializeUserMissions creates user mission entries for all active missions
func (s *MissionService) InitializeUserMissions(userID string) error {
	activeMissions, err := s.missionRepo.GetActiveMissions()
	if err != nil {
		return fmt.Errorf("failed to get active missions: %w", err)
	}

	for _, mission := range activeMissions {
		_, err := s.userMissionRepo.GetOrCreate(userID, mission.ID)
		if err != nil {
			return fmt.Errorf("failed to initialize user mission: %w", err)
		}
	}

	return nil
}

// GetActiveMissions returns all active missions (for admin/public use)
func (s *MissionService) GetActiveMissions() ([]models.Mission, error) {
	return s.missionRepo.GetActiveMissions()
}

// GetActiveCareerMissionsForUser returns one active career mission per mission type (the highest incomplete tier)
func (s *MissionService) GetActiveCareerMissionsForUser(userID string) ([]models.UserMissionWithDetails, error) {
	// Get all career mission types
	missionTypes := []models.MissionType{
		models.MissionTypeWatchTime,
		models.MissionTypeChatMessage,
		models.MissionTypePredictWinner,
	}

	var result []models.UserMissionWithDetails

	for _, missionType := range missionTypes {
		// Get all career missions of this type, ordered by tier
		careerMissions, err := s.missionRepo.GetActiveCareerMissionsByType(missionType)
		if err != nil {
			return nil, fmt.Errorf("failed to get career missions for type %s: %w", missionType, err)
		}

		// Find the highest incomplete tier mission
		var activeMission *models.Mission
		for i := len(careerMissions) - 1; i >= 0; i-- {
			mission := careerMissions[i]
			userMission, err := s.userMissionRepo.GetOrCreate(userID, mission.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get user mission: %w", err)
			}

			// If this mission is not completed, it's the active one
			if !userMission.IsCompleted() {
				activeMission = &mission
				break
			}
		}

		// If all tiers are completed, show the highest tier (or none if no missions exist)
		if activeMission == nil && len(careerMissions) > 0 {
			activeMission = &careerMissions[len(careerMissions)-1]
		}

		if activeMission != nil {
			userMission, err := s.userMissionRepo.GetOrCreate(userID, activeMission.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get user mission: %w", err)
			}

			result = append(result, models.UserMissionWithDetails{
				UserMission: *userMission,
				Mission:     *activeMission,
			})
		}
	}

	return result, nil
}

// GetWeeklyMissionsForUser returns all weekly missions for a user
func (s *MissionService) GetWeeklyMissionsForUser(userID string) ([]models.UserMissionWithDetails, error) {
	weeklyMissions, err := s.missionRepo.GetWeeklyMissions()
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly missions: %w", err)
	}

	var result []models.UserMissionWithDetails
	for _, mission := range weeklyMissions {
		userMission, err := s.userMissionRepo.GetOrCreate(userID, mission.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user mission: %w", err)
		}

		result = append(result, models.UserMissionWithDetails{
			UserMission: *userMission,
			Mission:     mission,
		})
	}

	return result, nil
}

