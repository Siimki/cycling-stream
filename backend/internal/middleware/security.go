package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

// SecurityHeaders adds security headers to responses
func SecurityHeaders() fiber.Handler {
	return helmet.New(helmet.Config{
		XSSProtection:             "1; mode=block",
		ContentTypeNosniff:        "nosniff",
		XFrameOptions:             "DENY",
		ReferrerPolicy:            "no-referrer",
		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "cross-origin",
		OriginAgentCluster:        "?1",
		XDNSPrefetchControl:       "off",
		XDownloadOptions:          "noopen",
		XPermittedCrossDomain:     "none",
		// Content Security Policy - adjust based on your needs
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' https://api.stripe.com; frame-src https://js.stripe.com https://hooks.stripe.com;",
		// HSTS - only enable in production with HTTPS
		HSTSMaxAge:            31536000, // 1 year
		HSTSExcludeSubdomains: false,
		HSTSPreloadEnabled:    false,
	})
}

