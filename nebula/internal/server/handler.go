package server

import (
	"log"
	"nebula/internal/client"
	"nebula/internal/hub"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeWS is the websocket handler of the server
// On new connection, it upgrades the HTTP connection
// to a websocket connection and registers the client to the hub
func ServeWS(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	// TODO: better way to handle and log errors
	if err != nil {
		log.Println(err)
		return
	}

	roomName := r.URL.Query().Get("room")
	if roomName == "" {
		roomName = "general"
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		username = "anonymous"
	}

	// Create new client and register it to the hub
	cl := client.NewClient(h, conn, username, roomName)
	h.RegisterClient(cl)

	// Start client's read and write pumps
	go cl.WritePump()
	go cl.ReadPump()
}
