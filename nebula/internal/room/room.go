package room

import "nebula/internal/types"

type Room struct {
	// Name of the room
	Name string
	// Clients in the room (Set-like)
	Clients map[types.Clienter]bool
}

func NewRoom(name string) *Room {
	return &Room{
		Name:    name,
		Clients: make(map[types.Clienter]bool),
	}
}
