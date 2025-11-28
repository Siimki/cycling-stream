package services

import (
	"testing"

	"github.com/cyclingstream/backend/internal/config"
	"github.com/cyclingstream/backend/internal/repository"
)

func TestXPService_GetLevelProgress(t *testing.T) {
	cfg := &config.LevelingConfig{
		BaseXP:           100,
		IncrementPerLevel: 20,
	}

	// Create a mock user repo (we won't actually use it for these tests)
	// In a real test, you'd use a mock or test database
	userRepo := &repository.UserRepository{}
	service := NewXPService(userRepo, cfg)

	tests := []struct {
		name                string
		xp                  int
		level               int
		expectedCurrentXP   int
		expectedNeededXP    int
	}{
		{
			name:              "Level 1, 50 XP",
			xp:                50,
			level:             1,
			expectedCurrentXP: 50,
			expectedNeededXP:  50, // Need 50 more to reach 100 (level 2)
		},
		{
			name:              "Level 1, 99 XP",
			xp:                99,
			level:             1,
			expectedCurrentXP: 99,
			expectedNeededXP:  1, // Need 1 more to reach 100 (level 2)
		},
		{
			name:              "Level 2, 100 XP",
			xp:                100,
			level:             2,
			expectedCurrentXP: 0, // Just reached level 2
			expectedNeededXP:  20, // Need 20 more to reach 120 (level 3)
		},
		{
			name:              "Level 2, 110 XP",
			xp:                110,
			level:             2,
			expectedCurrentXP: 10, // 10 XP into level 2
			expectedNeededXP:  10, // Need 10 more to reach 120 (level 3)
		},
		{
			name:              "Level 2, 119 XP",
			xp:                119,
			level:             2,
			expectedCurrentXP: 19, // 19 XP into level 2
			expectedNeededXP:  1,  // Need 1 more to reach 120 (level 3)
		},
		{
			name:              "Level 3, 120 XP",
			xp:                120,
			level:             3,
			expectedCurrentXP: 0, // Just reached level 3
			expectedNeededXP:  20, // Need 20 more to reach 140 (level 4)
		},
		{
			name:              "Level 3, 130 XP",
			xp:                130,
			level:             3,
			expectedCurrentXP: 10, // 10 XP into level 3
			expectedNeededXP:  10, // Need 10 more to reach 140 (level 4)
		},
		{
			name:              "Level 10, 260 XP",
			xp:                260,
			level:             10,
			expectedCurrentXP: 0, // Just reached level 10
			expectedNeededXP:  20, // Need 20 more to reach 280 (level 11)
		},
		{
			name:              "Level 10, 270 XP",
			xp:                270,
			level:             10,
			expectedCurrentXP: 10, // 10 XP into level 10
			expectedNeededXP:  10, // Need 10 more to reach 280 (level 11)
		},
		{
			name:              "Level 1, 0 XP",
			xp:                0,
			level:             1,
			expectedCurrentXP: 0,
			expectedNeededXP:  100, // Need 100 to reach level 2
		},
		{
			name:              "Invalid level (0), 50 XP",
			xp:                50,
			level:             0,
			expectedCurrentXP: 50, // Treated as level 1
			expectedNeededXP:  50,
		},
		{
			name:              "Invalid level (negative), 50 XP",
			xp:                50,
			level:             -1,
			expectedCurrentXP: 50, // Treated as level 1
			expectedNeededXP:  50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentXP, neededXP := service.GetLevelProgress(tt.xp, tt.level)
			if currentXP != tt.expectedCurrentXP {
				t.Errorf("GetLevelProgress(%d, %d) currentXP = %d, expected %d", tt.xp, tt.level, currentXP, tt.expectedCurrentXP)
			}
			if neededXP != tt.expectedNeededXP {
				t.Errorf("GetLevelProgress(%d, %d) neededXP = %d, expected %d", tt.xp, tt.level, neededXP, tt.expectedNeededXP)
			}
		})
	}
}

func TestXPService_CalculateLevel(t *testing.T) {
	cfg := &config.LevelingConfig{
		BaseXP:           100,
		IncrementPerLevel: 20,
	}

	userRepo := &repository.UserRepository{}
	service := NewXPService(userRepo, cfg)

	tests := []struct {
		name     string
		xp       int
		expected int
	}{
		{"0 XP", 0, 1},
		{"50 XP", 50, 1},
		{"99 XP", 99, 1},
		{"100 XP", 100, 2},
		{"119 XP", 119, 2},
		{"120 XP", 120, 3},
		{"139 XP", 139, 3},
		{"140 XP", 140, 4},
		{"260 XP", 260, 10},
		{"280 XP", 280, 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CalculateLevel(tt.xp)
			if result != tt.expected {
				t.Errorf("CalculateLevel(%d) = %d, expected %d", tt.xp, result, tt.expected)
			}
		})
	}
}

func TestXPService_GetXPForLevel(t *testing.T) {
	cfg := &config.LevelingConfig{
		BaseXP:           100,
		IncrementPerLevel: 20,
	}

	userRepo := &repository.UserRepository{}
	service := NewXPService(userRepo, cfg)

	tests := []struct {
		name     string
		level    int
		expected int
	}{
		{"Level 1", 1, 0},
		{"Level 2", 2, 100},
		{"Level 3", 3, 120},
		{"Level 4", 4, 140},
		{"Level 10", 10, 260},
		{"Level 11", 11, 280},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GetXPForLevel(tt.level)
			if result != tt.expected {
				t.Errorf("GetXPForLevel(%d) = %d, expected %d", tt.level, result, tt.expected)
			}
		})
	}
}

func TestXPService_GetXPForNextLevel(t *testing.T) {
	cfg := &config.LevelingConfig{
		BaseXP:           100,
		IncrementPerLevel: 20,
	}

	userRepo := &repository.UserRepository{}
	service := NewXPService(userRepo, cfg)

	tests := []struct {
		name     string
		level    int
		expected int
	}{
		{"From Level 1", 1, 100},
		{"From Level 2", 2, 120},
		{"From Level 3", 3, 140},
		{"From Level 10", 10, 280},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GetXPForNextLevel(tt.level)
			if result != tt.expected {
				t.Errorf("GetXPForNextLevel(%d) = %d, expected %d", tt.level, result, tt.expected)
			}
		})
	}
}

