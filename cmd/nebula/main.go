package main

import (
	"context"
	"log"
	"nebula/internal/hub"
	"nebula/internal/server"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Start the hub
	hub := hub.NewHub()
	go hub.Run()

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWS(hub, w, r)
	})
	mux.Handle("/", http.FileServer(http.Dir("web/static")))

	webServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Run the HTTP server in a goroutine
	go func() { webServer.ListenAndServe() }()
	log.Println("Server started on :8080")

	// Profiling server
	go func() {
		log.Println("Starting pprof server on localhost:6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Wait for termination signal
	<-signalChan

	log.Println("Shutting down server...")
	hub.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	webServer.Shutdown(ctx)

}
