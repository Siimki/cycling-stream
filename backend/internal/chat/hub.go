package chat

import (
	"sync"
)

// RoomAction represents an action to join or leave a room
type RoomAction struct {
	Client *Client
	RaceID string
}

// Hub maintains the set of active clients and broadcasts messages to rooms
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Race-specific rooms (race_id -> clients in that room)
	rooms map[string]map[*Client]bool

	// Inbound messages from clients
	broadcast chan []byte

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Join room requests
	joinRoom chan *RoomAction

	// Leave room requests
	leaveRoom chan *RoomAction

	// Mutex for thread-safe access
	mu sync.RWMutex
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		joinRoom:   make(chan *RoomAction),
		leaveRoom:  make(chan *RoomAction),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.RegisterClient(client)

		case client := <-h.unregister:
			h.UnregisterClient(client)

		case action := <-h.joinRoom:
			h.JoinRoom(action.Client, action.RaceID)

		case action := <-h.leaveRoom:
			h.LeaveRoom(action.Client, action.RaceID)

		case message := <-h.broadcast:
			// Broadcast to all clients (not used for room-specific messages)
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// RegisterClient adds a client to the hub
func (h *Hub) RegisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = true
}

// UnregisterClient removes a client from the hub and all rooms
func (h *Hub) UnregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)

		// Remove client from all rooms
		for raceID, roomClients := range h.rooms {
			if _, inRoom := roomClients[client]; inRoom {
				delete(roomClients, client)
				// Clean up empty rooms
				if len(roomClients) == 0 {
					delete(h.rooms, raceID)
				}
			}
		}
	}
}

// JoinRoom adds a client to a race-specific room
func (h *Hub) JoinRoom(client *Client, raceID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[raceID] == nil {
		h.rooms[raceID] = make(map[*Client]bool)
	}
	h.rooms[raceID][client] = true
}

// LeaveRoom removes a client from a race-specific room
func (h *Hub) LeaveRoom(client *Client, raceID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if roomClients, ok := h.rooms[raceID]; ok {
		delete(roomClients, client)
		// Clean up empty rooms
		if len(roomClients) == 0 {
			delete(h.rooms, raceID)
		}
	}
}

// BroadcastToRoom sends a message to all clients in a specific room
func (h *Hub) BroadcastToRoom(raceID string, message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if roomClients, ok := h.rooms[raceID]; ok {
		for client := range roomClients {
			select {
			case client.send <- message:
			default:
				// Client's send channel is full, close and remove
				close(client.send)
				delete(h.clients, client)
				delete(roomClients, client)
			}
		}
	}
}

// GetRoomClientCount returns the number of clients in a room
func (h *Hub) GetRoomClientCount(raceID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if roomClients, ok := h.rooms[raceID]; ok {
		return len(roomClients)
	}
	return 0
}

