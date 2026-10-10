# Development Documentation

Welcome to the **r/uqny** development documentation. This document provides an overview of the system architecture, component design, security mechanisms, and guidelines for testing and extending the codebase.

---

## Architecture Overview

**r/uqny** is organized into modular internal Go packages under `internal/`:

```text
internal/
├── config/      - Configuration defaults, JSON configuration loading/saving
├── identity/    - Ed25519 cryptographic key generation, storage, and Node ID derivation
├── logging/     - Structured logging initialization (slog)
├── node/        - Core Node lifecycle management, TCP listener, peer registry, session server
├── peer/        - Peer state tracking and session binding
├── protocol/    - Message framing, envelope definition, JSON serialization, and Handshake authentication
├── session/     - Authenticated node-to-node sessions and concurrency-safe message transport
├── transport/   - Transport layer abstraction (TCP listener, dialer, framing)
└── version/     - Version definitions
```

---

## Handshake & Authentication Protocol

The handshake protocol operates over any reliable connection transport. It provides node identity verification, challenge-response authentication, and replay protection.

### Handshake Structure

```go
type Handshake struct {
    Version   int       `json:"version"`
    NodeID    string    `json:"node_id"`
    PublicKey []byte    `json:"public_key"`
    Timestamp time.Time `json:"timestamp"`
    Challenge []byte    `json:"challenge"`
    Signature []byte    `json:"signature"`
}
```

### Security Mechanisms

1. **Identity Verification**:
   - `NodeID` must equal `hex(sha256(PublicKey))`.
   - Ed25519 signature is verified over canonical JSON representation of `(Version, NodeID, PublicKey, Timestamp, Challenge)`.

2. **Challenge-Response Authentication**:
   - A 32-byte cryptographically secure random nonce (`Challenge`) is generated via `crypto/rand` for every handshake.
   - Handshakes with invalid or missing challenge nonces are rejected.

3. **Replay Protection**:
   - Every handshake carries a UTC timestamp (`Timestamp`).
   - Timestamps outside the acceptable window (`MaxHandshakeAge = 5 * time.Minute`) are rejected to prevent replay attacks.

---

## Session & Node Lifecycle

1. **Starting a Node**:
   - Call `node.New(cfg, identity)` to initialize a node instance.
   - Call `node.Start()` to bind the TCP listener.
   - Call `node.Serve()` to continuously accept incoming peer connections.

2. **Establishing Sessions**:
   - **Client Mode**: `node.Connect(address)` dials a remote endpoint and runs `HandshakeAsClient()`.
   - **Server Mode**: `node.AcceptSession()` accepts a connection and runs `HandshakeAsServer()`.

3. **Message Transmission**:
   - Call `session.Send(msg)` to send framed JSON protocol messages.
   - Call `session.Receive()` to read messages from the peer.

---

## Development Workflow & Testing

Ensure all checks pass before submitting changes:

```bash
# Run all package unit tests
go test -v ./...

# Run race detector
go test -race ./...

# Run static analysis
go vet ./...
```
