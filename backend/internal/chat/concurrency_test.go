package chat

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestRateLimiter_ConcurrentAccess tests rate limiter under concurrent access
func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	rl := NewRateLimiter()
	defer rl.Stop()

	t.Run("Concurrent rate limit checks are thread-safe", func(t *testing.T) {
		identifier := "concurrent-user"
		numGoroutines := 20
		messagesPerGoroutine := 5

		var wg sync.WaitGroup
		allowedCount := 0
		var mu sync.Mutex

		// Concurrently check rate limit
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < messagesPerGoroutine; j++ {
					allowed := rl.CheckRateLimit(identifier)
					if allowed {
						mu.Lock()
						allowedCount++
						mu.Unlock()
					}
					// Small delay to increase chance of race conditions
					time.Sleep(time.Millisecond)
				}
			}()
		}

		wg.Wait()

		// Should not exceed the rate limit
		assert.LessOrEqual(t, allowedCount, ChatRateLimitMaxMessages,
			"Concurrent access should respect rate limit, got %d allowed messages", allowedCount)
	})

	t.Run("Multiple identifiers have separate rate limits concurrently", func(t *testing.T) {
		numIdentifiers := 10
		messagesPerIdentifier := ChatRateLimitMaxMessages

		var wg sync.WaitGroup
		successCount := 0
		var mu sync.Mutex

		// Each identifier should be able to send max messages
		for i := 0; i < numIdentifiers; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				identifier := "user-" + string(rune(id))
				localSuccess := 0
				for j := 0; j < messagesPerIdentifier; j++ {
					if rl.CheckRateLimit(identifier) {
						localSuccess++
					}
				}
				mu.Lock()
				successCount += localSuccess
				mu.Unlock()
			}(i)
		}

		wg.Wait()

		// Each identifier should have been able to send max messages
		expectedSuccess := numIdentifiers * ChatRateLimitMaxMessages
		assert.Equal(t, expectedSuccess, successCount,
			"Each identifier should have separate rate limit, expected %d, got %d", expectedSuccess, successCount)
	})

	t.Run("GetRemainingMessages is thread-safe under concurrent access", func(t *testing.T) {
		identifier := "remaining-user"
		numGoroutines := 10

		// Fill up the limit first
		for i := 0; i < ChatRateLimitMaxMessages; i++ {
			rl.CheckRateLimit(identifier)
		}

		var wg sync.WaitGroup
		remainingValues := make([]int, numGoroutines)

		// Concurrently get remaining messages
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				remainingValues[idx] = rl.GetRemainingMessages(identifier)
			}(i)
		}

		wg.Wait()

		// All should return 0 (limit is full)
		for i, remaining := range remainingValues {
			assert.Equal(t, 0, remaining, "Goroutine %d should return 0 remaining", i)
		}
	})
}

// TestHub_ConcurrentBroadcast tests hub broadcasting under concurrent access
func TestHub_ConcurrentBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer close(hub.register)

	t.Run("Broadcast during client registration/unregistration", func(t *testing.T) {
		raceID := "concurrent-race"
		numClients := 10
		clients := make([]*Client, numClients)

		// Create clients
		for i := 0; i < numClients; i++ {
			clients[i] = createTestClient(hub, nil, "User")
		}

		var wg sync.WaitGroup

		// Concurrently register and join room
		for i := 0; i < numClients; i++ {
			wg.Add(1)
			go func(client *Client) {
				defer wg.Done()
				hub.register <- client
				time.Sleep(5 * time.Millisecond)
				hub.joinRoom <- &RoomAction{Client: client, RaceID: raceID}
			}(clients[i])
		}

		// Concurrently broadcast while clients are joining
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 5; i++ {
				hub.BroadcastToRoom(raceID, []byte("broadcast message"))
				time.Sleep(10 * time.Millisecond)
			}
		}()

		wg.Wait()
		time.Sleep(20 * time.Millisecond)

		// Verify all clients are in room
		count := hub.GetRoomClientCount(raceID)
		assert.Equal(t, numClients, count, "All clients should be in room despite concurrent operations")
	})

	t.Run("Concurrent room joins and leaves", func(t *testing.T) {
		raceID := "concurrent-race-2"
		numClients := 15
		clients := make([]*Client, numClients)

		// Create and register clients
		for i := 0; i < numClients; i++ {
			clients[i] = createTestClient(hub, nil, "User")
			hub.register <- clients[i]
		}
		time.Sleep(10 * time.Millisecond)

		var wg sync.WaitGroup

		// Concurrently join and leave
		for i := 0; i < numClients; i++ {
			wg.Add(1)
			go func(client *Client) {
				defer wg.Done()
				hub.joinRoom <- &RoomAction{Client: client, RaceID: raceID}
				time.Sleep(5 * time.Millisecond)
				hub.leaveRoom <- &RoomAction{Client: client, RaceID: raceID}
			}(clients[i])
		}

		wg.Wait()
		time.Sleep(20 * time.Millisecond)

		// Room should be empty (all left)
		count := hub.GetRoomClientCount(raceID)
		assert.Equal(t, 0, count, "Room should be empty after all clients leave")
	})

	t.Run("Broadcast to multiple rooms concurrently", func(t *testing.T) {
		numRooms := 5
		clientsPerRoom := 3
		rooms := make([]string, numRooms)
		allClients := make([]*Client, numRooms*clientsPerRoom)

		// Create rooms and clients
		for i := 0; i < numRooms; i++ {
			rooms[i] = "room-" + string(rune(i))
			for j := 0; j < clientsPerRoom; j++ {
				idx := i*clientsPerRoom + j
				allClients[idx] = createTestClient(hub, nil, "User")
				hub.register <- allClients[idx]
			}
		}
		time.Sleep(10 * time.Millisecond)

		// Join clients to their rooms
		for i := 0; i < numRooms; i++ {
			for j := 0; j < clientsPerRoom; j++ {
				idx := i*clientsPerRoom + j
				hub.joinRoom <- &RoomAction{Client: allClients[idx], RaceID: rooms[i]}
			}
		}
		time.Sleep(10 * time.Millisecond)

		var wg sync.WaitGroup

		// Concurrently broadcast to all rooms
		for i := 0; i < numRooms; i++ {
			wg.Add(1)
			go func(roomID string) {
				defer wg.Done()
				hub.BroadcastToRoom(roomID, []byte("room-specific message"))
			}(rooms[i])
		}

		wg.Wait()
		time.Sleep(10 * time.Millisecond)

		// Verify each room has correct client count
		for i := 0; i < numRooms; i++ {
			count := hub.GetRoomClientCount(rooms[i])
			assert.Equal(t, clientsPerRoom, count, "Room %s should have %d clients", rooms[i], clientsPerRoom)
		}
	})
}

