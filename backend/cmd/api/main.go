package main

import (
	"github.com/cyclingstream/backend/internal/chat"
	"github.com/cyclingstream/backend/internal/config"
	"github.com/cyclingstream/backend/internal/database"
	"github.com/cyclingstream/backend/internal/logger"
	"github.com/cyclingstream/backend/internal/server"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		// Use standard log for fatal errors before logger is initialized
		panic("Failed to load config: " + err.Error())
	}

	// Initialize structured logger
	logger.Init(cfg.Env)
	logger.WithField("environment", cfg.Env).Info("Logger initialized")

	// Connect to database
	db, err := database.New(cfg.GetDSN())
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logger.Info("Database connection established")

	// Initialize chat hub and rate limiter
	hub := chat.NewHub()
	rateLimiter := chat.NewRateLimiter()

	// Start hub in background
	go hub.Run()

	logger.Info("Chat hub started")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			
			// Log error with context
			logger.WithFields(map[string]interface{}{
				"path":   c.Path(),
				"method": c.Method(),
				"status": code,
				"error":  err.Error(),
			}).Error("Request error")
			
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Setup routes
	server.SetupRoutes(app, db, cfg, hub, rateLimiter)

	// Start server
	addr := ":" + cfg.Port
	logger.WithField("address", addr).Info("Server starting")
	if err := app.Listen(addr); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}

