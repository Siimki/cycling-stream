package chat

import (
	"strings"
	"unicode"
)

const (
	MinMessageLength = 1
	MaxMessageLength = 500
)

// ValidateMessage validates a chat message
func ValidateMessage(message string) (string, error) {
	// Trim whitespace
	trimmed := strings.TrimSpace(message)

	// Check if message is empty or only whitespace first
	if len(trimmed) == 0 || strings.TrimFunc(trimmed, unicode.IsSpace) == "" {
		return "", ErrMessageEmpty
	}

	// Check length
	if len(trimmed) < MinMessageLength {
		return "", ErrMessageTooShort
	}

	if len(trimmed) > MaxMessageLength {
		return "", ErrMessageTooLong
	}

	// Basic sanitization - remove control characters except newline and tab
	var sanitized strings.Builder
	for _, r := range trimmed {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			continue
		}
		sanitized.WriteRune(r)
	}

	result := sanitized.String()

	// Final check after sanitization
	if len(strings.TrimSpace(result)) < MinMessageLength {
		return "", ErrMessageEmpty
	}

	return result, nil
}

// Validation errors
var (
	ErrMessageTooShort = &ValidationError{Message: "Message is too short"}
	ErrMessageTooLong  = &ValidationError{Message: "Message is too long (max 500 characters)"}
	ErrMessageEmpty    = &ValidationError{Message: "Message cannot be empty"}
)

// ValidationError represents a message validation error
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

