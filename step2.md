# 🚀 Step 2: High Performance & Sharding
**Breaking the Monolith**

| Status | Details |
| :--- | :--- |
| **Goal** | Scale to 5,000+ concurrent connections |
| **Estimated Time** | 4 - 6 hours |
| **Key Concepts** | Sharding, Hashing, sync.Pool, Lock Contention |

---

## 1. 🎯 Mission Brief

**The Reality Check:**
Your Step 1 benchmark crashed at **~100 concurrent users**.
Why? Because every single message forced the entire server to stop (Global Mutex) while it iterated over clients.

**The Solution:**
We will **Shard** the Hub. Instead of one giant traffic controller, we will have 32 (or 64) smaller, independent controllers running in parallel.
If Shard #1 is busy broadcasting, Shard #2 is still free to accept new clients.

---

## 2. ⚙️ Technical Specifications

### 2.1. The Sharding Strategy
*   **Concept**: Divide the `rooms` map into `N` smaller maps (Shards).
*   **Routing**: A room is assigned to a specific shard using a hash function.
    *   `ShardID = hash(RoomName) % NumShards`
*   **Independence**: Each Shard has its own `sync.RWMutex`.
    *   *Result*: Lock contention is divided by `N`.

### 2.2. Memory Optimization (Zero-Copy)
At 10,000 messages/sec, the Garbage Collector (GC) becomes your enemy.
*   **sync.Pool**: Reuse `BroadcastMessage` structs and byte buffers instead of allocating new ones for every message.
*   **Lazy Marshaling**: Don't call `json.Marshal` inside the loop for every client. Marshal **once** per message, then send the bytes.

---

## 3. 🏗️ Architecture Blueprint

**The New Hub Structure:**

```go
type Hub struct {
    shards []*Shard
    // ...
}

type Shard struct {
    mu    sync.RWMutex
    rooms map[string]*Room
    // ...
}
```

**The Flow:**
1.  **Client Connects** (`?room=general`):
    *   `hash("general") % 32` -> Returns Shard #5.
    *   Hub delegates registration to Shard #5.
2.  **Broadcast**:
    *   Message for "general" arrives.
    *   Hub routes it to Shard #5.
    *   Shard #5 locks **only itself**. Shard #1..#4 and #6..#32 remain unlocked.

---

## 4. 📝 Implementation Guide

### Phase A: The Hasher
Implement a fast string hashing function (like FNV-1a or DJB2). It needs to be deterministic (same room name = same hash).

### Phase B: Refactoring the Hub
1.  Create the `Shard` struct (it looks a lot like your old Hub).
2.  Update `Hub` to initialize `32` shards.
3.  Update `Register`, `Unregister`, and `Broadcast` to route calls to the correct shard.

### Phase C: Optimizing the Broadcast
1.  **Pre-calculate JSON**: In `BroadcastMessage`, marshal the JSON **before** entering the client loop.
2.  **Non-blocking Send**: Ensure `client.send` channel never blocks the shard (use `select default` or buffered channels).

---

## 5. ✅ Acceptance Criteria

To validate Step 2, you must beat your previous benchmark score:

- [ ] **Sharding**: The Hub is split into at least 32 shards.
- [ ] **Scalability**: k6 benchmark runs stable with **2,000 concurrent users** (20x improvement).
- [ ] **Throughput**: Handle **50,000+ messages/sec**.
- [ ] **Latency**: p95 latency remains under **200ms** at peak load.
- [ ] **Stability**: No `websocket: close sent` errors during the test.

### 🚀 Ready?
This step is where you turn a toy server into a production-grade engine.
**Good luck.**
