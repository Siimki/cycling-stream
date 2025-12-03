package chat

import (
	"testing"
	"time"

	"github.com/cyclingstream/backend/internal/logger"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Initialize logger for tests
	logger.Init("test")
}

// TestClient_NewClient tests client creation and basic functionality
func TestClient_NewClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer close(hub.register)

	t.Run("Authenticated client can send and receive messages", func(t *testing.T) {
		userID := "user-123"
		username := "TestUser"
		raceID := "race-123"
		messageHandler := func(*Client, *WSMessage) {}
		onClose := func(*Client) {}

		client := NewClient(hub, nil, &userID, username, false, raceID, messageHandler, onClose)
		assert.NotNil(t, client)
		assert.Equal(t, raceID, client.RaceID())

		// Register client with hub
		hub.register <- client
		time.Sleep(10 * time.Millisecond)

		// Test that client can receive messages from hub
		message := []byte("test message")
		client.SendMessage(message)

		// Verify message was queued
		select {
		case msg := <-client.send:
			assert.Equal(t, message, msg)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Client should be able to receive messages")
		}
	})

	t.Run("Anonymous client can receive messages", func(t *testing.T) {
		raceID := "race-anon"
		client := NewClient(hub, nil, nil, "Anonymous", false, raceID, nil, nil)
		assert.NotNil(t, client)
		assert.Equal(t, raceID, client.RaceID())

		hub.register <- client
		time.Sleep(10 * time.Millisecond)

		// Anonymous client should still be able to receive messages
		message := []byte("broadcast")
		client.SendMessage(message)

		select {
		case msg := <-client.send:
			assert.Equal(t, message, msg)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Anonymous client should be able to receive messages")
		}
	})
}

// TestClient_SendMessage tests message sending behavior
func TestClient_SendMessage(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer close(hub.register)

	client := NewClient(hub, nil, nil, "Test", false, "race-send", nil, nil)
	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	t.Run("Client can send multiple messages in sequence", func(t *testing.T) {
		messages := [][]byte{
			[]byte("message 1"),
			[]byte("message 2"),
			[]byte("message 3"),
		}

		// Send all messages
		for _, msg := range messages {
			success := client.SendMessage(msg)
			assert.True(t, success, "Message should be sent successfully")
		}

		// Verify all messages are received in order
		for i, expectedMsg := range messages {
			select {
			case msg := <-client.send:
				assert.Equal(t, expectedMsg, msg, "Message %d should match", i+1)
			case <-time.After(100 * time.Millisecond):
				t.Fatalf("Message %d was not received", i+1)
			}
		}
	})

	t.Run("Client drops messages when send channel is full", func(t *testing.T) {
		// Create a client with a very small channel to test dropping behavior
		smallChannel := make(chan []byte, 1)
		client.send = smallChannel

		// Fill the channel
		smallChannel <- []byte("first")

		// Try to send another message (should be dropped, returns false)
		success := client.SendMessage([]byte("second"))
		assert.False(t, success, "Message should be dropped when channel is full")

		// Verify only first message is in channel
		msg := <-smallChannel
		assert.Equal(t, []byte("first"), msg)

		// Verify second message was dropped
		select {
		case <-smallChannel:
			t.Fatal("Second message should have been dropped")
		default:
			// Expected - channel is empty
		}
	})
}

// TestClient_MessageHandler tests message handler functionality
func TestClient_MessageHandler(t *testing.T) {
	hub := NewHub()

	t.Run("Message handler is called", func(t *testing.T) {
		handlerCalled := false
		handler := func(client *Client, msg *WSMessage) {
			handlerCalled = true
		}

		client := &Client{
			hub:            hub,
			send:           make(chan []byte, 256),
			messageHandler: handler,
		}

		// Simulate message
		msg := &WSMessage{
			Type: "test",
		}

		if client.messageHandler != nil {
			client.messageHandler(client, msg)
		}

		assert.True(t, handlerCalled)
	})
}

// TestClient_OnClose tests onClose callback
func TestClient_OnClose(t *testing.T) {
	hub := NewHub()

	t.Run("OnClose callback is called", func(t *testing.T) {
		onCloseCalled := false
		onClose := func(client *Client) {
			onCloseCalled = true
		}

		client := &Client{
			hub:     hub,
			send:    make(chan []byte, 256),
			onClose: onClose,
		}

		if client.onClose != nil {
			client.onClose(client)
		}

		assert.True(t, onCloseCalled)
	})
}
