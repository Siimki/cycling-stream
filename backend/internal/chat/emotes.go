package chat

import (
	"regexp"
	"strings"
)

var specialEmotePattern = regexp.MustCompile(`(?i):(bike|fire|zap|bolt|clap|crown|rocket|heart|star|podium):`)

// IsSpecialEmoteMessage returns true when the message consists entirely of known emotes.
func IsSpecialEmoteMessage(message string) bool {
	trimmed := strings.TrimSpace(message)
	if trimmed == "" {
		return false
	}

	parts := strings.Fields(trimmed)
	if len(parts) == 0 {
		return false
	}

	for _, part := range parts {
		if !specialEmotePattern.MatchString(part) {
			return false
		}
	}
	return true
}
