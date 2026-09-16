package node

import "testing"

func TestNew(t *testing.T) {
    n := New("test-node-01")

    if n.ID != "test-node-01" {
        t.Fatalf("expected node ID %q, got %q", "test-node-01", n.ID)
    }
}
