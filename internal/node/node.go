package node

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/BITVEL22/r-uqny/internal/config"
	"github.com/BITVEL22/r-uqny/internal/identity"
	"github.com/BITVEL22/r-uqny/internal/peer"
	"github.com/BITVEL22/r-uqny/internal/session"
	"github.com/BITVEL22/r-uqny/internal/transport"
)

var (
	ErrAlreadyRunning = errors.New("node is already running")
	ErrNotRunning     = errors.New("node is not running")
	ErrPeerExists     = errors.New("peer already exists")
	ErrPeerNotFound   = errors.New("peer not found")
)

// Status represents the current lifecycle state of a node.
type Status string

const (
	StatusStopped Status = "stopped"
	StatusRunning Status = "running"
)

// Node represents a participant in the r/uqny network.
type Node struct {
	mu       sync.RWMutex
	ID       string
	Status   Status
	Config   config.Config
	Identity identity.Identity
	Listener *transport.TCPListener
	peers    map[string]peer.Peer
}

// New creates a new stopped node with the given configuration and identity.
func New(cfg config.Config, id identity.Identity) Node {
	return Node{
		ID:       id.NodeID,
		Status:   StatusStopped,
		Config:   cfg,
		Identity: id,
		peers:    make(map[string]peer.Peer),
	}
}

// Start changes the node state from stopped to running.
func (n *Node) Start() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.Status == StatusRunning {
		return ErrAlreadyRunning
	}

	listener, err := transport.ListenTCP(n.Config.ListenAddress)
	if err != nil {
		return err
	}

	n.Listener = listener
	n.Status = StatusRunning

	return nil
}

// Stop changes the node state from running to stopped.
func (n *Node) Stop() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.Status == StatusStopped {
		return ErrNotRunning
	}

	if n.Listener != nil {
		if err := n.Listener.Close(); err != nil {
			return err
		}
	}

	n.Listener = nil
	n.Status = StatusStopped

	return nil
}

// Accept waits for and accepts an incoming TCP connection.
func (n *Node) Accept() (*transport.Connection, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if n.Status != StatusRunning {
		return nil, ErrNotRunning
	}

	conn, err := n.Listener.Accept()
	if err != nil {
		return nil, err
	}

	return transport.NewConnection(conn), nil
}

// AddPeer adds a peer to the node's peer list.
func (n *Node) AddPeer(p peer.Peer) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if _, exists := n.peers[p.ID]; exists {
		return ErrPeerExists
	}

	n.peers[p.ID] = p

	return nil
}

// GetPeer returns a peer by its ID.
func (n *Node) GetPeer(id string) (peer.Peer, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	p, exists := n.peers[id]
	if !exists {
		return peer.Peer{}, ErrPeerNotFound
	}

	return p, nil
}

// RemovePeer removes a peer from the node's peer list.
func (n *Node) RemovePeer(id string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if _, exists := n.peers[id]; !exists {
		return ErrPeerNotFound
	}

	delete(n.peers, id)

	return nil
}

// PeerCount returns the number of known peers.
func (n *Node) PeerCount() int {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return len(n.peers)
}

// Peers returns a snapshot of all known peers.
func (n *Node) Peers() []peer.Peer {
	n.mu.RLock()
	defer n.mu.RUnlock()

	result := make([]peer.Peer, 0, len(n.peers))
	for _, p := range n.peers {
		result = append(result, p)
	}

	return result
}

// AcceptSession accepts an incoming connection and performs a server-side handshake.
func (n *Node) AcceptSession() (*session.Session, error) {
	conn, err := n.Accept()
	if err != nil {
		return nil, err
	}

	sess, err := session.New(n.Identity, conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	if err := sess.HandshakeAsServer(); err != nil {
		sess.Close()
		return nil, err
	}

	return sess, nil
}

// Connect dials a remote node and performs a client-side handshake.
func (n *Node) Connect(address string) (*session.Session, error) {
	rawConn, err := transport.DialTCP(address)
	if err != nil {
		return nil, err
	}

	conn := transport.NewConnection(rawConn)

	sess, err := session.New(n.Identity, conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	if err := sess.HandshakeAsClient(); err != nil {
		sess.Close()
		return nil, err
	}

	return sess, nil
}

// Serve starts accepting connections in a loop.
// It blocks until the node is stopped or an unrecoverable error occurs.
func (n *Node) Serve() error {
	for {
		conn, err := n.Accept()
		if err != nil {
			n.mu.RLock()
			status := n.Status
			n.mu.RUnlock()

			if status == StatusStopped {
				return nil
			}

			slog.Error("failed to accept connection", "error", err)
			continue
		}

		go n.handleConnection(conn)
	}
}

func (n *Node) handleConnection(conn *transport.Connection) {
	remoteAddr := conn.RemoteAddr().String()

	sess, err := session.New(n.Identity, conn)
	if err != nil {
		slog.Error("failed to create session",
			"error", err,
			"remote_addr", remoteAddr,
		)
		conn.Close()
		return
	}

	if err := sess.HandshakeAsServer(); err != nil {
		slog.Error("handshake failed",
			"error", err,
			"remote_addr", remoteAddr,
		)
		sess.Close()
		return
	}

	remoteID := sess.RemoteID()

	p, err := peer.New(remoteID, remoteAddr)
	if err != nil {
		slog.Error("failed to create peer",
			"error", err,
			"remote_addr", remoteAddr,
		)
		sess.Close()
		return
	}

	if err := p.AttachSession(sess); err != nil {
		slog.Error("failed to attach session",
			"error", err,
			"remote_id", remoteID,
		)
		sess.Close()
		return
	}

	p.MarkConnected()

	if err := n.AddPeer(p); err != nil {
		slog.Error("failed to add peer",
			"error", err,
			"remote_id", remoteID,
		)
		sess.Close()
		return
	}

	slog.Info("peer connected",
		"remote_id", remoteID,
		"remote_addr", remoteAddr,
	)
}

