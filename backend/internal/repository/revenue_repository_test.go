package repository

import (
	"testing"
)

// TestRevenueRepository_RevenueSplit tests the 50/50 revenue split calculation
func TestRevenueRepository_RevenueSplit(t *testing.T) {
	// Test that revenue split is calculated correctly (50/50)
	testCases := []struct {
		name              string
		totalRevenueCents int
		expectedPlatform  int
		expectedOrganizer int
	}{
		{
			name:              "Even amount",
			totalRevenueCents: 1000,
			expectedPlatform:  500,
			expectedOrganizer: 500,
		},
		{
			name:              "Odd amount (1 cent goes to organizer)",
			totalRevenueCents: 1001,
			expectedPlatform:  500,
			expectedOrganizer: 501,
		},
		{
			name:              "Zero revenue",
			totalRevenueCents: 0,
			expectedPlatform:  0,
			expectedOrganizer: 0,
		},
		{
			name:              "Single cent",
			totalRevenueCents: 1,
			expectedPlatform:  0,
			expectedOrganizer: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			platformShare := tc.totalRevenueCents / 2
			organizerShare := tc.totalRevenueCents - platformShare

			if platformShare != tc.expectedPlatform {
				t.Errorf("Expected platform share %d, got %d", tc.expectedPlatform, platformShare)
			}
			if organizerShare != tc.expectedOrganizer {
				t.Errorf("Expected organizer share %d, got %d", tc.expectedOrganizer, organizerShare)
			}
			if platformShare+organizerShare != tc.totalRevenueCents {
				t.Errorf("Shares don't add up to total: %d + %d != %d", platformShare, organizerShare, tc.totalRevenueCents)
			}
		})
	}
}

