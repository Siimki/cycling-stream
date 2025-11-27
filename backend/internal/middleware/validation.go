package middleware

import (
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateEmail checks if an email is valid
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	// Basic email validation
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	// Check for domain
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}

// ValidatePassword checks if a password meets requirements
func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}
	if len(password) > 128 {
		return false, "Password must be less than 128 characters"
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return false, "Password must contain at least one uppercase letter"
	}
	if !hasLower {
		return false, "Password must contain at least one lowercase letter"
	}
	if !hasNumber {
		return false, "Password must contain at least one number"
	}
	if !hasSpecial {
		return false, "Password must contain at least one special character"
	}

	return true, ""
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(s string, maxLength int) string {
	// Remove null bytes and control characters
	s = strings.ReplaceAll(s, "\x00", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimSpace(s)
	
	// Limit length
	if maxLength > 0 && len(s) > maxLength {
		s = s[:maxLength]
	}
	
	return s
}

// ValidateUUID checks if a string is a valid UUID
func ValidateUUID(id string) bool {
	if len(id) != 36 {
		return false
	}
	// Basic UUID format check (8-4-4-4-12)
	parts := strings.Split(id, "-")
	if len(parts) != 5 {
		return false
	}
	if len(parts[0]) != 8 || len(parts[1]) != 4 || len(parts[2]) != 4 || len(parts[3]) != 4 || len(parts[4]) != 12 {
		return false
	}
	return true
}

// ValidateRequest validates common request fields
func ValidateRequest(c *fiber.Ctx) error {
	// This is a helper that can be used in handlers
	// For now, we'll use it as a reference
	return c.Next()
}

// ValidateJSONBody validates that the request has a valid JSON body
func ValidateJSONBody(c *fiber.Ctx) error {
	if c.Get("Content-Type") != "application/json" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Content-Type must be application/json",
		})
	}
	return c.Next()
}

