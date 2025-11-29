package hub

import (
	"encoding/json"
	"nebula/internal/message"
	"testing"
	"time"
)

// MockClient implements types.Clienter for testing
type MockClient struct {
	username string
	room     string
	send     chan message.BroadcastMessage
}

func NewMockClient(username, room string) *MockClient {
	return &MockClient{
		username: username,
		room:     room,
		send:     make(chan message.BroadcastMessage, 10), // Buffer to avoid blocking
	}
}

func (m *MockClient) SendChannel() chan<- message.BroadcastMessage {
	return m.send
}

func (m *MockClient) Username() string {
	return m.username
}

func (m *MockClient) CurrentRoom() string {
	return m.room
}

func (m *MockClient) CloseSendChannel() {
	close(m.send)
}

func expectMessage(t *testing.T, client *MockClient) message.BroadcastMessage {
	t.Helper()
	select {
	case msg := <-client.send:
		return msg
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for message")
		return message.BroadcastMessage{}
	}
}

// -------------------------------------------------------------------------
// Step 1 Validation Tests
// -------------------------------------------------------------------------

// Criteria 1: Protocol Compliance
// The server must send/receive valid JSON.
func TestStep1_ProtocolCompliance(t *testing.T) {
	// Test Marshaling (Server -> Client)
	msg := message.BroadcastMessage{
		Type:    "message",
		Room:    "general",
		Sender:  "Alice",
		Content: "Hello",
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	// Verify JSON structure
	var decoded map[string]interface{}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if decoded["type"] != "message" {
		t.Errorf("Expected type 'message', got %v", decoded["type"])
	}
	if decoded["content"] != "Hello" {
		t.Errorf("Expected content 'Hello', got %v", decoded["content"])
	}

	// Test omitempty for Users
	msgPresence := message.BroadcastMessage{
		Type:  "presence",
		Room:  "general",
		Users: []string{"Alice", "Bob"},
	}
	bytes, err = json.Marshal(msgPresence)
	if err != nil {
		t.Fatalf("Failed to marshal presence: %v", err)
	}

	// Check that "content" is NOT present (omitempty)
	var decodedPresence map[string]interface{}
	if err := json.Unmarshal(bytes, &decodedPresence); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Content should be present as empty string because it does NOT have omitempty in struct definition
	// Wait, let's check types.go again.
	// Content string `json:"content"` -> No omitempty.
	// So it SHOULD be present.
	if _, ok := decodedPresence["content"]; !ok {
		t.Errorf("Expected content field to be present (even if empty), but it was missing")
	}
	if decodedPresence["content"] != "" {
		t.Errorf("Expected empty content, got '%v'", decodedPresence["content"])
	}
}

// Criteria 2: Room Isolation
// Users in "dev" cannot see messages from "general".
func TestStep1_RoomIsolation(t *testing.T) {
	h := NewHub()
	go h.Run()
	defer h.Close()

	alice := NewMockClient("Alice", "general")
	bob := NewMockClient("Bob", "general")
	charlie := NewMockClient("Charlie", "dev")

	h.RegisterClient(alice)
	expectMessage(t, alice) // join
	expectMessage(t, alice) // presence

	h.RegisterClient(bob)
	expectMessage(t, bob)   // join
	expectMessage(t, bob)   // presence
	expectMessage(t, alice) // join (Bob)
	expectMessage(t, alice) // presence

	h.RegisterClient(charlie)
	expectMessage(t, charlie) // join
	expectMessage(t, charlie) // presence

	// Alice sends a message to general
	msg := message.BroadcastMessage{
		Room:    "general",
		Type:    "message",
		Sender:  "Alice",
		Content: "Hello General",
	}
	h.BroadcastMessage(msg)

	// Bob should receive it
	m := expectMessage(t, bob)
	if m.Content != "Hello General" {
		t.Errorf("Bob received wrong message: %s", m.Content)
	}

	// Charlie should NOT receive it
	select {
	case m := <-charlie.send:
		t.Errorf("Charlie received message from another room: %v", m)
	case <-time.After(100 * time.Millisecond):
		// OK
	}
}

// Criteria 3 & 4: Presence System & Real-time Updates
// When I join, I see who is already there.
// When I leave, others see my name disappear immediately.
func TestStep1_PresenceAndRealTime(t *testing.T) {
	h := NewHub()
	go h.Run()
	defer h.Close()

	alice := NewMockClient("Alice", "general")
	bob := NewMockClient("Bob", "general")

	// 1. Alice Joins
	h.RegisterClient(alice)

	// Alice receives Join(Alice)
	msg := expectMessage(t, alice)
	if msg.Type != "join" || msg.Sender != "Alice" {
		t.Errorf("Expected join(Alice), got %v", msg)
	}

	// Alice receives Presence([Alice])
	msg = expectMessage(t, alice)
	if msg.Type != "presence" {
		t.Errorf("Expected presence, got %v", msg)
	}
	if len(msg.Users) != 1 || msg.Users[0] != "Alice" {
		t.Errorf("Expected users [Alice], got %v", msg.Users)
	}

	// 2. Bob Joins
	h.RegisterClient(bob)

	// Bob receives Join(Bob)
	expectMessage(t, bob)
	// Bob receives Presence([Alice, Bob]) (Order depends on map iteration, so check containment)
	msg = expectMessage(t, bob)
	if len(msg.Users) != 2 {
		t.Errorf("Expected 2 users, got %v", msg.Users)
	}

	// Alice receives Join(Bob)
	msg = expectMessage(t, alice)
	if msg.Type != "join" || msg.Sender != "Bob" {
		t.Errorf("Expected join(Bob), got %v", msg)
	}
	// Alice receives Presence([Alice, Bob])
	msg = expectMessage(t, alice)
	if len(msg.Users) != 2 {
		t.Errorf("Expected 2 users, got %v", msg.Users)
	}

	// 3. Alice Leaves
	h.UnregisterClient(alice)

	// Bob receives Leave(Alice)
	msg = expectMessage(t, bob)
	if msg.Type != "leave" || msg.Sender != "Alice" {
		t.Errorf("Expected leave(Alice), got %v", msg)
	}

	// Bob receives Presence([Bob])
	msg = expectMessage(t, bob)
	if msg.Type != "presence" {
		t.Errorf("Expected presence, got %v", msg)
	}
	if len(msg.Users) != 1 || msg.Users[0] != "Bob" {
		t.Errorf("Expected users [Bob], got %v", msg.Users)
	}
}

// Criteria 6: Graceful Shutdown
// Ctrl+C stops the server without errors or zombie goroutines.
func TestStep1_GracefulShutdown(t *testing.T) {
	h := NewHub()
	go h.Run()

	c1 := NewMockClient("Alice", "general")
	h.RegisterClient(c1)

	// Drain initial messages
	expectMessage(t, c1)
	expectMessage(t, c1)

	// Close hub
	h.Close()

	// Verify client channel is closed
	select {
	case _, ok := <-c1.send:
		if ok {
			t.Error("Client channel should be closed")
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for channel close")
	}
}
