# r/uqny (reliable unified quality network)

**r/uqny** is an experimental, self-hosted and peer-to-peer communication project focused on resilient communication across different network transports.

The project aims to explore how messaging and voice communication can work without depending on third-party communication services, while remaining secure, portable, and adaptable to different connectivity conditions.

> **Status: Early development / Experimental**

## Project Goals

The long-term goal of r/uqny is to build a communication system that can operate across multiple transport mediums:

* Internet / IP networks
* Bluetooth
* Radio
* Potentially other transports in the future

The application layer should remain as independent as possible from the underlying transport.

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

This architecture is intended to allow r/uqny to eventually switch between available transports without requiring the communication layer to be rewritten.

---

# Roadmap

The project is developed incrementally. Early versions prioritize learning, correctness, and stability over feature count.

## v0.0 — Foundation

Establish the project and protocol foundations.

* [ ] Create initial protocol design
* [ ] Define peer and session concepts
* [ ] Define basic message format
* [ ] Implement basic logging
* [ ] Create initial test structure
* [ ] Establish development documentation
* [ ] Create a minimal CLI prototype

**Goal:** Two r/uqny nodes can establish a basic communication session.

---

## v0.1 — Peer-to-Peer Communication

Focus on basic communication without depending on third-party communication services.

### v0.1.0 — Messaging

* [ ] Peer discovery on LAN
* [ ] P2P connection establishment
* [ ] Peer identification
* [ ] Text messaging
* [ ] Message IDs
* [ ] Timestamps
* [ ] Basic acknowledgements
* [ ] Connection timeout handling

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

### v0.2.0 — Secure Identity

* [ ] Peer identity
* [ ] Public/private key pairs
* [ ] Key exchange
* [ ] Peer authentication
* [ ] Secure session establishment

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

Introduce a transport abstraction so the communication protocol is not tightly coupled to IP networking.

* [ ] Transport abstraction
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

> r/uqny can intelligently handle different network conditions.

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

Explore communication without relying on a central infrastructure.

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
├── docs/
├── src/
├── tests/
└── .github/
    └── workflows/
```

The repository structure will evolve as the implementation becomes more mature.

---

# Development Status

| Component             | Status     |
| --------------------- | ---------- |
| Project foundation    | 🟡 Planned |
| P2P messaging         | ⚪ Planned  |
| Voice communication   | ⚪ Planned  |
| Internet P2P          | ⚪ Planned  |
| End-to-end encryption | ⚪ Planned  |
| Transport abstraction | ⚪ Planned  |
| Bluetooth             | ⚪ Planned  |
| Radio                 | ⚪ Planned  |
| Mesh networking       | ⚪ Planned  |

---

# Disclaimer

r/uqny is an experimental networking and communication project created for educational, research, and development purposes.

Network behavior, connectivity, and compatibility may vary depending on the underlying network, operating system, hardware, and transport medium.

---

# License

This project is licensed under the MIT License. See [`LICENSE`](LICENSE) for details.
