package repository

import (
	"testing"

	"github.com/cyclingstream/backend/internal/config"
)

func TestGetLevelFromXP(t *testing.T) {
	cfg := &config.LevelingConfig{
		BaseXP:           100,
		IncrementPerLevel: 20,
	}

	tests := []struct {
		name     string
		xp       int
		expected int
	}{
		{"Level 1: 0 XP", 0, 1},
		{"Level 1: 50 XP", 50, 1},
		{"Level 1: 99 XP", 99, 1},
		{"Level 2: 100 XP", 100, 2},
		{"Level 2: 119 XP", 119, 2},
		{"Level 3: 120 XP", 120, 3},
		{"Level 3: 159 XP", 159, 3},
		{"Level 4: 160 XP", 160, 4},
		{"Level 4: 219 XP", 219, 4},
		{"Level 5: 220 XP", 220, 5},
		{"Level 10: 820 XP", 820, 10},
		{"Level 11: 1000 XP", 1000, 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLevelFromXP(tt.xp, cfg)
			if result != tt.expected {
				t.Errorf("GetLevelFromXP(%d) = %d, expected %d", tt.xp, result, tt.expected)
			}
		})
	}
}

func TestGetXPForLevel(t *testing.T) {
	cfg := &config.LevelingConfig{
		BaseXP:           100,
		IncrementPerLevel: 20,
	}

	tests := []struct {
		name     string
		level    int
		expected int
	}{
		{"Level 1: 0 XP", 1, 0},
		{"Level 2: 100 XP", 2, 100},
		{"Level 3: 120 XP", 3, 120},
		{"Level 4: 160 XP", 4, 160},
		{"Level 5: 220 XP", 5, 220},
		{"Level 10: 820 XP", 10, 820},
		{"Level 11: 1000 XP", 11, 1000},
		{"Level 0 or negative: 0 XP", 0, 0},
		{"Level -1: 0 XP", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetXPForLevel(tt.level, cfg)
			if result != tt.expected {
				t.Errorf("GetXPForLevel(%d) = %d, expected %d", tt.level, result, tt.expected)
			}
		})
	}
}

func TestGetXPForNextLevel(t *testing.T) {
	cfg := &config.LevelingConfig{
		BaseXP:           100,
		IncrementPerLevel: 20,
	}

	tests := []struct {
		name     string
		level    int
		expected int
	}{
		{"From Level 1 to 2: 100 XP", 1, 100},
		{"From Level 2 to 3: 120 XP", 2, 120},
		{"From Level 3 to 4: 160 XP", 3, 160},
		{"From Level 4 to 5: 220 XP", 4, 220},
		{"From Level 10 to 11: 1000 XP", 10, 1000},
		{"From Level 0 to 1: 0 XP", 0, 0},
		{"From Level -1 to 0: 0 XP", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetXPForNextLevel(tt.level, cfg)
			if result != tt.expected {
				t.Errorf("GetXPForNextLevel(%d) = %d, expected %d", tt.level, result, tt.expected)
			}
		})
	}
}

func TestGetLevelFromXP_CustomConfig(t *testing.T) {
	// Test with different config values
	cfg := &config.LevelingConfig{
		BaseXP:           200,
		IncrementPerLevel: 50,
	}

	tests := []struct {
		name     string
		xp       int
		expected int
	}{
		{"Level 1: 0-199 XP", 199, 1},
		{"Level 2: 200 XP", 200, 2},
		{"Level 2: 249 XP", 249, 2},
		{"Level 3: 250 XP", 250, 3},
		{"Level 3: 349 XP", 349, 3},
		{"Level 4: 400 XP", 400, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLevelFromXP(tt.xp, cfg)
			if result != tt.expected {
				t.Errorf("GetLevelFromXP(%d) = %d, expected %d", tt.xp, result, tt.expected)
			}
		})
	}
}

