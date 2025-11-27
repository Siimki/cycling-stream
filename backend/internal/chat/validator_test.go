package chat

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMessage(t *testing.T) {
	t.Run("Valid message passes", func(t *testing.T) {
		message := "Hello, world!"
		result, err := ValidateMessage(message)
		assert.NoError(t, err)
		assert.Equal(t, message, result)
	})

	t.Run("Empty message fails", func(t *testing.T) {
		message := ""
		_, err := ValidateMessage(message)
		assert.Error(t, err)
		assert.Equal(t, ErrMessageEmpty, err)
	})

	t.Run("Whitespace-only message fails", func(t *testing.T) {
		message := "   \n\t  "
		_, err := ValidateMessage(message)
		assert.Error(t, err)
		assert.Equal(t, ErrMessageEmpty, err)
	})

	t.Run("Message too long fails", func(t *testing.T) {
		message := strings.Repeat("a", MaxMessageLength+1)
		_, err := ValidateMessage(message)
		assert.Error(t, err)
		assert.Equal(t, ErrMessageTooLong, err)
	})

	t.Run("Message at max length passes", func(t *testing.T) {
		message := strings.Repeat("a", MaxMessageLength)
		result, err := ValidateMessage(message)
		assert.NoError(t, err)
		assert.Equal(t, message, result)
	})

	t.Run("Message with leading/trailing whitespace is trimmed", func(t *testing.T) {
		message := "  Hello, world!  "
		result, err := ValidateMessage(message)
		assert.NoError(t, err)
		assert.Equal(t, "Hello, world!", result)
	})

	t.Run("Message with newlines passes", func(t *testing.T) {
		message := "Line 1\nLine 2"
		result, err := ValidateMessage(message)
		assert.NoError(t, err)
		assert.Equal(t, message, result)
	})

	t.Run("Message with tabs passes", func(t *testing.T) {
		message := "Tab\there"
		result, err := ValidateMessage(message)
		assert.NoError(t, err)
		assert.Equal(t, message, result)
	})

	t.Run("Control characters are removed", func(t *testing.T) {
		message := "Hello\x00World\x01Test"
		result, err := ValidateMessage(message)
		assert.NoError(t, err)
		assert.Equal(t, "HelloWorldTest", result)
	})

	t.Run("Single character message passes", func(t *testing.T) {
		message := "a"
		result, err := ValidateMessage(message)
		assert.NoError(t, err)
		assert.Equal(t, message, result)
	})
}

