package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// GetTestDB returns a test database connection
// It uses environment variables or defaults for test database connection
func GetTestDB(t *testing.T) *sql.DB {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "cyclingstream"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "cyclingstream_dev"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "cyclingstream"
	}

	dbSSLMode := os.Getenv("DB_SSLMODE")
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	return db
}

// CleanupChatMessages removes all chat messages from the database (for test cleanup)
func CleanupChatMessages(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DELETE FROM chat_messages")
	if err != nil {
		t.Logf("Warning: Failed to cleanup chat_messages: %v", err)
	}
}

// CleanupRaces removes test races from the database (for test cleanup)
func CleanupRaces(t *testing.T, db *sql.DB, raceIDs []string) {
	for _, raceID := range raceIDs {
		_, err := db.Exec("DELETE FROM races WHERE id = $1", raceID)
		if err != nil {
			t.Logf("Warning: Failed to cleanup race %s: %v", raceID, err)
		}
	}
}

// CleanupUsers removes test users from the database (for test cleanup)
func CleanupUsers(t *testing.T, db *sql.DB, userIDs []string) {
	for _, userID := range userIDs {
		_, err := db.Exec("DELETE FROM users WHERE id = $1", userID)
		if err != nil {
			t.Logf("Warning: Failed to cleanup user %s: %v", userID, err)
		}
	}
}

// CreateTestRace creates a test race in the database and returns the race ID
func CreateTestRace(t *testing.T, db *sql.DB, name string) string {
	query := `
		INSERT INTO races (id, name, description, start_date, end_date, location, category, is_free, price_cents, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, 'Test race', NOW(), NOW() + INTERVAL '2 hours', 'Test Location', 'Test', true, 0, NOW(), NOW())
		RETURNING id
	`

	var raceID string
	err := db.QueryRow(query, name).Scan(&raceID)
	if err != nil {
		t.Fatalf("Failed to create test race: %v", err)
	}

	return raceID
}

// CreateTestUser creates a test user in the database and returns the user ID
// Password will be hashed using bcrypt
func CreateTestUser(t *testing.T, db *sql.DB, email, password, name string) string {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	query := `
		INSERT INTO users (id, email, password_hash, name, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, NOW(), NOW())
		RETURNING id
	`

	var userID string
	err = db.QueryRow(query, email, string(hashedPassword), name).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return userID
}

// CreateTestStream creates a test stream for a race
func CreateTestStream(t *testing.T, db *sql.DB, raceID, status string) {
	query := `
		INSERT INTO streams (id, race_id, status, origin_url, cdn_url, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, 'http://test.com/stream.m3u8', 'http://cdn.test.com/stream.m3u8', NOW(), NOW())
		ON CONFLICT (race_id) DO UPDATE SET status = $2, updated_at = NOW()
	`

	_, err := db.Exec(query, raceID, status)
	if err != nil {
		t.Fatalf("Failed to create test stream: %v", err)
	}
}

