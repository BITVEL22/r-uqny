package node

import (
	"errors"
	"sync"

	"github.com/BITVEL22/r-uqny/internal/config"
	"github.com/BITVEL22/r-uqny/internal/identity"
	"github.com/BITVEL22/r-uqny/internal/peer"
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
