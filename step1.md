# 🟢 Step 1: The Functional Core
**Building the Engine Foundation**

| Status | Details |
| :--- | :--- |
| **Goal** | A robust, single-node WebSocket server |
| **Estimated Time** | 4 - 6 hours |
| **Key Concepts** | Goroutines, Channels, JSON Marshaling, Mutexes |

---

## 1. 🎯 Mission Brief

Your first mission is to build the **minimum viable engine** for Nebula.
Before thinking about distributed systems or 20k connections, we need a solid core that runs perfectly on a single machine.

**The Goal**: A user can open a browser, join a specific room (e.g., "general"), and chat with other users in that room in real-time.

---

## 2. ⚙️ Technical Specifications

### 2.1. The Constraints
*   **No Web Frameworks**: Use `net/http` for the server.
*   **WebSocket Library**: `github.com/gorilla/websocket` is authorized.
*   **Concurrency**:
    *   1 Goroutine per Client (Read Pump).
    *   1 Goroutine per Client (Write Pump).
    *   1 Goroutine for the Hub (Broadcast Loop).
*   **Protocol**: Strict JSON. No raw text.

### 2.2. The JSON Protocol
You must implement a structured protocol. Every message sent over the WebSocket must follow this schema:

```json
{
  "type": "message",       // "message", "join", "leave", "presence"
  "room": "general",       // Target room
  "sender": "Alice",       // Username
  "content": "Hello!",     // Actual text (optional for join/leave)
  "users": ["Alice", "Bob"] // Only for "presence" events
}
```

### 2.3. Expected Behavior
1.  **Connection**: Client connects to `/?room=general&username=Alice`.
2.  **Join**: Server adds Alice to "general" room and broadcasts a `join` event.
3.  **Presence**: Server sends the updated list of users in "general" to **everyone** in that room.
4.  **Chat**: Alice sends "Hello". Server broadcasts it to Bob and Charlie (in "general").
5.  **Leave**: Alice closes tab. Server detects it, removes her, and broadcasts `leave` + new `presence`.

---

## 3. 🏗️ Architecture Blueprint

Your project structure for Step 1 should look like this:

```bash
nebula/
├── cmd/nebula/
│   └── main.go           # Entry point: HTTP server + Graceful Shutdown
├── internal/
│   ├── server/           # HTTP Handlers (Upgrade to WS)
│   ├── hub/              # The Brain: Manages rooms & Broadcasts
│   ├── client/           # The Worker: ReadPump & WritePump
│   ├── room/             # Room state (Set of clients)
│   └── message/          # JSON Struct definitions
└── web/static/           # The provided HTML/JS client
```

---

## 4. 📝 Implementation Guide

### Phase A: The Message Package
Define your Go structs in `internal/message`. Use `json` tags.
*   *Tip*: Use `omitempty` to keep payloads light.

### Phase B: The Hub & Rooms
The Hub is the central traffic controller.
*   It needs a `Register` channel, an `Unregister` channel, and a `Broadcast` channel.
*   It must manage a map of Rooms: `map[string]*Room`.
*   **Critical**: Accessing maps is **not thread-safe**. You must use `sync.RWMutex` or confine access to a single "run loop" goroutine.

### Phase C: The Client Pumps
*   **ReadPump**: Reads from WS, unmarshals JSON, validates it, and sends it to the Hub.
*   **WritePump**: Listens on a Go channel, marshals to JSON, and writes to WS.
*   *Why separate pumps?* So that a blocked writer doesn't stop us from reading (and vice-versa).

### Phase D: Graceful Shutdown
In `main.go`, listen for `os.Interrupt`.
When received:
1.  Close the HTTP server.
2.  Close all Client connections cleanly.
3.  Exit.

---

## 5. ✅ Acceptance Criteria

To validate Step 1, you must check all these boxes:

- [ ] **Protocol Compliance**: The server sends/receives valid JSON, not plain text.
- [ ] **Room Isolation**: Users in "dev" cannot see messages from "general".
- [ ] **Presence System**: When I join, I see who is already there.
- [ ] **Real-time Updates**: When I leave, others see my name disappear immediately.
- [ ] **Concurrency Safety**: `go test -race ./...` returns NO errors.
- [ ] **Graceful Shutdown**: Ctrl+C stops the server without errors or zombie goroutines.

### 🚀 Ready for Step 2?
Once this checklist is complete, you are ready to tackle **Sharding & High Performance**.
