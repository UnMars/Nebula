# GO EXPERT PROJECT — NEBULA v3
**Distributed, high-performance, and resilient real-time chat**
Complete subject to follow step-by-step — Level ENSIMAG / École 42 / Go Backend Internship

Total duration: 4 to 6 weeks
Final objective: a cluster of Go servers capable of handling 20,000+ simultaneous WebSocket connections, with persistence, node restart tolerance, monitoring, and zero goroutine leaks.

This subject is designed to be followed sequentially. You validate each step before moving on to the next. By the end, you will have a publishable project on GitHub that will impress any Go recruiter.

## Final project objective (what you will have on the last day)

- `docker-compose up --scale nebula=3` → 3 instances + Redis + Prometheus/Grafana
- 20,000+ simultaneous WebSocket connections distributed across the 3 nodes
- Users in the same room see each other even if they are connected to different nodes
- Node restart → zero lost messages, history restored
- Real-time metrics (connections, messages/s, latency, goroutines, RAM)
- `go test -race ./...` → 100% clean
- 5-page PDF report + flamegraphs + k6 results

## Final features (realistic and ordered scope)

| Feature                                        | Step  | Mandatory   |
|------------------------------------------------|-------|-------------|
| WebSocket connection + rooms                   | 1     | Yes         |
| Broadcast in a room (local)                    | 1     | Yes         |
| Presence (list of connected users)             | 1     | Yes         |
| Username + signed cookie (no heavy JWT)        | 1     | Yes         |
| Rate limiting by IP/user                       | 2     | Yes         |
| Room sharding (32–64 shards)                   | 2     | Yes         |
| Asynchronous persistence (AOF or SQLite)       | 3     | Yes         |
| History replay on startup                      | 3     | Yes         |
| Cluster via Redis Pub/Sub                      | 4     | Yes         |
| Prometheus metrics + /metrics endpoint         | 5     | Yes         |
| Complete graceful shutdown (flush + clean disconnect) | 5 | Yes         |
| Automatic reconnection with same userID        | 6     | Bonus       |
| Protobuf instead of JSON                       | 6     | Bonus       |
| Simple HTMX dashboard                          | 6     | Bonus       |

## Non-negotiable technical constraints

1. Zero web framework (just `net/http` + `nhooyr/websocket` or `gorilla/websocket`)
2. Zero goroutine leak → `go test -race ./...` must be flawless
3. Complete graceful shutdown in a cluster
4. < 100 allocs/op on the critical path (broadcast)
5. Everything must run with `docker-compose up`

## Final project structure

```bash
nebula/
├── cmd/
│   └── nebula/
│       └── main.go                  # server + config + graceful shutdown
├── internal/
│   ├── server/                      # HTTP handlers + WS upgrade
│   │   ├── handler.go
│   │   └── middleware.go
│   ├── hub/
│   │   ├── shard.go                 # one shard = one goroutine + map[room]*Room
│   │   ├── sharding.go              # hash function roomName → shardID
│   │   └── redis_bridge.go          # Redis pub/sub for the cluster
│   ├── room/
│   │   └── room.go                  # local clients + broadcast channel
│   ├── client/
│   │   └── client.go
│   ├── persistence/
│   │   ├── writer.go                # buffered async writer (AOF or SQLite)
│   │   └── replayer.go
│   ├── auth/
│   │   └── userid.go                # signed cookie → userID
│   ├── metrics/
│   │   └── prometheus.go
│   └── protocol/
│       └── message.proto            # (Protobuf bonus)
├── pkg/
│   └── rate/
│       └── limiter.go               # token bucket
├── web/
│   └── static/
│       ├── index.html
│       └── app.js                   # ultra-simple test client
├── scripts/
│   ├── loadtest.k6.js
│   └── cluster.sh
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── Makefile
└── README.md
```

## The 6 detailed steps (to be validated one by one)

### STEP 1 — Functional Local MVP (4–6 days)

Objective: a chat that works perfectly on a single instance.

Deliverables:
- WebSocket connection
- Create/join a room
- Send/receive messages
- Real-time presence
- Username via signed cookie
- Working HTML/JS page

Simple architecture: a single global `Hub` with `sync.RWMutex` (it's OK at this stage).

### STEP 2 — Local Performance & Sharding (5–7 days)

Objective: go from 500 → 8,000+ connections on one machine.

To do:
- Remove global mutex → 32 or 64 shards (each shard has its own goroutine and map)
- Use `sync.Pool` for message buffers
- Token bucket rate limiter (pkg/rate/limiter.go)
- Measure with pprof → aim for < 100 allocs/op on broadcast
- Local load test with k6 → 8,000 simultaneous connections

### STEP 3 — Asynchronous Persistence & Replay (5–7 days)

Objective: no longer lose history on restart.

Two options (choose the one you want):
- Simple option: Append-Only File (AOF) → one JSON line per message
- Cleaner option: SQLite with WAL mode

Implementation:
- A writer goroutine with buffered channel (10,000 messages)
- On every message → send to writer + broadcast
- On startup → replay file/database to reconstruct state

### STEP 4 — Cluster Mode with Redis (6–8 days)

Objective: two different instances can exchange messages.

Architecture:
- Each node manages its local clients
- When a client sends a message → local broadcast + Redis publication `channel:roomName`
- Each node is subscribed to all existing rooms (or via pattern)
- Message received from Redis → local broadcast only

Realistic bonus: automatic room discovery via Redis

### STEP 5 — Observability & Cluster Graceful Shutdown (4–6 days)

To add:
- Prometheus `/metrics` endpoint (number of connections, messages/s, goroutines, etc.)
- `/healthz` healthcheck
- Graceful shutdown:
  - Stop accepting new connections
  - Wait for all writers to flush
  - Close all WebSocket connections cleanly
  - Close Redis
- docker-compose with preconfigured Prometheus + Grafana

### STEP 6 — Bonus, Tests, Cleanup & Final Report (5–10 days)

- Unit + integration tests (especially hub + persistence)
- Broadcaster benchmarks
- GitHub Actions CI with race detector
- 5-page PDF report (architecture, difficulties, flamegraphs, k6 results)
- Epic README with GIFs

## Final deliverables required (to have on the last day)

1. Clean public GitHub repo
2. `docker-compose up --scale nebula=3` that works in < 15 seconds
3. README with:
   - GIF of chat in cluster
   - Test commands
   - CI badges
4. 4–6 page PDF report (ENSIMAG submission style)
5. k6 results (20k+ connections, p95 latency < 150ms)

## Skills you will have mastered 100%

- Advanced Go concurrency (sharding, channels, context)
- Memory management (sync.Pool, zero-copy when possible)
- Asynchronous persistence
- Simple distributed systems (Redis pub/sub)
- Observability (Prometheus, pprof)
- Graceful shutdown in a cluster
- Real load tests
- Docker + docker-compose