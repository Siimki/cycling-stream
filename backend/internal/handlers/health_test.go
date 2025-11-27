package handlers

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_GetHealth(t *testing.T) {
	app := fiber.New()
	handler := NewHealthHandler(nil) // nil is acceptable for health check
	app.Get("/health", handler.GetHealth)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "ok")
}

