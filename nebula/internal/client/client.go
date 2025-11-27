package client

import (
	"bytes"
	"log"
	"nebula/internal/message"
	"nebula/internal/types"
	"time"

	"github.com/gorilla/websocket"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client struct, implements Clienter interface
type Client struct {
	// Hub of the client
	hub types.Hubber
	// Websocket connection
	conn *websocket.Conn
	// Send channel
	send chan []byte
	// Username of the client
	username string
	// Current room name of the client
	room string
}

// Close the send channel
func (c *Client) CloseSendChannel() {
	close(c.send)
}

// Get the send channel
func (c *Client) SendChannel() chan<- []byte {
	return c.send
}

// Get the current room of the client
func (c *Client) CurrentRoom() string {
	return c.room
}

// Get the username of the client
func (c *Client) Username() string {
	return c.username
}

func NewClient(hub types.Hubber, conn *websocket.Conn, username string, room string) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
		room:     room,
	}
}

// The Read pump function handles incoming messages from the client
// It is triggered when a client send a message
func (c *Client) ReadPump() {
	defer func() {
		c.hub.UnregisterClient(c)
		c.conn.Close()
	}()

	// Set connection options
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Time{})
	c.conn.SetPongHandler(func(string) error { return c.conn.SetReadDeadline(time.Time{}) })

	for {
		// Read message from the client
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// Format message
		msg = bytes.TrimSpace(bytes.ReplaceAll(msg, newline, space))
		if len(msg) == 0 {
			continue
		}
		// Broadcast message to the room
		c.hub.BroadcastMessage(message.BroadcastMessage{
			Room: c.room,
			Data: append([]byte(c.username+": "), msg...),
		})

	}
}

// The Write pump function handles outgoing messages to the client
// It is triggered when a message is sent to the client (via a room)
// It gets messages from the send channel and send them to the client
func (c *Client) WritePump() {
	// Heartbeat each 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.send:
			// If send channel is closed, close the connection
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// Set write deadline of 10s to avoid blocking in case of slow client
			err := c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				return
			}

			// Get next writer
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			// Write message
			w.Write(msg)

			// Write pending messages
			for range len(c.send) {
				w.Write(<-c.send)
			}

			// Close writer
			err = w.Close()
			if err != nil {
				return
			}

		case <-ticker.C:
			// Send new ping message
			err := c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				return
			}
			err = c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}
		}
	}
}
