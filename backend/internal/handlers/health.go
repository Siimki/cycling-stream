package handlers

import (
	"runtime"
	"time"

	"github.com/cyclingstream/backend/internal/database"
	"github.com/gofiber/fiber/v2"
)

var (
	startTime = time.Now()
)

type HealthHandler struct {
	db *database.DB
}

func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) GetHealth(c *fiber.Ctx) error {
	start := time.Now()
	
	// Check database connectivity
	dbStatus := "ok"
	dbLatency := int64(0)
	if h.db != nil {
		dbStart := time.Now()
		err := h.db.Ping()
		dbLatency = time.Since(dbStart).Milliseconds()
		if err != nil {
			dbStatus = "error"
		}
	} else {
		dbStatus = "not_configured"
	}

	// Calculate uptime
	uptime := time.Since(startTime)

	// Get system info
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Calculate response time
	responseTime := time.Since(start).Milliseconds()

	// Determine overall status
	overallStatus := "ok"
	if dbStatus == "error" {
		overallStatus = "degraded"
	}

	statusCode := fiber.StatusOK
	if dbStatus == "error" {
		statusCode = fiber.StatusServiceUnavailable
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"status": overallStatus,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime_seconds": int64(uptime.Seconds()),
		"response_time_ms": responseTime,
		"services": fiber.Map{
			"database": fiber.Map{
				"status": dbStatus,
				"latency_ms": dbLatency,
			},
		},
		"system": fiber.Map{
			"go_version": runtime.Version(),
			"num_goroutines": runtime.NumGoroutine(),
			"memory": fiber.Map{
				"allocated_mb": bToMb(m.Alloc),
				"total_allocated_mb": bToMb(m.TotalAlloc),
				"sys_mb": bToMb(m.Sys),
				"num_gc": m.NumGC,
			},
		},
	})
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

