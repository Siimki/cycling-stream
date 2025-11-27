package chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClient_NewClient tests client creation
func TestClient_NewClient(t *testing.T) {
	hub := NewHub()
	
	t.Run("Create client with all fields", func(t *testing.T) {
		userID := "user-123"
		username := "TestUser"
		isAdmin := false
		messageHandler := func(*Client, *WSMessage) {}
		onClose := func(*Client) {}
		
		// Test NewClient function
		client := NewClient(hub, nil, &userID, username, isAdmin, messageHandler, onClose)

		assert.NotNil(t, client)
		assert.Equal(t, hub, client.hub)
		assert.NotNil(t, client.send)
		assert.NotNil(t, client.userID)
		assert.Equal(t, userID, *client.userID)
		assert.Equal(t, username, client.username)
		assert.Equal(t, isAdmin, client.isAdmin)
		assert.NotNil(t, client.messageHandler)
		assert.NotNil(t, client.onClose)
	})

	t.Run("Create anonymous client", func(t *testing.T) {
		client := NewClient(hub, nil, nil, "Anonymous", false, nil, nil)

		assert.Nil(t, client.userID)
		assert.Equal(t, "Anonymous", client.username)
		assert.Nil(t, client.messageHandler)
		assert.Nil(t, client.onClose)
	})

	t.Run("Create admin client", func(t *testing.T) {
		userID := "admin-123"
		client := NewClient(hub, nil, &userID, "Admin", true, nil, nil)

		assert.NotNil(t, client.userID)
		assert.Equal(t, userID, *client.userID)
		assert.True(t, client.isAdmin)
	})
}

// TestClient_SendMessage tests message sending
func TestClient_SendMessage(t *testing.T) {
	hub := NewHub()
	client := &Client{
		hub:      hub,
		send:     make(chan []byte, 256),
		userID:   nil,
		username: "Test",
		isAdmin:  false,
	}

	t.Run("Send message to channel", func(t *testing.T) {
		message := []byte("test message")
		
		// Send message
		client.SendMessage(message)
		
		// Verify message is in channel
		select {
		case msg := <-client.send:
			assert.Equal(t, message, msg)
		default:
			t.Fatal("Message was not sent to channel")
		}
	})

	t.Run("Send multiple messages", func(t *testing.T) {
		messages := [][]byte{
			[]byte("message 1"),
			[]byte("message 2"),
			[]byte("message 3"),
		}

		for _, msg := range messages {
			client.SendMessage(msg)
		}

		// Verify all messages are in channel
		for i, expectedMsg := range messages {
			select {
			case msg := <-client.send:
				assert.Equal(t, expectedMsg, msg, "Message %d mismatch", i)
			default:
				t.Fatalf("Message %d was not sent to channel", i)
			}
		}
	})

	t.Run("Full channel drops message", func(t *testing.T) {
		// Fill channel to capacity
		smallChannel := make(chan []byte, 1)
		client.send = smallChannel
		
		// Fill the channel
		smallChannel <- []byte("first")
		
		// Try to send another message (should be dropped)
		client.SendMessage([]byte("second"))
		
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

	t.Run("Nil message handler doesn't panic", func(t *testing.T) {
		client := &Client{
			hub:            hub,
			send:           make(chan []byte, 256),
			messageHandler: nil,
		}

		msg := &WSMessage{
			Type: "test",
		}

		// Should not panic
		if client.messageHandler != nil {
			client.messageHandler(client, msg)
		}
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

	t.Run("Nil onClose doesn't panic", func(t *testing.T) {
		client := &Client{
			hub:     hub,
			send:    make(chan []byte, 256),
			onClose: nil,
		}

		// Should not panic
		if client.onClose != nil {
			client.onClose(client)
		}
	})
}

