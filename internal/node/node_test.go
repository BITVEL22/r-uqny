package node

import "testing"

func TestNew(t *testing.T) {
	n := New("test-node-01")

	if n.ID != "test-node-01" {
		t.Fatalf("expected node ID %q, got %q", "test-node-01", n.ID)
	}

	if n.Status != StatusStopped {
		t.Fatalf("expected initial status %q, got %q", StatusStopped, n.Status)
	}
}

func TestStart(t *testing.T) {
	n := New("test-node-01")

	if err := n.Start(); err != nil {
		t.Fatalf("unexpected error starting node: %v", err)
	}

	if n.Status != StatusRunning {
		t.Fatalf("expected status %q, got %q", StatusRunning, n.Status)
	}
}

func TestStop(t *testing.T) {
	n := New("test-node-01")

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
	n := New("test-node-01")

	if err := n.Start(); err != nil {
		t.Fatalf("unexpected error starting node: %v", err)
	}

	if err := n.Start(); err != ErrAlreadyRunning {
		t.Fatalf("expected ErrAlreadyRunning, got %v", err)
	}
}

func TestStopAlreadyStopped(t *testing.T) {
	n := New("test-node-01")

	if err := n.Stop(); err != ErrNotRunning {
		t.Fatalf("expected ErrNotRunning, got %v", err)
	}
}
