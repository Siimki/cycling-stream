package testutil

import (
	"encoding/json"
)

// WSMessage represents a WebSocket message for testing
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

// Note: For full WebSocket integration tests, we recommend using github.com/gorilla/websocket
// as a test dependency. This provides better client support for connecting to Fiber WebSocket servers.
//
// To add it: go get github.com/gorilla/websocket
//
// Example usage with gorilla/websocket:
//
//	import "github.com/gorilla/websocket"
//
//	dialer := websocket.Dialer{}
//	conn, _, err := dialer.Dial("ws://localhost:8080/races/test/chat/ws", nil)
//	// ... use conn for testing

// CreateSendMessage creates a send_message WebSocket message
func CreateSendMessage(text string) WSMessage {
	return WSMessage{
		Type: "send_message",
		Data: map[string]string{
			"message": text,
		},
	}
}

// CreatePingMessage creates a ping message
func CreatePingMessage() WSMessage {
	return WSMessage{
		Type: "ping",
	}
}

// ParseWSMessage parses a JSON WebSocket message
func ParseWSMessage(data []byte) (WSMessage, error) {
	var msg WSMessage
	err := json.Unmarshal(data, &msg)
	return msg, err
}

// MarshalWSMessage marshals a WebSocket message to JSON
func MarshalWSMessage(msg WSMessage) ([]byte, error) {
	return json.Marshal(msg)
}

