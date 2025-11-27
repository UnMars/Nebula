package message

type BroadcastMessage struct {
	// Room name
	Room string
	// Message data
	Data []byte
}
