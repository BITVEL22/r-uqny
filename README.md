# r/uqny (reliable/unified quality network)

**r/uqny** is an experimental, self-hosted, peer-to-peer communication project focused on resilient communication across different network transports.

The project explores how messaging and voice communication can work without depending on third-party communication services, while remaining secure, portable, and adaptable to different connectivity conditions.

> **Status: Early development / Experimental**

---

# Project Goals

The long-term goal of r/uqny is to build a communication system that can operate across multiple transport mediums:

* Internet / IP networks
* Bluetooth
* Radio
* Potentially other transports in the future

The application and communication layers should remain as independent as possible from the underlying transport.

Conceptually:

```text
                    r/uqny
                       │
              Communication Layer
                       │
                Session / Protocol
                       │
                  Security Layer
                       │
                Transport Layer
                       │
          ┌────────────┼────────────┐
          │            │            │
        Internet    Bluetooth      Radio
```

This architecture is intended to allow r/uqny to eventually use different transports without requiring the communication layer to be fundamentally rewritten.

---

# Roadmap

The project is developed incrementally. Early versions prioritize learning, correctness, interoperability, and stability over feature count.

## v0.0 — Foundation

Establish the core project, node, protocol, identity, transport, peer, and session foundations.

### Completed

* [x] Create initial protocol design
* [x] Define basic message format
* [x] Create initial test structure
* [x] Create a minimal CLI prototype
* [x] Implement node lifecycle
* [x] Implement configuration system
* [x] Implement cryptographic node identity
* [x] Implement Ed25519 node identity keys
* [x] Derive node IDs from public keys
* [x] Implement protocol message envelope
* [x] Implement JSON message encoding/decoding
* [x] Implement TCP listener
* [x] Implement TCP dialer
* [x] Implement length-prefixed TCP message framing
* [x] Implement transport connection abstraction
* [x] Implement authenticated handshake foundation
* [x] Sign handshake data using node identity
* [x] Verify handshake signatures and node identity
* [x] Validate remote handshakes before establishing sessions
* [x] Define peer and session concepts
* [x] Implement peer management
* [x] Implement session abstraction
* [x] Implement session state management
* [x] Integrate authenticated handshake with sessions
* [x] Connect peers with sessions

### Remaining

* [ ] Implement basic logging
* [ ] Establish development documentation
* [ ] Implement node-level connection/session management
* [ ] Define connection lifecycle after handshake
* [ ] Implement connection timeout handling
* [ ] Implement connection error handling
* [ ] Implement challenge-response authentication
* [ ] Implement replay protection

**Current goal:**

> Two r/uqny nodes can establish a basic authenticated communication session.

> **Note:** The current handshake is a foundation for peer authentication. It validates the remote handshake signature and node identity, but it is not yet the final production security protocol and does not yet provide full replay protection or connection-specific challenge/response.

---

## v0.1 — Peer-to-Peer Communication

Focus on basic communication without depending on third-party communication services.

### v0.1.0 — Messaging

* [ ] Peer discovery on LAN
* [x] P2P connection establishment foundation
* [x] Peer identification
* [x] Text messaging foundation
* [x] Message IDs
* [x] Timestamps
* [ ] Basic acknowledgements
* [ ] Connection timeout handling
* [ ] Connection error handling
* [ ] Peer connection state management

### v0.1.1 — Voice

* [ ] Microphone input
* [ ] Audio playback
* [ ] Audio encoding/decoding
* [ ] Voice packetization
* [ ] Basic jitter handling
* [ ] Packet-loss handling
* [ ] Call establishment
* [ ] Call termination
* [ ] Mute functionality

### v0.1.2 — Internet Connectivity

* [ ] Internet P2P connectivity
* [ ] NAT traversal research
* [ ] Connection establishment across different networks
* [ ] Reconnection
* [ ] Network interruption handling
* [ ] Connection quality monitoring

### v0.1.x — Stabilization

* [ ] Crash handling
* [ ] Invalid packet handling
* [ ] Network interruption recovery
* [ ] Audio quality improvements
* [ ] Stress testing
* [ ] Automated tests
* [ ] Documentation improvements
* [ ] Compatibility testing

**v0.1 goal:**

> Reliable peer-to-peer messaging and voice communication.

---

## v0.2 — Security

Focus on making communication private and authenticated.

The cryptographic node identity and initial authenticated handshake introduced during v0.0 provide the foundation for this phase.

### v0.2.0 — Secure Identity

* [x] Cryptographic node identity foundation
* [x] Public/private key pairs
* [x] Node identity verification foundation
* [x] Authenticated handshake foundation
* [ ] Key exchange
* [ ] Secure session establishment
* [ ] Session key management

### v0.2.1 — End-to-End Encrypted Messaging

* [ ] E2EE text messages
* [ ] Message integrity protection
* [ ] Replay protection
* [ ] Secure key handling
* [ ] Key rotation

### v0.2.2 — Encrypted Voice

* [ ] Encrypted voice transport
* [ ] Secure session keys
* [ ] Authentication
* [ ] Integrity protection
* [ ] Security testing

**v0.2 goal:**

> Messages and voice communication are protected by end-to-end encryption.

r/uqny will use established cryptographic primitives and protocols rather than implementing cryptographic algorithms from scratch.

---

## v0.3 — Resilient Networking

Introduce a stronger transport abstraction so the communication protocol is not tightly coupled to IP networking.

The current TCP implementation and `Connection` abstraction are early foundations for this phase.

* [x] Initial transport abstraction
* [x] TCP transport implementation
* [x] Connection abstraction
* [ ] Transport interface
* [ ] Connection monitoring
* [ ] Automatic transport selection
* [ ] Automatic reconnection
* [ ] Network quality detection
* [ ] Multi-path communication research

Potential transports:

```text
Internet
   │
   ├── Wi-Fi
   ├── Ethernet
   └── Cellular
```

Future transports:

```text
Bluetooth
Radio
Other experimental transports
```

**v0.3 goal:**

> r/uqny can intelligently handle different network conditions and transport types.

---

## v0.4 — Bluetooth

Add Bluetooth as an alternative transport.

* [ ] Bluetooth peer discovery
* [ ] Bluetooth connection
* [ ] Messaging over Bluetooth
* [ ] Encrypted communication over Bluetooth
* [ ] Voice communication research
* [ ] Automatic fallback

**v0.4 goal:**

> Communication can continue through Bluetooth when IP connectivity is unavailable.

---

## v0.5 — Radio

Research and implement radio-based communication.

Because radio can have significantly lower bandwidth and higher packet loss than Internet connections, this phase will focus on adapting the protocol rather than simply replacing the transport.

### Initial goals

* [ ] Radio transport abstraction
* [ ] Peer discovery
* [ ] Short messages
* [ ] Packet fragmentation
* [ ] Packet acknowledgement
* [ ] Error recovery
* [ ] Store-and-forward research

### Experimental goals

* [ ] Low-bitrate voice research
* [ ] Adaptive codecs
* [ ] Voice prioritization

**v0.5 goal:**

> Establish a practical foundation for communication over constrained radio links.

---

## v0.6 — Decentralized Communication

Explore communication without relying on centralized infrastructure.

Potential research areas:

* [ ] Peer-to-peer discovery
* [ ] Peer relay
* [ ] Distributed routing
* [ ] Store-and-forward messaging
* [ ] Mesh networking
* [ ] Offline communication

This phase is intentionally exploratory and may change significantly based on earlier experiments.

---

## v1.0 — Stable Release

The first major stable release.

Potential requirements:

* [ ] Stable communication protocol
* [ ] Reliable messaging
* [ ] Reliable voice communication
* [ ] Secure communication
* [ ] Internet P2P
* [ ] Alternative transport support
* [ ] Comprehensive documentation
* [ ] Installation instructions
* [ ] Protocol specification
* [ ] Security documentation
* [ ] Automated testing
* [ ] Release binaries

The exact requirements for v1.0 will be determined as the project develops.

---

# Design Principles

r/uqny is developed around several principles:

### 1. Peer-to-peer first

Communication should not depend on a third-party communication provider.

### 2. Transport independence

The communication protocol should not assume that the underlying connection is always IP-based.

### 3. Security by design

Security should be considered during protocol design rather than added as an afterthought.

### 4. Use established cryptography

r/uqny should not implement cryptographic primitives from scratch when established, audited solutions are available.

### 5. Graceful degradation

The system should be able to adapt to poor connectivity and constrained transports.

### 6. Open development

Architecture, experiments, limitations, and design decisions should be documented publicly.

### 7. Incremental development

Each version should produce a working and testable improvement instead of attempting to build the entire system at once.

---

# Project Structure

```text
r-uqny/
├── README.md
├── LICENSE
├── .gitignore
├── go.mod
├── cmd/
│   └── r-uqny/
└── internal/
    ├── config/
    ├── identity/
    ├── node/
    ├── peer/
    ├── protocol/
    ├── session/
    ├── transport/
    └── version/
```

The repository structure will evolve as the implementation becomes more mature.

---

# Development Status

| Component                          | Status         |
| ---------------------------------- | -------------- |
| Project foundation                 | 🟢 Complete    |
| Node lifecycle                     | 🟢 Complete    |
| Configuration                      | 🟢 Complete    |
| Node identity                      | 🟢 Complete    |
| Protocol envelope                  | 🟢 Complete    |
| TCP transport                      | 🟢 Complete    |
| Connection abstraction             | 🟢 Complete    |
| Authenticated handshake foundation | 🟢 Complete    |
| Peer abstraction                   | 🟢 Complete    |
| Peer management                    | 🟢 Complete    |
| Session abstraction                | 🟢 Complete    |
| Session state management           | 🟢 Complete    |
| Node-to-node session               | 🟡 In progress |
| P2P messaging                      | 🟡 In progress |
| Voice communication                | ⚪ Planned      |
| Internet P2P                       | ⚪ Planned      |
| End-to-end encryption              | ⚪ Planned      |
| Transport abstraction              | 🟡 Foundation  |
| Bluetooth                          | ⚪ Planned      |
| Radio                              | ⚪ Planned      |
| Mesh networking                    | ⚪ Planned      |

### Status Legend

* 🟢 **Complete** — implemented and tested
* 🟡 **In progress** — actively being developed
* ⚪ **Planned** — not yet implemented

---

# Development

r/uqny is currently being developed in Go.

Basic checks used during development:

```bash
go test ./...
go vet ./...
go test -race ./...
```

The project currently prioritizes small, independently testable components before integrating them into a complete communication system.

---

# Disclaimer

r/uqny is an experimental networking and communication project created for educational, research, and development purposes.

Network behavior, connectivity, and compatibility may vary depending on the underlying network, operating system, hardware, and transport medium.

Security-related components should be considered experimental until they have undergone appropriate review and testing.

---

# License

This project is licensed under the MIT License. See [`LICENSE`](LICENSE) for details.
````
