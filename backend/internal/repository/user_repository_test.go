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
		{"Level 3: 139 XP", 139, 3},
		{"Level 4: 140 XP", 140, 4},
		{"Level 4: 159 XP", 159, 4},
		{"Level 5: 160 XP", 160, 5},
		{"Level 10: 260 XP", 260, 10},
		{"Level 11: 280 XP", 280, 11},
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
		{"Level 4: 140 XP", 4, 140},
		{"Level 5: 160 XP", 5, 160},
		{"Level 10: 260 XP", 10, 260},
		{"Level 11: 280 XP", 11, 280},
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
		{"From Level 3 to 4: 140 XP", 3, 140},
		{"From Level 4 to 5: 160 XP", 4, 160},
		{"From Level 10 to 11: 280 XP", 10, 280},
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
		{"Level 3: 299 XP", 299, 3},
		{"Level 4: 300 XP", 300, 4},
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

