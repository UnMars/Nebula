# 🌌 Project Nebula: High-Performance Distributed WebSocket Engine
**Advanced Go Systems Programming — Final Capstone Project**

| Metadata | Details |
| :--- | :--- |
| **Level** | Advanced (Master 2 / End of Studies) |
| **Focus** | Concurrency, Distributed Systems, High Availability |
| **Stack** | Go 1.25+, Redis, Docker, Prometheus |
| **Duration** | 4 to 6 weeks |
| **Team** | Solo or Pair |

---

## 1. 📝 Introduction

### The Pitch
Modern real-time applications (Discord, Slack, WhatsApp) handle millions of concurrent connections. They don't achieve this by accident; they rely on highly optimized, distributed architectures.

Your mission is to build **Nebula**: a distributed, fault-tolerant, and high-performance WebSocket chat engine.

Unlike a standard "bootcamp chat app", Nebula is designed to be **production-grade**. It must handle **20,000+ simultaneous connections**, survive node failures, and provide real-time observability. You will not use high-level web frameworks (like Gin or Fiber). You will build the core engine using the Go standard library to understand exactly how memory and concurrency work under the hood.

### Pedagogical Objectives
By the end of this project, you will have mastered:
*   **Advanced Go Concurrency**: Channels, Mutexes, WaitGroups, and avoiding race conditions.
*   **Memory Optimization**: Profiling (pprof), minimizing allocations, and using `sync.Pool`.
*   **Distributed Systems**: Pub/Sub patterns, eventual consistency, and sharding.
*   **Resilience**: Graceful shutdowns, crash recovery, and signal handling.
*   **DevOps**: Docker composition, load testing (k6), and monitoring (Prometheus/Grafana).

---

## 2. ⚖️ General Rules & Constraints

### Technical Constraints
1.  **Language**: Go (latest stable version).
2.  **Web Frameworks**: **FORBIDDEN**. You must use `net/http`.
3.  **WebSocket Library**: `github.com/gorilla/websocket` or `github.com/nhooyr/websocket` are allowed.
4.  **Database**: Redis (for Pub/Sub) and a file-based persistence (AOF or SQLite).
5.  **Code Quality**:
    *   `go vet` and `staticcheck` must pass.
    *   **Zero Data Races**: `go test -race ./...` must be 100% clean.
    *   Code must be formatted with `gofmt`.

### Project Structure
You must adhere to the standard Go project layout:
```
nebula/
├── cmd/nebula/       # Main entry point
├── internal/         # Private application and library code
│   ├── server/       # HTTP & WebSocket transport layer
│   ├── hub/          # Core logic (Room management)
│   ├── protocol/     # Serialization & Message definitions
│   └── ...
├── pkg/              # Library code safe to use by external apps
├── web/              # Frontend assets (for testing)
└── deployments/      # Docker & CI/CD configurations
```

---

## 3. 🗺️ The Roadmap

The project is divided into **6 distinct milestones**. You must validate each step before moving to the next.

### 🟢 Step 1: The Functional MVP
**Objective**: A working chat server on a single node.
*   **Features**:
    *   WebSocket handshake & upgrade.
    *   Concept of **Rooms**: Users can join/leave specific channels.
    *   **Broadcast**: Messages sent by one user are received by all others in the room.
    *   **Presence**: System messages when users join/leave ("Alice joined the room").
    *   **Protocol**: Define a strict JSON protocol (e.g., `{"type": "msg", "room": "general", "payload": "..."}`).
*   **Deliverable**: A server that works with the provided HTML/JS client.

### 🔵 Step 2: High Performance & Sharding
**Objective**: Scale from 500 to 10,000+ concurrent connections on a single machine.
*   **The Problem**: A single `sync.Mutex` on the central Hub will become a bottleneck.
*   **The Solution**: **Sharding**.
    *   Split the Hub into 32 or 64 "Shards".
    *   Route rooms to shards using a hash function (e.g., `CRC32(roomName) % numShards`).
    *   Each shard manages its own locks and clients.
*   **Optimization**:
    *   Implement `sync.Pool` for reusing message buffers (reduce GC pressure).
    *   Implement a Rate Limiter (Token Bucket) to prevent spam.
*   **Validation**: Run a local benchmark. The server must handle 5k connections with low latency.

### 🟠 Step 3: Persistence & Reliability
**Objective**: Zero data loss on restart.
*   **Feature**: When the server restarts, previous messages in a room should be replayed to new users.
*   **Implementation**:
    *   Create an **Asynchronous Writer** (don't block the broadcast loop!).
    *   Storage Strategy: **Append-Only File (AOF)** or **SQLite (WAL mode)**.
    *   On startup: Read the storage and repopulate the in-memory history of active rooms.
*   **Constraint**: The write operation must not degrade broadcast performance.

### 🔴 Step 4: The Distributed Cluster
**Objective**: Horizontal scaling. Run multiple Nebula nodes that talk to each other.
*   **Scenario**: User A connects to Node 1. User B connects to Node 2. They are in the same room. They must be able to chat.
*   **Architecture**:
    *   Use **Redis Pub/Sub** as the inter-node bus.
    *   When Node 1 receives a message for "Room X", it publishes it to Redis channel `nebula:room:X`.
    *   Node 2 (subscribed to `nebula:room:X`) receives the payload and broadcasts it to its local clients.
*   **Deliverable**: `docker-compose up --scale nebula=3`.

### 🟣 Step 5: Production Readiness
**Objective**: Observability and Graceful Shutdown.
*   **Metrics**: Expose a `/metrics` endpoint for **Prometheus**.
    *   Gauges: `connected_clients`, `active_rooms`, `goroutines`.
    *   Counters: `messages_total`, `errors_total`.
    *   Histograms: `broadcast_duration_seconds`.
*   **Graceful Shutdown**:
    *   Catch `SIGINT`/`SIGTERM`.
    *   Stop accepting new connections.
    *   Flush pending writes to disk.
    *   Close existing WebSockets with a "Server Shutting Down" frame.
    *   Exit only when safe.

### 🏆 Step 6: Bonus & Excellence
**Objective**: Go above and beyond.
*   **Binary Protocol**: Replace JSON with **Protobuf** for tighter payloads.
*   **Dynamic Resharding**: Handle adding/removing nodes from the cluster dynamically.
*   **Dashboard**: A real-time admin dashboard (using HTMX or React) showing cluster health.
*   **CI/CD**: GitHub Actions pipeline running tests and linters.

---

## 4. 🧪 Evaluation Criteria

Your project will be evaluated on:
1.  **Stability**: Does it crash under load? (It shouldn't).
2.  **Performance**: Is the latency acceptable with 10k users?
3.  **Code Quality**: Is the code idiomatic Go? Are mutexes used correctly?
4.  **Architecture**: Is the separation of concerns respected?

## 5. 📚 Resources

*   [The Go Memory Model](https://go.dev/ref/mem)
*   [Gorilla WebSocket Documentation](https://pkg.go.dev/github.com/gorilla/websocket)
*   [Redis Pub/Sub](https://redis.io/docs/manual/pubsub/)
*   [Prometheus Go Client](https://github.com/prometheus/client_golang)
*   [1M connections in Go (WebSocket)](https://www.freecodecamp.org/news/million-websockets-and-go-cc58418460bb/)

---
*Good luck. The system is waiting.*