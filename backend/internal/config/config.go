package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	Env                 string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPass              string
	DBName              string
	DBSSLMode           string
	JWTSecret           string
	StripeKey           string
	StripeWebhookSecret string
	FrontendURL         string
	XP                  *XPConfig
	Bunny               *BunnyConfig
}

type BunnyConfig struct {
	APIKey    string
	LibraryID string
	BaseURL   string
}

func Load() (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	env := getEnv("ENV", "development")
	isProduction := env == "production"

	config := &Config{
		Port:   getEnv("PORT", "8080"),
		Env:    env,
		DBHost: getEnv("DB_HOST", "localhost"),
		// Default to the Docker-exposed Postgres port from docker-compose (5434:5432)
		DBPort:              getEnv("DB_PORT", "5434"),
		DBUser:              getEnv("DB_USER", "cyclingstream"),
		DBPass:              getEnv("DB_PASSWORD", "cyclingstream_dev"),
		DBName:              getEnv("DB_NAME", "cyclingstream"),
		DBSSLMode:           getEnv("DB_SSLMODE", "disable"),
		JWTSecret:           getEnv("JWT_SECRET", "change-me-in-production"),
		StripeKey:           getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		FrontendURL:         getEnv("FRONTEND_URL", "http://localhost:3000"),
		XP:                  LoadXPConfig(),
		Bunny:               LoadBunnyConfig(),
	}

	// Validate configuration
	if err := config.Validate(isProduction); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// Validate checks that all required configuration values are set appropriately
func (c *Config) Validate(isProduction bool) error {
	var errors []string

	// JWT Secret validation
	if c.JWTSecret == "" {
		errors = append(errors, "JWT_SECRET is required")
	} else if isProduction && (c.JWTSecret == "change-me-in-production" || len(c.JWTSecret) < 32) {
		errors = append(errors, "JWT_SECRET must be a secure random string (at least 32 characters) in production")
	}

	// Database configuration validation
	if c.DBHost == "" {
		errors = append(errors, "DB_HOST is required")
	}
	if c.DBPort == "" {
		errors = append(errors, "DB_PORT is required")
	}
	if c.DBUser == "" {
		errors = append(errors, "DB_USER is required")
	}
	if c.DBPass == "" {
		errors = append(errors, "DB_PASSWORD is required")
	}
	if c.DBName == "" {
		errors = append(errors, "DB_NAME is required")
	}

	// Stripe configuration validation (required for payment features)
	// Note: In production, Stripe keys should be set, but we allow empty in development
	// for testing without payment features
	if isProduction {
		if c.StripeKey == "" {
			errors = append(errors, "STRIPE_SECRET_KEY is required in production")
		} else if !strings.HasPrefix(c.StripeKey, "sk_") {
			errors = append(errors, "STRIPE_SECRET_KEY must start with 'sk_'")
		}
		if c.StripeWebhookSecret == "" {
			errors = append(errors, "STRIPE_WEBHOOK_SECRET is required in production")
		} else if !strings.HasPrefix(c.StripeWebhookSecret, "whsec_") {
			errors = append(errors, "STRIPE_WEBHOOK_SECRET must start with 'whsec_'")
		}
	} else {
		// In development, warn if Stripe keys are partially set
		if (c.StripeKey != "" && c.StripeWebhookSecret == "") || (c.StripeKey == "" && c.StripeWebhookSecret != "") {
			// This is just a warning, not an error
			fmt.Fprintf(os.Stderr, "WARNING: Stripe configuration is incomplete. Payment features may not work correctly.\n")
		}
	}

	// Frontend URL validation
	if c.FrontendURL == "" {
		errors = append(errors, "FRONTEND_URL is required")
	} else if isProduction && !strings.HasPrefix(c.FrontendURL, "https://") {
		errors = append(errors, "FRONTEND_URL must use HTTPS in production")
	}

	// Port validation
	if port, err := strconv.Atoi(c.Port); err != nil || port < 1 || port > 65535 {
		errors = append(errors, "PORT must be a valid port number (1-65535)")
	}

	// Bunny config (warn in dev, require in prod if enabled)
	if c.Bunny != nil {
		if isProduction {
			if c.Bunny.APIKey == "" {
				errors = append(errors, "BUNNY_API_KEY is required in production for analytics sync")
			}
			if c.Bunny.LibraryID == "" {
				errors = append(errors, "BUNNY_LIBRARY_ID is required in production for analytics sync")
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func LoadBunnyConfig() *BunnyConfig {
	return &BunnyConfig{
		APIKey:    getEnv("BUNNY_API_KEY", ""),
		LibraryID: getEnv("BUNNY_LIBRARY_ID", ""),
		BaseURL:   getEnv("BUNNY_API_BASE_URL", "https://video.bunnycdn.com/library"),
	}
}
