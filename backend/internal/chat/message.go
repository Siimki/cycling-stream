package chat

import (
	"encoding/json"
	"time"

	"github.com/cyclingstream/backend/internal/models"
)

// MessageType represents the type of WebSocket message
type MessageType string

const (
	MessageTypeMessage   MessageType = "message"
	MessageTypeError     MessageType = "error"
	MessageTypeJoined    MessageType = "joined"
	MessageTypeLeft      MessageType = "left"
	MessageTypePing      MessageType = "ping"
	MessageTypePong      MessageType = "pong"
	MessageTypeSendMessage MessageType = "send_message"
)

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

// ChatMessageData represents chat message data in WebSocket messages
type ChatMessageData struct {
	ID        string    `json:"id"`
	RaceID    string    `json:"race_id"`
	UserID    *string   `json:"user_id,omitempty"`
	Username  string    `json:"username"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// ErrorData represents error data in WebSocket messages
type ErrorData struct {
	Message string `json:"message"`
}

// UserActionData represents user join/leave data
type UserActionData struct {
	Username string `json:"username"`
}

// SendMessageData represents data sent by client to send a message
type SendMessageData struct {
	Message string `json:"message"`
}

// NewMessageWSMessage creates a WebSocket message for a chat message
func NewMessageWSMessage(msg *models.ChatMessage) *WSMessage {
	return &WSMessage{
		Type: string(MessageTypeMessage),
		Data: ChatMessageData{
			ID:        msg.ID,
			RaceID:    msg.RaceID,
			UserID:    msg.UserID,
			Username:  msg.Username,
			Message:   msg.Message,
			CreatedAt: msg.CreatedAt,
		},
	}
}

// NewErrorWSMessage creates a WebSocket error message
func NewErrorWSMessage(message string) *WSMessage {
	return &WSMessage{
		Type: string(MessageTypeError),
		Data: ErrorData{
			Message: message,
		},
	}
}

// NewJoinedWSMessage creates a WebSocket message for user joined
func NewJoinedWSMessage(username string) *WSMessage {
	return &WSMessage{
		Type: string(MessageTypeJoined),
		Data: UserActionData{
			Username: username,
		},
	}
}

// NewLeftWSMessage creates a WebSocket message for user left
func NewLeftWSMessage(username string) *WSMessage {
	return &WSMessage{
		Type: string(MessageTypeLeft),
		Data: UserActionData{
			Username: username,
		},
	}
}

// NewPongWSMessage creates a WebSocket pong message
func NewPongWSMessage() *WSMessage {
	return &WSMessage{
		Type: string(MessageTypePong),
	}
}

// UnmarshalWSMessage unmarshals JSON to WSMessage
func UnmarshalWSMessage(data []byte) (*WSMessage, error) {
	var msg WSMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// ParseSendMessageData parses SendMessageData from WSMessage
func ParseSendMessageData(msg *WSMessage) (*SendMessageData, error) {
	if msg.Type != string(MessageTypeSendMessage) {
		return nil, nil
	}

	dataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		return nil, err
	}

	var data SendMessageData
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

