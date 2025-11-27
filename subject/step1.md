# STEP 1 — FUNCTIONAL LOCAL MVP

**Session Objective**: Have a complete WebSocket chat that works perfectly on a single instance (one machine, one process).
Target duration: 4 to 6 hours (one evening or a short weekend)

By the end of this step you must be able to:

- Run `go run ./cmd/nebula`
- Open two (or ten) browser tabs
- Create / join a room
- See messages in real-time in both directions
- See presence ("alice joined", "bob left")
- See the list of users connected in the room updated instantly

### 1. Step 1 Constraints

| Constraint                                          | Mandatory |
| --------------------------------------------------- | ----------- |
| Zero web framework (Gin, Echo, Fiber, …)            | Yes         |
| Use only `net/http` + a WebSocket lib               | Yes         |
| `go test -race ./...` must pass without data race   | Yes         |
| Basic graceful shutdown (Ctrl+C closes cleanly)     | Yes         |
| No lost messages as long as the server is running   | Yes         |
| Readable code well structured in packages           | Yes         |

### 2. Exact structure to create from the start

```bash
nebula/
├── cmd/
│   └── nebula/
│       └── main.go                  # server launch + graceful shutdown
├── internal/
│   ├── server/
│   │   ├── handler.go               # /ws endpoint + upgrade
│   │   └── middleware.go            # (optional for now)
│   ├── hub/
│   │   └── hub.go                   # Central Hub for step 1
│   ├── client/
│   │   └── client.go
│   ├── room/
│   │   └── room.go
│   └── message/
│       └── types.go                 # structures of exchanged messages
├── web/
│   └── static/
│       ├── index.html
│       └── app.js                   # minimal test client
├── go.mod
├── Makefile
└── README.md
```

### 3. Technical choices to make immediately

1. WebSocket Library

   - [ ] `github.com/gorilla/websocket` (recommended for speed)
   - [ ] `github.com/nhooyr/websocket` (more modern, to consider later)

2. Message Protocol (step 1)
   - Simple JSON (accept the cost, we will optimize later)
   - Message types to support:
     - `join`
     - `leave`
     - `message`
     - `presence` (list of users)

### 4. Design questions to solve yourself

1. How to represent a `Client`?
   ID uuid
   Name string
   createdAt timestamp

2. How to represent a `Room`?
   clients map[clients]
   createdAt timestamp
   messages []string

3. How to represent the central `Hub`?
   List of all clients
   register channel
   unregister channel
   broadcast channel

4. Which channels are needed in the Hub? (register / unregister / broadcast ?)
   register channel for a Client to join a Room
   unregister channel for a Client to leave a Room
   broadcast channel to send a message to all current Clients

5. Should there be a dedicated goroutine per Room or one for the whole Hub?
   One for the whole Hub, it's the Hub that manages the rooms, rooms are just an ID essentially.

6. How to cleanly handle client disconnection?
   Close their channels, make them leave all rooms, unregister them from the hub.

7. How to propagate presence events?
   Hub broadcasts to the list of Clients in a Room?

8. How to simply sign a username (signed cookie or UUID)?
   UUID

### 5. Expected WebSocket connection flow

1. HTTP Request → `/ws`
2. Upgrade to WebSocket
3. Creation of a `*Client`
4. Registration of the client in the Hub
5. Launch of the two pumps (`readPump` and `writePump`)
6. Reading messages → processing → broadcast
7. Writing messages received from the client's channel

### 6. Validation Checklist (to check before moving to step 2)

- [ ] `go run ./cmd/nebula` starts on :8080
- [ ] The `index.html` page loads and connects automatically
- [ ] Two tabs in the same room see messages instantly
- [ ] When a tab is closed → others see "X left the room"
- [ ] The list of present users is updated in real-time
- [ ] Ctrl+C stops the server without zombie goroutines
- [ ] `go test -race ./...` passes without warning
- [ ] No panic observed after 10 minutes of intensive testing

### 7. Authorized Resources (to keep open)

- https://pkg.go.dev/github.com/gorilla/websocket
- https://github.com/gorilla/websocket/tree/master/examples/chat (structure only)
- https://gowebexamples.com/websockets/
- https://go.dev/blog/pipelines
- https://pkg.go.dev/sync#RWMutex
- https://pkg.go.dev/context (for timeouts)

### 8. Expected Deliverable at the end of Step 1

- Git Repo with a clean first commit titled `feat: MVP WebSocket chat with rooms and presence`
- README with a screenshot or GIF showing two tabs chatting
- No copy-pasted code (all typed and understood)

When you have checked the entire checklist → you move to step 2 (sharding + 8000 connections).

Good luck, tonight you are laying the foundations for everything else.
When you are finished: simply write **"step 1 finished"**.
