package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Build database connection string
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5434")
	dbUser := getEnv("DB_USER", "cyclingstream")
	dbPassword := getEnv("DB_PASSWORD", "cyclingstream_dev")
	dbName := getEnv("DB_NAME", "cyclingstream")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// Connect to database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to ping database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Verifying stream status migration...")
	fmt.Println("")

	// Count statuses
	rows, err := db.Query(`
		SELECT status, COUNT(*) as count 
		FROM streams 
		GROUP BY status 
		ORDER BY status
	`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to query stream statuses: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	fmt.Println("Current stream status distribution:")
	total := 0
	invalidStatuses := []string{}

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to scan row: %v\n", err)
			os.Exit(1)
		}
		total += count
		fmt.Printf("  %s: %d\n", status, count)

		// Check for invalid statuses
		if status != "live" && status != "offline" && status != "upcoming" {
			invalidStatuses = append(invalidStatuses, status)
		}
	}

	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error iterating rows: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nTotal streams: %d\n", total)

	// Check for races without streams
	var racesWithoutStreams int
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM races r 
		LEFT JOIN streams s ON r.id = s.race_id 
		WHERE s.id IS NULL
	`).Scan(&racesWithoutStreams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check races without streams: %v\n", err)
		os.Exit(1)
	}

	if racesWithoutStreams > 0 {
		fmt.Printf("\n⚠️  Warning: %d race(s) have no stream record\n", racesWithoutStreams)
	} else {
		fmt.Println("\n✅ All races have stream records")
	}

	// Check for invalid statuses
	if len(invalidStatuses) > 0 {
		fmt.Printf("\n❌ Invalid statuses found: %v\n", invalidStatuses)
		fmt.Println("   Expected only: live, offline, upcoming")
		os.Exit(1)
	} else {
		fmt.Println("\n✅ All statuses are valid (live, offline, upcoming)")
	}

	// Verify default value
	var defaultValue string
	err = db.QueryRow(`
		SELECT column_default 
		FROM information_schema.columns 
		WHERE table_name = 'streams' 
		AND column_name = 'status'
	`).Scan(&defaultValue)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check default value: %v\n", err)
		os.Exit(1)
	}

	if defaultValue == "'upcoming'::character varying" || defaultValue == "'upcoming'" {
		fmt.Println("✅ Default status is set to 'upcoming'")
	} else {
		fmt.Printf("⚠️  Default status is: %s (expected 'upcoming')\n", defaultValue)
	}

	fmt.Println("\n✅ Migration verification complete!")
}

