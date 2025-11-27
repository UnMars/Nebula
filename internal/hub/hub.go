package hub

import (
	"nebula/internal/message"
	"nebula/internal/room"
	"nebula/internal/types"
	"sync"
)

// Hub of the server, it handles clients registration,
// unregistration and broadcast
// Implements Hubber interface
type Hub struct {
	// Mutex for concurrent access to the rooms map
	mu sync.RWMutex
	// All rooms of the hub (Set-like)
	rooms map[string]*room.Room
	// Channel to broadcast messages
	broadcast chan message.BroadcastMessage
}

// Register a client to the hub
func (h *Hub) RegisterClient(client types.Clienter) {
	roomName := client.CurrentRoom()
	h.mu.Lock()

	// Create room if it doesn't exist
	if _, ok := h.rooms[roomName]; !ok {
		h.rooms[roomName] = room.NewRoom(roomName)
	}
	r := h.rooms[roomName]

	// Add client to the room
	r.Clients[client] = true
	h.mu.Unlock()

	// Send welcome message to the room
	welcomeMessage := message.BroadcastMessage{
		Room: roomName,
		Data: []byte(client.Username() + " joined the room."),
	}
	h.BroadcastMessage(welcomeMessage)
}

// Unregister a client from the hub
func (h *Hub) UnregisterClient(client types.Clienter) {
	roomName := client.CurrentRoom()

	h.mu.Lock()
	// Remove client from the room
	if r, ok := h.rooms[roomName]; ok {
		delete(r.Clients, client)
		// Remove room if it's empty
		if len(r.Clients) == 0 {
			delete(h.rooms, roomName)
		}
	}
	h.mu.Unlock()

	// Send leave message to the room
	leaveMessage := message.BroadcastMessage{
		Room: roomName,
		Data: []byte(client.Username() + " left the room."),
	}
	h.BroadcastMessage(leaveMessage)

	// Close client's send channel
	client.CloseSendChannel()
}

// Broadcast a message to clients of the room of the message
func (h *Hub) BroadcastMessage(msg message.BroadcastMessage) {
	h.broadcast <- msg
}

// Create a new hub
func NewHub() *Hub {
	return &Hub{
		rooms:     make(map[string]*room.Room),
		broadcast: make(chan message.BroadcastMessage),
	}
}

// Hub main function
func (h *Hub) Run() {
	for msg := range h.broadcast {
		// Read lock (for rooms map)
		h.mu.RLock()
		var deadClients []types.Clienter

		if room, ok := h.rooms[msg.Room]; ok {
			for client := range room.Clients {
				select {
				case client.SendChannel() <- msg.Data:
				// Dead / unresponsive client, remove it
				default:
					deadClients = append(deadClients, client)
				}
			}
		}
		h.mu.RUnlock()

		// Unregister dead clients
		for _, client := range deadClients {
			h.UnregisterClient(client)
		}
	}
}
