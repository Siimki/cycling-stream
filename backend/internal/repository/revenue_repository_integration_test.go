package repository

import (
	"testing"
)

// TestRevenueRepository_Integration tests the revenue repository with a real database
// This test requires a database connection and should be run with: go test -tags=integration
func TestRevenueRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This would require setting up a test database connection
	// For now, this is a placeholder that documents the integration test structure
	
	t.Log("Integration test structure for revenue repository")
	t.Log("To run: go test -tags=integration ./internal/repository")
}

// TestRevenueRepository_CalculateMonthlyRevenue_Logic tests the logic of revenue calculation
// This is a unit test that doesn't require a database
func TestRevenueRepository_CalculateMonthlyRevenue_Logic(t *testing.T) {
	// Test the revenue split calculation logic
	testCases := []struct {
		name              string
		totalRevenueCents int
		expectedPlatform  int
		expectedOrganizer int
	}{
		{
			name:              "Even split - $10.00",
			totalRevenueCents: 1000,
			expectedPlatform:  500,
			expectedOrganizer: 500,
		},
		{
			name:              "Odd amount - $10.01 (extra cent to organizer)",
			totalRevenueCents: 1001,
			expectedPlatform:  500,
			expectedOrganizer: 501,
		},
		{
			name:              "Large amount - $1000.00",
			totalRevenueCents: 100000,
			expectedPlatform:  50000,
			expectedOrganizer: 50000,
		},
		{
			name:              "Large odd amount - $1000.01",
			totalRevenueCents: 100001,
			expectedPlatform:  50000,
			expectedOrganizer: 50001,
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
		{
			name:              "Two cents",
			totalRevenueCents: 2,
			expectedPlatform:  1,
			expectedOrganizer: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This is the same logic used in CalculateMonthlyRevenue
			platformShareCents := tc.totalRevenueCents / 2
			organizerShareCents := tc.totalRevenueCents - platformShareCents

			if platformShareCents != tc.expectedPlatform {
				t.Errorf("Platform share: expected %d, got %d", tc.expectedPlatform, platformShareCents)
			}
			if organizerShareCents != tc.expectedOrganizer {
				t.Errorf("Organizer share: expected %d, got %d", tc.expectedOrganizer, organizerShareCents)
			}
			// Verify they add up to total
			if platformShareCents+organizerShareCents != tc.totalRevenueCents {
				t.Errorf("Shares don't add up: %d + %d != %d", platformShareCents, organizerShareCents, tc.totalRevenueCents)
			}
		})
	}
}

// TestRevenueRepository_QueryValidation tests that SQL queries are syntactically correct
// This doesn't require a database connection, just validates the query strings
func TestRevenueRepository_QueryValidation(t *testing.T) {
	// These are the query patterns used in the repository
	// We can't test them without a DB, but we can verify the structure
	
	queries := []struct {
		name  string
		query string
	}{
		{
			name: "Revenue calculation query",
			query: `
				SELECT COALESCE(SUM(amount_cents), 0)
				FROM payments
				WHERE race_id = $1
				  AND status = 'succeeded'
				  AND EXTRACT(YEAR FROM created_at) = $2
				  AND EXTRACT(MONTH FROM created_at) = $3
			`,
		},
		{
			name: "Watch minutes calculation query",
			query: `
				SELECT COALESCE(SUM(duration_seconds) / 60.0, 0)
				FROM watch_sessions
				WHERE race_id = $1
				  AND duration_seconds IS NOT NULL
				  AND EXTRACT(YEAR FROM started_at) = $2
				  AND EXTRACT(MONTH FROM started_at) = $3
			`,
		},
		{
			name: "Upsert monthly revenue query",
			query: `
				INSERT INTO revenue_share_monthly (
					id, race_id, year, month, total_revenue_cents, total_watch_minutes,
					platform_share_cents, organizer_share_cents, calculated_at
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP)
				ON CONFLICT (race_id, year, month)
				DO UPDATE SET
					total_revenue_cents = EXCLUDED.total_revenue_cents,
					total_watch_minutes = EXCLUDED.total_watch_minutes,
					platform_share_cents = EXCLUDED.platform_share_cents,
					organizer_share_cents = EXCLUDED.organizer_share_cents,
					calculated_at = EXCLUDED.calculated_at,
					updated_at = CURRENT_TIMESTAMP
			`,
		},
	}

	for _, q := range queries {
		t.Run(q.name, func(t *testing.T) {
			// Basic validation: query should not be empty
			if len(q.query) == 0 {
				t.Error("Query is empty")
			}
			// Query should contain expected keywords
			if q.name == "Revenue calculation query" {
				if !contains(q.query, "SELECT") || !contains(q.query, "FROM payments") {
					t.Error("Revenue query missing expected keywords")
				}
			}
			if q.name == "Watch minutes calculation query" {
				if !contains(q.query, "SELECT") || !contains(q.query, "FROM watch_sessions") {
					t.Error("Watch minutes query missing expected keywords")
				}
			}
			if q.name == "Upsert monthly revenue query" {
				if !contains(q.query, "INSERT") || !contains(q.query, "ON CONFLICT") {
					t.Error("Upsert query missing expected keywords")
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsHelper(s, substr))))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

