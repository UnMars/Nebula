package main

import (
	"log"
	"nebula/internal/hub"
	"nebula/internal/server"
	"net/http"
)

func main() {
	// Run the hub
	hub := hub.NewHub()
	go hub.Run()

	// Setup routes
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWS(hub, w, r)
	})
	http.Handle("/", http.FileServer(http.Dir("web/static")))

	// Run the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
