package server

import (
	"log"

	"github.com/cyclingstream/backend/internal/chat"
	"github.com/cyclingstream/backend/internal/config"
	"github.com/cyclingstream/backend/internal/database"
	"github.com/cyclingstream/backend/internal/handlers"
	"github.com/cyclingstream/backend/internal/middleware"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/cyclingstream/backend/internal/services"
	"github.com/cyclingstream/backend/internal/services/analytics"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App, db *database.DB, cfg *config.Config, hub *chat.Hub, rateLimiter *chat.RateLimiter) {
	// Global middleware - order matters!
	app.Use(middleware.StructuredLogger())
	app.Use(middleware.SecurityHeaders()) // Security headers first

	// CORS configuration - environment-aware
	corsOrigins := cfg.FrontendURL
	if cfg.Env == "development" {
		// In development, allow localhost:3000 by default if not set in FrontendURL,
		// but since FrontendURL defaults to localhost:3000, we can just use that.
		// If we need multiple, we can comma separate them.
		if corsOrigins == "*" {
			// If someone explicitly set FRONTEND_URL=*, we can't use AllowCredentials=true
			// But let's assume we want credentials.
			corsOrigins = "http://localhost:3000"
		}
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization,X-CSRF-Token",
		AllowCredentials: true,
	}))

	// Initialize repositories
	raceRepo := repository.NewRaceRepository(db.DB)
	streamRepo := repository.NewStreamRepository(db.DB)
	streamProviderRepo := repository.NewStreamProviderRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)
	paymentRepo := repository.NewPaymentRepository(db.DB)
	entitlementRepo := repository.NewEntitlementRepository(db.DB)
	watchSessionRepo := repository.NewWatchSessionRepository(db.DB)
	revenueRepo := repository.NewRevenueRepository(db.DB)
	viewerSessionRepo := repository.NewViewerSessionRepository(db.DB)
	costRepo := repository.NewCostRepository(db.DB)
	playbackEventRepo := repository.NewPlaybackEventRepository(db.DB)
	streamStatsRepo := repository.NewStreamStatsRepository(db.DB)
	bunnyStatsRepo := repository.NewBunnyStatsRepository(db.DB)
	chatRepo := repository.NewChatRepository(db.DB)
	achievementRepo := repository.NewAchievementRepository(db.DB)
	userPrefsRepo := repository.NewUserPreferencesRepository(db.DB)
	userFavRepo := repository.NewUserFavoriteRepository(db.DB)
	watchHistoryRepo := repository.NewWatchHistoryRepository(db.DB)
	recommendationService := services.NewRecommendationService(raceRepo, watchHistoryRepo, userFavRepo, streamRepo)
	missionRepo := repository.NewMissionRepository(db.DB)
	userMissionRepo := repository.NewUserMissionRepository(db.DB)
	xpService := services.NewXPService(userRepo, &cfg.XP.Leveling)
	missionService := services.NewMissionService(missionRepo, userMissionRepo, userRepo, xpService)
	weeklyRepo := repository.NewWeeklyRepository(db.DB)
	streakRepo := repository.NewStreakRepository(db.DB)
	weeklyService := services.NewWeeklyService(weeklyRepo, streakRepo, userRepo, xpService, cfg.XP)
	achievementService := services.NewAchievementService(achievementRepo, chatRepo, watchSessionRepo, userRepo)
	if err := achievementService.SeedDefaults(); err != nil {
		log.Printf("failed to seed achievements: %v", err)
	}
	xpService.SetAchievementService(achievementService)
	missionTriggers := services.NewMissionTriggers(missionService, xpService, weeklyService, achievementService, cfg.XP)
	predictionRepo := repository.NewPredictionRepository(db.DB)
	predictionService := services.NewPredictionService(predictionRepo, userRepo, xpService, missionTriggers, cfg.XP)
	analyticsAggregator := analytics.NewAggregator(playbackEventRepo, streamStatsRepo, streamRepo)
	var bunnyImporter *analytics.BunnyImporter
	bunnyEnabled := cfg.Bunny != nil && cfg.Bunny.APIKey != "" && cfg.Bunny.LibraryID != ""
	if bunnyEnabled {
		bunnyClient := analytics.NewBunnyClient(cfg.Bunny)
		bunnyImporter = analytics.NewBunnyImporter(bunnyClient, streamProviderRepo, bunnyStatsRepo, streamStatsRepo)
	}

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(db)
	raceHandler := handlers.NewRaceHandler(raceRepo, streamRepo, entitlementRepo)
	streamHandler := handlers.NewStreamHandler(streamRepo)
	authHandler := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)
	adminHandler := handlers.NewAdminHandler(raceRepo, streamRepo, revenueRepo)
	paymentHandler := handlers.NewPaymentHandler(
		paymentRepo,
		entitlementRepo,
		raceRepo,
		cfg.StripeKey,
		cfg.StripeWebhookSecret,
	)
	watchHandler := handlers.NewWatchHandler(watchSessionRepo, streamRepo, userRepo, missionTriggers)
	viewerHandler := handlers.NewViewerHandler(viewerSessionRepo)
	analyticsHandler := handlers.NewAnalyticsHandler(
		raceRepo,
		viewerSessionRepo,
		watchSessionRepo,
		revenueRepo,
		streamRepo,
		playbackEventRepo,
		streamStatsRepo,
		analyticsAggregator,
		bunnyImporter,
		bunnyEnabled,
	)
	costHandler := handlers.NewCostHandler(costRepo, raceRepo)
	pollManager := chat.NewPollManager()
	chatHandler := handlers.NewChatHandler(chatRepo, raceRepo, streamRepo, userRepo, entitlementRepo, hub, rateLimiter, missionTriggers, pollManager)
	userHandler := handlers.NewUserHandler(userRepo, watchSessionRepo)
	userPrefsHandler := handlers.NewUserPreferencesHandler(userPrefsRepo)
	userFavHandler := handlers.NewUserFavoritesHandler(userFavRepo)
	watchHistoryHandler := handlers.NewWatchHistoryHandler(watchHistoryRepo)
	recommendationsHandler := handlers.NewRecommendationsHandler(recommendationService)
	missionsHandler := handlers.NewMissionsHandler(missionService)
	xpHandler := handlers.NewXPHandler(userRepo, xpService)
	weeklyHandler := handlers.NewWeeklyHandler(weeklyService)
	predictionsHandler := handlers.NewPredictionsHandler(predictionService)
	achievementsHandler := handlers.NewAchievementsHandler(achievementService)
	analyticsIngestionHandler := handlers.NewAnalyticsIngestionHandler(streamRepo, playbackEventRepo)

	// Middleware instances
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)
	userAuthMiddleware := middleware.UserAuthMiddleware(cfg.JWTSecret)
	optionalUserAuthMiddleware := middleware.OptionalUserAuthMiddleware(cfg.JWTSecret)
	chatAuthMiddleware := middleware.ChatAuthMiddleware(cfg.JWTSecret)
	csrfProtection := middleware.CSRFProtection(cfg.JWTSecret)

	// Setup route groups
	setupPublicRoutes(app, healthHandler, raceHandler, userHandler, missionsHandler)
	setupAuthRoutes(app, authHandler)
	setupViewerRoutes(app, viewerHandler, optionalUserAuthMiddleware)
	setupStreamRoutes(app, raceHandler, streamHandler, optionalUserAuthMiddleware)
	setupChatRoutes(app, chatHandler, chatAuthMiddleware, authMiddleware, userAuthMiddleware)
	setupUserRoutes(app, authHandler, paymentHandler, watchHandler, userPrefsHandler, userFavHandler, watchHistoryHandler, recommendationsHandler, missionsHandler, xpHandler, weeklyHandler, achievementsHandler, userAuthMiddleware, csrfProtection)
	setupPredictionRoutes(app, predictionsHandler, optionalUserAuthMiddleware, userAuthMiddleware, csrfProtection)
	setupWebhookRoutes(app, paymentHandler)
	setupAdminRoutes(app, adminHandler, analyticsHandler, costHandler, authMiddleware, csrfProtection)
	setupAnalyticsRoutes(app, analyticsIngestionHandler)
}

func setupPublicRoutes(app *fiber.App, healthHandler *handlers.HealthHandler, raceHandler *handlers.RaceHandler, userHandler *handlers.UserHandler, missionsHandler *handlers.MissionsHandler) {
	// Public routes with lenient rate limiting
	public := app.Group("", middleware.LenientRateLimiter())
	public.Get("/health", healthHandler.GetHealth)
	public.Get("/races", raceHandler.GetRaces)
	public.Get("/leaderboard", userHandler.GetLeaderboard)
	// Public user profile (no auth required) - uses /profiles to avoid conflict with authenticated /users group
	public.Get("/profiles/:id", userHandler.GetPublicProfile)
	// General race routes (must be after more specific routes in other groups, but here strict ordering depends on framework)
	// Fiber matches first, so specific routes should be registered before parameterized routes if they conflict.
	// /races/:id is quite generic, so ensuring it doesn't conflict is key.
	public.Get("/races/:id", raceHandler.GetRaceByID)
	// Public missions endpoint
	public.Get("/missions/active", missionsHandler.GetActiveMissions)
}

func setupAuthRoutes(app *fiber.App, authHandler *handlers.AuthHandler) {
	// Auth routes with strict rate limiting (prevent brute force)
	auth := app.Group("/auth", middleware.StrictRateLimiter())
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
}

func setupViewerRoutes(app *fiber.App, viewerHandler *handlers.ViewerHandler, optionalAuth fiber.Handler) {
	// Viewer tracking routes (public, but supports authenticated users)
	// Must be before /races/:id to avoid route conflicts
	viewer := app.Group("", middleware.LenientRateLimiter())
	viewer.Post("/viewers/sessions/start", optionalAuth, viewerHandler.StartSession)
	viewer.Post("/viewers/sessions/end", optionalAuth, viewerHandler.EndSession)
	viewer.Post("/viewers/sessions/heartbeat", optionalAuth, viewerHandler.Heartbeat)
	viewer.Get("/races/:id/viewers/concurrent", viewerHandler.GetConcurrentViewers)
	viewer.Get("/races/:id/viewers/unique", viewerHandler.GetUniqueViewers)
}

func setupStreamRoutes(app *fiber.App, raceHandler *handlers.RaceHandler, streamHandler *handlers.StreamHandler, optionalAuth fiber.Handler) {
	// Stream endpoint - optional auth (checks inside handler)
	stream := app.Group("", middleware.LenientRateLimiter())
	stream.Get("/races/:id/stream", optionalAuth, raceHandler.GetRaceStream)
	stream.Get("/races/:id/stream/status", streamHandler.GetStreamStatus)
}

func setupChatRoutes(app *fiber.App, chatHandler *handlers.ChatHandler, chatAuth fiber.Handler, adminAuth fiber.Handler, userAuth fiber.Handler) {
	// Chat routes - MUST be before /races/:id to avoid route conflicts
	chatRoutes := app.Group("", middleware.LenientRateLimiter())
	chatRoutes.Get("/races/:id/chat/ws", chatAuth, chatHandler.HandleWebSocket)
	chatRoutes.Get("/races/:id/chat/history", chatHandler.GetChatHistory)
	chatRoutes.Get("/races/:id/chat/stats", chatHandler.GetChatStats)
	chatRoutes.Post("/races/:id/chat/polls", adminAuth, chatHandler.CreatePoll)
	chatRoutes.Post("/races/:id/chat/polls/:pollId/close", adminAuth, chatHandler.ClosePoll)
	chatRoutes.Post("/races/:id/chat/polls/:pollId/vote", userAuth, chatHandler.CastPollVote)
}

func setupUserRoutes(app *fiber.App, authHandler *handlers.AuthHandler, paymentHandler *handlers.PaymentHandler, watchHandler *handlers.WatchHandler, userPrefsHandler *handlers.UserPreferencesHandler, userFavHandler *handlers.UserFavoritesHandler, watchHistoryHandler *handlers.WatchHistoryHandler, recommendationsHandler *handlers.RecommendationsHandler, missionsHandler *handlers.MissionsHandler, xpHandler *handlers.XPHandler, weeklyHandler *handlers.WeeklyHandler, achievementsHandler *handlers.AchievementsHandler, userAuth fiber.Handler, csrf fiber.Handler) {
	// Protected user routes with standard rate limiting and CSRF protection
	user := app.Group("/users", userAuth, middleware.StandardRateLimiter(), csrf)
	user.Get("/me", authHandler.GetProfile)
	user.Post("/me/password", authHandler.ChangePassword)
	user.Post("/me/points/tick", authHandler.AwardWatchPoints)  // 10 points for watching
	user.Post("/me/points/bonus", authHandler.AwardBonusPoints) // 50 points for claim bonus
	user.Post("/payments/create-checkout", paymentHandler.CreateCheckout)
	user.Post("/watch/sessions/start", watchHandler.StartSession)
	user.Post("/watch/sessions/end", watchHandler.EndSession)
	user.Get("/watch/sessions/stats/:race_id", watchHandler.GetStats)

	// User preferences routes
	user.Get("/me/preferences", userPrefsHandler.GetPreferences)
	user.Put("/me/preferences", userPrefsHandler.UpdatePreferences)
	user.Post("/me/onboarding/complete", userPrefsHandler.CompleteOnboarding)

	// User favorites routes
	user.Get("/me/favorites", userFavHandler.GetFavorites)
	user.Post("/me/favorites", userFavHandler.AddFavorite)
	user.Delete("/me/favorites/:type/:id", userFavHandler.RemoveFavorite)

	// Watch history routes
	user.Get("/me/watch-history", watchHistoryHandler.GetWatchHistory)

	// Recommendations routes
	user.Get("/me/recommendations", recommendationsHandler.GetAllRecommendations)
	user.Get("/me/recommendations/continue-watching", recommendationsHandler.GetContinueWatching)
	user.Get("/me/recommendations/upcoming", recommendationsHandler.GetUpcoming)
	user.Get("/me/recommendations/replays", recommendationsHandler.GetReplays)

	// Missions routes
	user.Get("/me/missions", missionsHandler.GetUserMissions)
	user.Get("/me/missions/career", missionsHandler.GetCareerMissions)
	user.Get("/me/missions/weekly", missionsHandler.GetWeeklyMissions)
	user.Post("/me/missions/:missionId/claim", missionsHandler.ClaimMissionReward)

	// XP routes
	user.Get("/me/xp", xpHandler.GetUserXPProgress)

	// Achievements
	user.Get("/me/achievements", achievementsHandler.GetUserAchievements)

	// Weekly routes
	user.Get("/me/weekly", weeklyHandler.GetWeeklyProgress)
	user.Post("/me/weekly/claim", weeklyHandler.ClaimWeeklyReward)
}

func setupWebhookRoutes(app *fiber.App, paymentHandler *handlers.PaymentHandler) {
	// Webhooks (no auth required, uses Stripe signature)
	// CSRF is automatically skipped for webhooks
	webhook := app.Group("/webhooks", middleware.StrictRateLimiter())
	webhook.Post("/stripe", paymentHandler.HandleWebhook)
}

func setupPredictionRoutes(app *fiber.App, predictionsHandler *handlers.PredictionsHandler, optionalAuth fiber.Handler, userAuth fiber.Handler, csrf fiber.Handler) {
	// Prediction routes - public for markets, protected for bets
	predictions := app.Group("", middleware.LenientRateLimiter())
	predictions.Get("/races/:id/predictions", optionalAuth, predictionsHandler.GetRacePredictions)

	// Protected bet routes - placed at race level (before /races/:id route to avoid conflicts)
	betRoutes := app.Group("", userAuth, middleware.StandardRateLimiter(), csrf)
	betRoutes.Post("/races/:id/predictions/:marketId/bet", predictionsHandler.PlaceBet)

	// User predictions route (in /users group)
	userPredictions := app.Group("/users", userAuth, middleware.StandardRateLimiter(), csrf)
	userPredictions.Get("/me/predictions", predictionsHandler.GetUserPredictions)
}

func setupAdminRoutes(app *fiber.App, adminHandler *handlers.AdminHandler, analyticsHandler *handlers.AnalyticsHandler, costHandler *handlers.CostHandler, adminAuth fiber.Handler, csrf fiber.Handler) {
	// Admin routes (protected) with standard rate limiting and CSRF protection
	admin := app.Group("/admin", adminAuth, middleware.StandardRateLimiter(), csrf)

	// Races
	admin.Post("/races", adminHandler.CreateRace)
	admin.Put("/races/:id", adminHandler.UpdateRace)
	admin.Delete("/races/:id", adminHandler.DeleteRace)

	// Streams
	admin.Post("/races/:id/stream", adminHandler.UpdateStream)
	admin.Put("/races/:id/stream/status", adminHandler.UpdateStreamStatus)

	// Revenue
	admin.Get("/revenue", adminHandler.GetRevenue)
	admin.Get("/revenue/races/:id", adminHandler.GetRevenueByRace)
	admin.Get("/revenue/races/:id/summary", adminHandler.GetRevenueSummaryByRace)
	admin.Post("/revenue/recalculate", adminHandler.RecalculateRevenue)
	admin.Post("/revenue/recalculate/:year/:month", adminHandler.RecalculateRevenueForPeriod)

	// Analytics
	admin.Get("/analytics/races", analyticsHandler.GetRaceAnalytics)
	admin.Get("/analytics/watch-time", analyticsHandler.GetWatchTimeAnalytics)
	admin.Get("/analytics/revenue", analyticsHandler.GetRevenueAnalytics)
	admin.Get("/analytics/streams", analyticsHandler.GetStreamAnalytics)
	admin.Get("/analytics/streams/summary", analyticsHandler.GetStreamAnalyticsSummary)
	admin.Post("/analytics/streams/bunny-sync", analyticsHandler.SyncBunnyAnalytics)

	// Costs
	admin.Post("/costs", costHandler.CreateCost)
	admin.Get("/costs", costHandler.GetCosts)
	admin.Get("/costs/summary", costHandler.GetCostSummary)
	admin.Get("/costs/:id", costHandler.GetCostByID)
	admin.Put("/costs/:id", costHandler.UpdateCost)
	admin.Delete("/costs/:id", costHandler.DeleteCost)
	admin.Get("/costs/races/:race_id", costHandler.GetCostsByRace)
}

func setupAnalyticsRoutes(app *fiber.App, analyticsHandler *handlers.AnalyticsIngestionHandler) {
	analytics := app.Group("/analytics", middleware.LenientRateLimiter())
	// Handle OPTIONS for CORS preflight
	analytics.Options("/events", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})
	analytics.Post("/events", analyticsHandler.IngestEvents)
}
