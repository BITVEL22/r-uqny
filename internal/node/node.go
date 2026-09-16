package node

import (
	"errors"

	"github.com/BITVEL22/r-uqny/internal/config"
	"github.com/BITVEL22/r-uqny/internal/identity"
	"github.com/BITVEL22/r-uqny/internal/transport"
)

var (
	ErrAlreadyRunning = errors.New("node is already running")
	ErrNotRunning     = errors.New("node is not running")
)

// Status represents the current lifecycle state of a node.
type Status string

const (
	StatusStopped Status = "stopped"
	StatusRunning Status = "running"
)

// Node represents a participant in the r/uqny network.
type Node struct {
	ID       string
	Status   Status
	Config   config.Config
	Identity identity.Identity
	Listener *transport.TCPListener
}

// New creates a new stopped node with the given configuration and identity.
func New(cfg config.Config, id identity.Identity) Node {
	return Node{
		ID:       id.NodeID,
		Status:   StatusStopped,
		Config:   cfg,
		Identity: id,
	}
}

// Start changes the node state from stopped to running.
func (n *Node) Start() error {
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
	if n.Status != StatusRunning {
		return nil, ErrNotRunning
	}

	conn, err := n.Listener.Accept()
	if err != nil {
		return nil, err
	}

	return transport.NewConnection(conn), nil
}
