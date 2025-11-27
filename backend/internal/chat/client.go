package chat

import (
	"encoding/json"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/cyclingstream/backend/internal/logger"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// MessageHandler is a function that handles incoming messages from a client
type MessageHandler func(*Client, *WSMessage)

// Client is a middleman between the websocket connection and the hub
type Client struct {
	hub            *Hub
	conn           *websocket.Conn
	send           chan []byte
	userID         *string
	username       string
	isAdmin        bool
	messageHandler MessageHandler
	onClose        func(*Client)
}

// NewClient creates a new Client
func NewClient(hub *Hub, conn *websocket.Conn, userID *string, username string, isAdmin bool, messageHandler MessageHandler, onClose func(*Client)) *Client {
	return &Client{
		hub:            hub,
		conn:           conn,
		send:           make(chan []byte, 256),
		userID:         userID,
		username:       username,
		isAdmin:        isAdmin,
		messageHandler: messageHandler,
		onClose:        onClose,
	}
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		if c.onClose != nil {
			c.onClose(c)
		}
		c.hub.unregister <- c
		c.conn.Close()
	}()

	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	c.conn.SetReadLimit(maxMessageSize)

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.WithError(err).Error("WebSocket read error")
			}
			break
		}

		// Parse the message
		var msg WSMessage
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			logger.WithError(err).Error("Failed to unmarshal WebSocket message")
			errorMsg := NewErrorWSMessage("Invalid message format")
			if errorBytes, err := json.Marshal(errorMsg); err == nil {
				c.send <- errorBytes
			}
			continue
		}

		// Handle ping messages
		if msg.Type == string(MessageTypePing) {
			pongMsg := NewPongWSMessage()
			if pongBytes, err := json.Marshal(pongMsg); err == nil {
				c.send <- pongBytes
			}
			continue
		}

		// Forward other messages to the message handler
		if c.messageHandler != nil {
			c.messageHandler(c, &msg)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write([]byte{'\n'})
				_, _ = w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Start starts the client's read and write pumps
func (c *Client) Start() {
	// Register client with hub
	c.hub.register <- c

	// Start read and write pumps
	go c.writePump()
	c.readPump()
}

// SendMessage sends a message to the client
// Returns true if the message was successfully queued, false if the channel is full
func (c *Client) SendMessage(message []byte) bool {
	select {
	case c.send <- message:
		return true
	default:
		// Channel is full, message will be dropped
		// Log this as a warning since it indicates the client is not reading fast enough
		logger.WithFields(map[string]interface{}{
			"username": c.username,
			"user_id":  c.userID,
		}).Warn("Client send channel full, message dropped")
		return false
	}
}

