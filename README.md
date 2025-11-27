# Nebula

**Nebula** is a high-performance, distributed real-time chat server written in Go.
This project is designed to scale to 20,000+ simultaneous WebSocket connections across multiple nodes.

## 🚀 Features (Step 1 MVP)

- **WebSocket Support**: Native WebSocket handling without heavy frameworks.
- **Real-time Chat**: Instant messaging with low latency.
- **Rooms**: Support for multiple chat rooms (`?room=name`).
- **Concurrency**: Thread-safe Hub implementation using `sync.RWMutex`.
- **Clean Architecture**: Modular design separating `hub`, `client`, `room`, and `server`.

## 🛠️ Installation & Run

### Prerequisites
- Go 1.25+

### Running the Server

```bash
# Clone the repository
git clone https://github.com/UnMars/Nebula.git
cd Nebula

# Install dependencies
go mod tidy

# Run the server
go run cmd/nebula/main.go
```

The server will start on `http://localhost:8080`.

### Connecting

Open your browser and navigate to:
`http://localhost:8080` (defaults to `general` room and random username)

Or specify room and username:
`http://localhost:8080/?room=devops&username=alice`

## 📂 Project Structure

```
nebula/
├── cmd/
│   └── nebula/      # Main entry point
├── internal/
│   ├── client/      # Client WebSocket handling (Read/Write pumps)
│   ├── hub/         # Central Hub for managing rooms and broadcast
│   ├── room/        # Room structure
│   ├── server/      # HTTP handlers
│   └── types/       # Shared interfaces
└── web/
    └── static/      # Simple frontend for testing
```

## 📚 Learning Goals

This project follows a strict "Cahier des Charges" to master:
- Go Concurrency (Goroutines, Channels, Mutexes)
- WebSockets internals
- Distributed Systems (Sharding, Redis Pub/Sub - *Coming in later steps*)
- Optimization (Memory management, Zero-copy)
