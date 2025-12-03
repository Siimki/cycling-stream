package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cyclingstream/backend/internal/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestChatAuthMiddleware(t *testing.T) {
	logger.Init("test")
	secret := "test-secret"
	app := fiber.New()

	app.Get("/ws", ChatAuthMiddleware(secret), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	makeReq := func(token string, useHeader bool) *http.Response {
		req := httptest.NewRequest("GET", "/ws", nil)
		req.Header.Set("Connection", "Upgrade")
		req.Header.Set("Upgrade", "websocket")
		if token != "" && !useHeader {
			req.URL.RawQuery = "token=" + token
		}
		if token != "" && useHeader {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		resp, _ := app.Test(req)
		return resp
	}

	t.Run("allows anonymous when no token provided", func(t *testing.T) {
		resp := makeReq("", false)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("rejects invalid token format", func(t *testing.T) {
		resp := makeReq("not-a-token", true)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("rejects expired token", func(t *testing.T) {
		expired := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
			UserID: "user-1",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
			},
		})
		tokenStr, _ := expired.SignedString([]byte(secret))
		resp := makeReq(tokenStr, true)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("accepts valid token", func(t *testing.T) {
		valid := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
			UserID:  "user-1",
			IsAdmin: false,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			},
		})
		tokenStr, _ := valid.SignedString([]byte(secret))
		resp := makeReq(tokenStr, true)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}
