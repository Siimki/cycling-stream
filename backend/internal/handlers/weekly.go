package handlers

import (
	"github.com/cyclingstream/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

type WeeklyHandler struct {
	weeklyService *services.WeeklyService
}

func NewWeeklyHandler(weeklyService *services.WeeklyService) *WeeklyHandler {
	return &WeeklyHandler{
		weeklyService: weeklyService,
	}
}

// GetWeeklyProgress returns the current week's progress and streak info for the authenticated user
func (h *WeeklyHandler) GetWeeklyProgress(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	progress, err := h.weeklyService.GetWeeklyProgress(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to get weekly progress",
		})
	}

	// Map to frontend expected format
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":               userID,
		"week_number":           progress.WeekNumber,
		"watch_minutes":         progress.WatchMinutes,
		"watch_minutes_goal":    progress.WatchMinutesGoal,
		"chat_messages":         progress.ChatMessages,
		"chat_messages_goal":    progress.ChatMessagesGoal,
		"weekly_goal_completed": progress.GoalCompleted,
		"current_streak_weeks":  progress.CurrentStreakWeeks,
		"best_streak_weeks":     progress.BestStreakWeeks,
		"can_claim_reward":      progress.GoalCompleted && !progress.RewardClaimed,
		"reward_xp":             progress.RewardXP,
		"reward_points":         progress.RewardPoints,
	})
}

// ClaimWeeklyReward claims the weekly goal reward
func (h *WeeklyHandler) ClaimWeeklyReward(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	// Get current week progress to validate
	progress, err := h.weeklyService.GetWeeklyProgress(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to get weekly progress",
		})
	}

	// Check if goal is completed
	if !progress.GoalCompleted {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{
			Error: "Weekly goal not completed",
		})
	}

	// Check if already claimed
	if progress.RewardClaimed {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{
			Error: "Weekly reward already claimed",
		})
	}

	// Claim the reward (this will award XP and points)
	if err := h.weeklyService.ClaimWeeklyReward(userID, progress.WeekNumber); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to claim weekly reward",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Weekly reward claimed successfully",
	})
}


