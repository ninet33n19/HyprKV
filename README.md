# HyprKV

HyprKV is a small Redis-inspired key-value server written in Go.

This project exists as a learning exercise: I am trying to implement Redis from scratch in Golang to understand the protocol, command handling, networking, and storage model at a lower level.

It is not a full Redis replacement. The current codebase focuses on a minimal TCP server, RESP parsing/encoding, and an in-memory store with basic expiration support.

## Current Milestone

Current milestone: RESP parsing plus basic string-command support over TCP.

Working in this milestone:

- RESP request decoding
- RESP response encoding
- TCP connection handling
- In-memory key/value storage
- Basic Redis-style string commands: `PING`, `ECHO`, `SET`, `GET`
- Minimal `COMMAND` handling for client compatibility probes

## Current Status

Implemented today:

- TCP server listening on `:7379`
- RESP decoding and encoding
- In-memory key-value storage
- Basic command routing
- Supported commands:
  - `PING`
  - `ECHO`
  - `SET`
  - `GET`
  - `COMMAND`
- TTL field exists in storage items and the storage layer can handle expiration durations

Not implemented yet:

- Full Redis command coverage
- Persistence
- Replication
- Transactions
- Pub/Sub
- Eviction policies
- Proper Redis-compatible config loading
- Full TTL option parsing in the command layer

## Project Structure

```text
cmd/server/main.go          Entrypoint that boots the storage and TCP server
internal/server/            Connection handling and command routing
internal/storage/           In-memory store implementation
internal/resp/              RESP encoder/decoder and tests
internal/config/            Simple config type
Taskfile.yml                Convenience task for running the server
load_test.sh                Basic netcat-based connection script
```

## How It Works

### 1. TCP server

The server accepts raw TCP connections and reads RESP messages from clients.

### 2. RESP protocol layer

The `internal/resp` package parses incoming RESP values and encodes responses back to the client. This is the foundation that makes the server speak a Redis-like protocol.

### 3. Command routing

Decoded RESP arrays are routed in `internal/server/handler.go`. The server currently understands a small subset of Redis-style commands.

### 4. Storage

The `internal/storage` package keeps data in memory using a `map[string]*Item` protected by a `sync.RWMutex`.

## Supported Commands

### `PING`

```text
PING
PING hello
```

Returns `PONG` or echoes the provided bulk string.

### `ECHO`

```text
ECHO hello
```

Returns the provided message.

### `SET`

```text
SET mykey hello
```

Stores a value in memory.

### `GET`

```text
GET mykey
```

Returns the stored value or a null bulk string when the key does not exist.

### `COMMAND`

Currently returns an empty RESP array to satisfy Redis clients that probe command support.

## Running The Server

### With Go

```bash
go run cmd/server/main.go
```

The server listens on:

```text
127.0.0.1:7379
```

### With Task

If you use `task`:

```bash
task run
```

## Testing

Run the RESP tests with:

```bash
CGO_ENABLED=0 GOCACHE=/tmp/go-build go test ./...
```

The `CGO_ENABLED=0` and `GOCACHE=/tmp/go-build` settings were needed in my sandboxed environment. In a normal local setup, plain `go test ./...` may be enough.

## Example With redis-cli

If `redis-cli` is available locally:

```bash
redis-cli -p 7379 PING
redis-cli -p 7379 SET name hyprkv
redis-cli -p 7379 GET name
redis-cli -p 7379 ECHO hello
```

## Notes And Limitations

- This project is intentionally incomplete and educational
- The implementation is only partially Redis-compatible
- The server currently reads commands into a fixed buffer and handles a minimal subset of cases
- Expiration exists in storage, but the command parser does not yet expose Redis-style expiry options such as `EX` or `PX`
- Expired keys are treated as missing on read, but are not actively cleaned up

## Why This Project

The goal is to learn Redis by building the core pieces manually in Go:

- socket handling
- RESP parsing
- command dispatch
- concurrency control
- in-memory data structures
- protocol-compatible server behavior

## Next Steps

Natural improvements for the project:

1. Add `SET EX`, `SET PX`, or `EXPIRE`
2. Improve RESP parsing for streaming/partial reads
3. Add more Redis commands
4. Implement active cleanup of expired keys
5. Add persistence with snapshots or append-only logging
6. Expand integration testing with real client interactions
