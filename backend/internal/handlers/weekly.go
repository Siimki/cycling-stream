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
		"chat_messages":         progress.ChatMessages,
		"weekly_goal_completed": progress.GoalCompleted,
		"current_streak_weeks":  progress.CurrentStreakWeeks,
		"best_streak_weeks":     progress.BestStreakWeeks,
		"can_claim_reward":      progress.GoalCompleted && !progress.RewardClaimed,
		"reward_xp":             150,
		"reward_points":         200,
	})
}


