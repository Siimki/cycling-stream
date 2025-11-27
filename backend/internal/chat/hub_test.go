package chat

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// createTestClient creates a mock client for testing
func createTestClient(hub *Hub, userID *string, username string) *Client {
	return &Client{
		hub:      hub,
		conn:     nil, // Mock connection
		send:     make(chan []byte, 256),
		userID:   userID,
		username: username,
		isAdmin:  false,
	}
}

// TestHub_RegisterClient tests client registration
func TestHub_RegisterClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer close(hub.register)

	t.Run("Hub initializes correctly", func(t *testing.T) {
		assert.NotNil(t, hub)
		assert.NotNil(t, hub.clients)
		assert.NotNil(t, hub.rooms)
		assert.NotNil(t, hub.register)
		assert.NotNil(t, hub.unregister)
	})

	t.Run("Hub has correct channel sizes", func(t *testing.T) {
		// Verify channels are created
		assert.NotNil(t, hub.broadcast)
		assert.NotNil(t, hub.joinRoom)
		assert.NotNil(t, hub.leaveRoom)
	})

	t.Run("Register client adds to hub", func(t *testing.T) {
		client := createTestClient(hub, nil, "TestUser")
		hub.register <- client

		// Give hub time to process
		time.Sleep(10 * time.Millisecond)

		hub.mu.RLock()
		_, exists := hub.clients[client]
		hub.mu.RUnlock()
		assert.True(t, exists, "Client should be registered")
	})

	t.Run("Unregister client removes from hub", func(t *testing.T) {
		client := createTestClient(hub, nil, "TestUser")
		hub.register <- client
		time.Sleep(10 * time.Millisecond)

		hub.unregister <- client
		time.Sleep(10 * time.Millisecond)

		hub.mu.RLock()
		_, exists := hub.clients[client]
		hub.mu.RUnlock()
		assert.False(t, exists, "Client should be unregistered")
	})
}

// TestHub_JoinRoom tests room joining
func TestHub_JoinRoom(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer close(hub.register)

	t.Run("Room is created when first client joins", func(t *testing.T) {
		raceID := "test-race-1"
		client := createTestClient(hub, nil, "TestUser")

		// Register client first
		hub.register <- client
		time.Sleep(10 * time.Millisecond)

		// Join room
		action := &RoomAction{Client: client, RaceID: raceID}
		hub.joinRoom <- action
		time.Sleep(10 * time.Millisecond)

		// Verify room was created and client is in it
		hub.mu.RLock()
		roomClients, exists := hub.rooms[raceID]
		hub.mu.RUnlock()
		assert.True(t, exists, "Room should exist")
		assert.NotNil(t, roomClients)
		_, inRoom := roomClients[client]
		assert.True(t, inRoom, "Client should be in room")
	})

	t.Run("Multiple clients can join same room", func(t *testing.T) {
		raceID := "test-race-2"
		client1 := createTestClient(hub, nil, "User1")
		client2 := createTestClient(hub, nil, "User2")
		client3 := createTestClient(hub, nil, "User3")

		// Register all clients
		hub.register <- client1
		hub.register <- client2
		hub.register <- client3
		time.Sleep(10 * time.Millisecond)

		// Join all to same room
		hub.joinRoom <- &RoomAction{Client: client1, RaceID: raceID}
		hub.joinRoom <- &RoomAction{Client: client2, RaceID: raceID}
		hub.joinRoom <- &RoomAction{Client: client3, RaceID: raceID}
		time.Sleep(10 * time.Millisecond)

		// Verify all clients are in room
		count := hub.GetRoomClientCount(raceID)
		assert.Equal(t, 3, count, "Room should have 3 clients")
	})

	t.Run("Client can join multiple rooms", func(t *testing.T) {
		raceID1 := "test-race-3"
		raceID2 := "test-race-4"
		client := createTestClient(hub, nil, "MultiRoomUser")

		hub.register <- client
		time.Sleep(10 * time.Millisecond)

		hub.joinRoom <- &RoomAction{Client: client, RaceID: raceID1}
		hub.joinRoom <- &RoomAction{Client: client, RaceID: raceID2}
		time.Sleep(10 * time.Millisecond)

		assert.Equal(t, 1, hub.GetRoomClientCount(raceID1))
		assert.Equal(t, 1, hub.GetRoomClientCount(raceID2))
	})
}

// TestHub_BroadcastToRoom tests room broadcasting
func TestHub_BroadcastToRoom(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer close(hub.register)

	t.Run("Broadcast to empty room does nothing", func(t *testing.T) {
		raceID := "test-race-3"
		message := []byte("test message")

		// Should not panic
		hub.BroadcastToRoom(raceID, message)
	})

	t.Run("Broadcast to non-existent room does nothing", func(t *testing.T) {
		raceID := "non-existent-race"
		message := []byte("test message")

		// Should not panic
		hub.BroadcastToRoom(raceID, message)
	})

	t.Run("Broadcast to room with clients sends message", func(t *testing.T) {
		raceID := "test-race-5"
		client1 := createTestClient(hub, nil, "User1")
		client2 := createTestClient(hub, nil, "User2")

		hub.register <- client1
		hub.register <- client2
		time.Sleep(10 * time.Millisecond)

		hub.joinRoom <- &RoomAction{Client: client1, RaceID: raceID}
		hub.joinRoom <- &RoomAction{Client: client2, RaceID: raceID}
		time.Sleep(10 * time.Millisecond)

		message := []byte("broadcast message")
		hub.BroadcastToRoom(raceID, message)

		// Verify both clients received the message
		time.Sleep(10 * time.Millisecond)
		select {
		case msg := <-client1.send:
			assert.Equal(t, message, msg)
		case <-time.After(100 * time.Millisecond):
			t.Error("Client1 did not receive message")
		}

		select {
		case msg := <-client2.send:
			assert.Equal(t, message, msg)
		case <-time.After(100 * time.Millisecond):
			t.Error("Client2 did not receive message")
		}
	})

	t.Run("Broadcast only sends to clients in same room", func(t *testing.T) {
		raceID1 := "test-race-6"
		raceID2 := "test-race-7"
		client1 := createTestClient(hub, nil, "User1")
		client2 := createTestClient(hub, nil, "User2")

		hub.register <- client1
		hub.register <- client2
		time.Sleep(10 * time.Millisecond)

		hub.joinRoom <- &RoomAction{Client: client1, RaceID: raceID1}
		hub.joinRoom <- &RoomAction{Client: client2, RaceID: raceID2}
		time.Sleep(10 * time.Millisecond)

		message := []byte("room-specific message")
		hub.BroadcastToRoom(raceID1, message)

		// Client1 should receive, client2 should not
		time.Sleep(10 * time.Millisecond)
		select {
		case msg := <-client1.send:
			assert.Equal(t, message, msg)
		case <-time.After(100 * time.Millisecond):
			t.Error("Client1 did not receive message")
		}

		select {
		case <-client2.send:
			t.Error("Client2 should not receive message from different room")
		case <-time.After(100 * time.Millisecond):
			// Expected - client2 should not receive
		}
	})
}

// TestHub_GetRoomClientCount tests room client counting
func TestHub_GetRoomClientCount(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer close(hub.register)

	t.Run("Empty room returns zero", func(t *testing.T) {
		raceID := "test-race-4"
		count := hub.GetRoomClientCount(raceID)
		assert.Equal(t, 0, count)
	})

	t.Run("Non-existent room returns zero", func(t *testing.T) {
		raceID := "non-existent-race"
		count := hub.GetRoomClientCount(raceID)
		assert.Equal(t, 0, count)
	})

	t.Run("Returns correct count for room with clients", func(t *testing.T) {
		raceID := "test-race-8"
		client1 := createTestClient(hub, nil, "User1")
		client2 := createTestClient(hub, nil, "User2")

		hub.register <- client1
		hub.register <- client2
		time.Sleep(10 * time.Millisecond)

		hub.joinRoom <- &RoomAction{Client: client1, RaceID: raceID}
		time.Sleep(10 * time.Millisecond)
		assert.Equal(t, 1, hub.GetRoomClientCount(raceID))

		hub.joinRoom <- &RoomAction{Client: client2, RaceID: raceID}
		time.Sleep(10 * time.Millisecond)
		assert.Equal(t, 2, hub.GetRoomClientCount(raceID))
	})
}

// TestHub_ConcurrentAccess tests thread safety
func TestHub_ConcurrentAccess(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer close(hub.register)

	t.Run("Concurrent room operations", func(t *testing.T) {
		raceID := "test-race-9"
		numClients := 10
		clients := make([]*Client, numClients)

		// Create and register clients
		for i := 0; i < numClients; i++ {
			clients[i] = createTestClient(hub, nil, "User")
			hub.register <- clients[i]
		}
		time.Sleep(10 * time.Millisecond)

		// Concurrently join rooms
		var wg sync.WaitGroup
		for i := 0; i < numClients; i++ {
			wg.Add(1)
			go func(client *Client) {
				defer wg.Done()
				hub.joinRoom <- &RoomAction{Client: client, RaceID: raceID}
			}(clients[i])
		}
		wg.Wait()
		time.Sleep(10 * time.Millisecond)

		// Verify all clients joined
		count := hub.GetRoomClientCount(raceID)
		assert.Equal(t, numClients, count, "All clients should be in room")

		// Concurrently leave rooms
		for i := 0; i < numClients; i++ {
			wg.Add(1)
			go func(client *Client) {
				defer wg.Done()
				hub.leaveRoom <- &RoomAction{Client: client, RaceID: raceID}
			}(clients[i])
		}
		wg.Wait()
		time.Sleep(10 * time.Millisecond)

		// Verify room is cleaned up
		count = hub.GetRoomClientCount(raceID)
		assert.Equal(t, 0, count, "Room should be empty")
	})
}

// TestHub_Cleanup tests cleanup of empty rooms
func TestHub_Cleanup(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer close(hub.register)

	t.Run("Empty rooms are cleaned up when last client leaves", func(t *testing.T) {
		raceID := "test-race-10"
		client := createTestClient(hub, nil, "User")

		hub.register <- client
		time.Sleep(10 * time.Millisecond)

		// Join room
		hub.joinRoom <- &RoomAction{Client: client, RaceID: raceID}
		time.Sleep(10 * time.Millisecond)

		// Verify room exists
		hub.mu.RLock()
		_, exists := hub.rooms[raceID]
		hub.mu.RUnlock()
		assert.True(t, exists, "Room should exist")

		// Leave room
		hub.leaveRoom <- &RoomAction{Client: client, RaceID: raceID}
		time.Sleep(10 * time.Millisecond)

		// Verify room is cleaned up
		hub.mu.RLock()
		_, exists = hub.rooms[raceID]
		hub.mu.RUnlock()
		assert.False(t, exists, "Empty room should be cleaned up")
	})

	t.Run("Room persists when clients remain", func(t *testing.T) {
		raceID := "test-race-11"
		client1 := createTestClient(hub, nil, "User1")
		client2 := createTestClient(hub, nil, "User2")

		hub.register <- client1
		hub.register <- client2
		time.Sleep(10 * time.Millisecond)

		hub.joinRoom <- &RoomAction{Client: client1, RaceID: raceID}
		hub.joinRoom <- &RoomAction{Client: client2, RaceID: raceID}
		time.Sleep(10 * time.Millisecond)

		// Leave one client
		hub.leaveRoom <- &RoomAction{Client: client1, RaceID: raceID}
		time.Sleep(10 * time.Millisecond)

		// Room should still exist
		hub.mu.RLock()
		_, exists := hub.rooms[raceID]
		hub.mu.RUnlock()
		assert.True(t, exists, "Room should persist with remaining client")
		assert.Equal(t, 1, hub.GetRoomClientCount(raceID))
	})

	t.Run("Unregister client cleans up from all rooms", func(t *testing.T) {
		raceID1 := "test-race-12"
		raceID2 := "test-race-13"
		client := createTestClient(hub, nil, "User")

		hub.register <- client
		time.Sleep(10 * time.Millisecond)

		hub.joinRoom <- &RoomAction{Client: client, RaceID: raceID1}
		hub.joinRoom <- &RoomAction{Client: client, RaceID: raceID2}
		time.Sleep(10 * time.Millisecond)

		// Unregister should remove from all rooms
		hub.unregister <- client
		time.Sleep(10 * time.Millisecond)

		// Both rooms should be cleaned up
		assert.Equal(t, 0, hub.GetRoomClientCount(raceID1))
		assert.Equal(t, 0, hub.GetRoomClientCount(raceID2))
	})
}

// TestHub_Run tests hub main loop
func TestHub_Run(t *testing.T) {
	hub := NewHub()

	// Start hub in background
	go hub.Run()

	// Give it a moment to start
	time.Sleep(10 * time.Millisecond)

	// Verify hub is running (channels are ready)
	client := createTestClient(hub, nil, "Test")
	select {
	case hub.register <- client:
		// Channel is open, hub is running
		t.Log("Hub is running")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Hub register channel is not ready")
	}

	// Verify client was registered
	time.Sleep(10 * time.Millisecond)
	hub.mu.RLock()
	_, exists := hub.clients[client]
	hub.mu.RUnlock()
	assert.True(t, exists, "Client should be registered")

	// Cleanup: unregister client
	// Note: We don't close channels as that would cause hub.Run() to panic
	// The hub goroutine will continue running, which is acceptable for a test
	hub.unregister <- client
	time.Sleep(10 * time.Millisecond)
}
