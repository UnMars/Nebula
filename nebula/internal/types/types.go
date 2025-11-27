package types

import "nebula/internal/message"

// Client interface
type Clienter interface {
	// Send channel to send messages to the client
	SendChannel() chan<- []byte
	// Username of the client
	Username() string
	// Current room of the client
	CurrentRoom() string
	// Close the send channel
	CloseSendChannel()
}

// Hub interface
type Hubber interface {
	// Register a client to the hub
	RegisterClient(Clienter)
	// Unregister a client from the hub
	UnregisterClient(Clienter)
	// Broadcast a message to clients of the room of the message
	BroadcastMessage(message.BroadcastMessage)
}
