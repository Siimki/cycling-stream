package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/cyclingstream/backend/internal/config"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/joho/godotenv"
)

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

	// Load XP config (using defaults from config package)
	xpConfig := config.LoadXPConfig()

	// Get all users
	rows, err := db.Query("SELECT id, email, name, xp_total, level FROM users ORDER BY xp_total DESC")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to query users: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	fmt.Println("Recalculating user levels based on XP totals...")
	fmt.Println("================================================")

	var updated, unchanged, errors int

	for rows.Next() {
		var id, email string
		var name sql.NullString
		var xpTotal, currentLevel int

		if err := rows.Scan(&id, &email, &name, &xpTotal, &currentLevel); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to scan user: %v\n", err)
			errors++
			continue
		}

		// Get display name
		displayName := email
		if name.Valid {
			displayName = name.String
		}

		// Calculate correct level from XP using new formula
		correctLevel := repository.GetLevelFromXP(xpTotal, &xpConfig.Leveling)

		if correctLevel != currentLevel {
			// Update level
			_, err := db.Exec("UPDATE users SET level = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
				correctLevel, id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to update user %s (%s): %v\n", displayName, id, err)
				errors++
				continue
			}

			fmt.Printf("✓ Updated %s (email: %s): Level %d → %d (XP: %d)\n",
				displayName, email, currentLevel, correctLevel, xpTotal)
			updated++
		} else {
			fmt.Printf("  %s (email: %s): Level %d (XP: %d) - correct\n",
				displayName, email, currentLevel, xpTotal)
			unchanged++
		}
	}

	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error iterating users: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("================================================")
	fmt.Printf("Recalculation complete!\n")
	fmt.Printf("  Updated: %d users\n", updated)
	fmt.Printf("  Unchanged: %d users\n", unchanged)
	if errors > 0 {
		fmt.Printf("  Errors: %d users\n", errors)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
