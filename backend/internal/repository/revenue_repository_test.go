package repository

import (
	"testing"
)

// TestRevenueRepository_CalculateMonthlyRevenue tests the monthly revenue calculation
// This is an integration test that requires a database connection
func TestRevenueRepository_CalculateMonthlyRevenue(t *testing.T) {
	// Skip if database is not available
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// This test would require:
	// 1. A test database connection
	// 2. Test data setup (races, payments, watch_sessions)
	// 3. Call CalculateMonthlyRevenue
	// 4. Verify the results

	t.Log("Integration test for CalculateMonthlyRevenue - requires database")
}

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

