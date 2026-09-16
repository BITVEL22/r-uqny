package node

import (
	"errors"
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
	ID     string
	Status Status
}

// New creates a new stopped node with the given ID.
func New(id string) Node {
	return Node{
		ID:     id,
		Status: StatusStopped,
	}
}

// Start changes the node state from stopped to running.
func (n *Node) Start() error {
	if n.Status == StatusRunning {
		return ErrAlreadyRunning
	}

	n.Status = StatusRunning
	return nil
}

// Stop changes the node state from running to stopped.
func (n *Node) Stop() error {
	if n.Status == StatusStopped {
		return ErrNotRunning
	}

	n.Status = StatusStopped
	return nil
}
