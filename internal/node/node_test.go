package node

import (
	"testing"

	"github.com/BITVEL22/r-uqny/internal/config"
	"github.com/BITVEL22/r-uqny/internal/identity"
	"github.com/BITVEL22/r-uqny/internal/transport"
)

func newTestNode(t *testing.T) Node {
	t.Helper()

	id, err := identity.Generate()
	if err != nil {
		t.Fatalf("unexpected error generating identity: %v", err)
	}

	return New(config.Config{
		NodeID:        id.NodeID,
		ListenAddress: "127.0.0.1:0",
		DataDirectory: "./data",
	}, id)
}

func TestNew(t *testing.T) {
	n := newTestNode(t)

	if n.ID != n.Identity.NodeID {
		t.Fatalf("expected node ID to match identity NodeID")
	}

	if n.Status != StatusStopped {
		t.Fatalf("expected initial status %q, got %q", StatusStopped, n.Status)
	}
}

func TestStart(t *testing.T) {
	n := newTestNode(t)

	if err := n.Start(); err != nil {
		t.Fatalf("unexpected error starting node: %v", err)
	}
	defer n.Stop()

	if n.Status != StatusRunning {
		t.Fatalf("expected status %q, got %q", StatusRunning, n.Status)
	}
}

func TestStop(t *testing.T) {
	n := newTestNode(t)

	if err := n.Start(); err != nil {
		t.Fatalf("unexpected error starting node: %v", err)
	}

	if err := n.Stop(); err != nil {
		t.Fatalf("unexpected error stopping node: %v", err)
	}

	if n.Status != StatusStopped {
		t.Fatalf("expected status %q, got %q", StatusStopped, n.Status)
	}
}

func TestStartAlreadyRunning(t *testing.T) {
	n := newTestNode(t)

	if err := n.Start(); err != nil {
		t.Fatalf("unexpected error starting node: %v", err)
	}
	defer n.Stop()

	if err := n.Start(); err != ErrAlreadyRunning {
		t.Fatalf("expected ErrAlreadyRunning, got %v", err)
	}
}

func TestStopAlreadyStopped(t *testing.T) {
	n := newTestNode(t)

	if err := n.Stop(); err != ErrNotRunning {
		t.Fatalf("expected ErrNotRunning, got %v", err)
	}
}

func TestStartCreatesListener(t *testing.T) {
	n := newTestNode(t)

	if err := n.Start(); err != nil {
		t.Fatalf("unexpected error starting node: %v", err)
	}
	defer n.Stop()

	if n.Status != StatusRunning {
		t.Fatalf("expected status %q, got %q", StatusRunning, n.Status)
	}

	if n.Listener == nil {
		t.Fatal("expected listener to be created")
	}
}

func TestStopClosesListener(t *testing.T) {
	n := newTestNode(t)

	if err := n.Start(); err != nil {
		t.Fatalf("unexpected error starting node: %v", err)
	}

	if n.Listener == nil {
		t.Fatal("expected listener to be created")
	}

	if err := n.Stop(); err != nil {
		t.Fatalf("unexpected error stopping node: %v", err)
	}

	if n.Status != StatusStopped {
		t.Fatalf("expected status %q, got %q", StatusStopped, n.Status)
	}

	if n.Listener != nil {
		t.Fatal("expected listener to be nil after stop")
	}
}

func TestAcceptConnection(t *testing.T) {
	n := newTestNode(t)

	if err := n.Start(); err != nil {
		t.Fatalf("unexpected error starting node: %v", err)
	}
	defer n.Stop()

	done := make(chan error, 1)

	go func() {
		conn, err := n.Accept()
		if err != nil {
			done <- err
			return
		}

		defer conn.Close()
		done <- nil
	}()

	rawConn, err := transport.DialTCP(n.Listener.Address().String())
	if err != nil {
		t.Fatalf("unexpected error dialing node: %v", err)
	}
	defer rawConn.Close()

	if err := <-done; err != nil {
		t.Fatalf("unexpected Accept() error: %v", err)
	}
}
