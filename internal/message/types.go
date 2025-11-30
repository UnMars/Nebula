package message

type BroadcastMessage struct {
	// Message type ("message", "join", "leave", "presence")
	Type string `json:"type"`
	// Room name
	Room string `json:"room"`
	// Sender username
	Sender string `json:"sender"`
	// Message content
	Content string `json:"content"`
	// Users (only for presence messages)
	Users []string `json:"users,omitempty"`
	// Timestamp
	SendAt int64 `json:"sendAt"`
}
