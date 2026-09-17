package peer

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/BITVEL22/r-uqny/internal/identity"
	"github.com/BITVEL22/r-uqny/internal/session"
	"github.com/BITVEL22/r-uqny/internal/transport"
)

const testPeerID = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func TestNew(t *testing.T) {
	p, err := New(testPeerID, "127.0.0.1:9000")
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	if p.ID != testPeerID {
		t.Fatalf("unexpected peer ID: %s", p.ID)
	}

	if p.Address != "127.0.0.1:9000" {
		t.Fatalf("unexpected peer address: %s", p.Address)
	}

	if p.State != StateDisconnected {
		t.Fatalf("unexpected initial state: %s", p.State)
	}

	if !p.LastConnected.IsZero() {
		t.Fatal("LastConnected should initially be zero")
	}

	if !p.LastSeen.IsZero() {
		t.Fatal("LastSeen should initially be zero")
	}
}

func TestNewValidation(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		address string
	}{
		{
			name:    "empty ID",
			id:      "",
			address: "127.0.0.1:9000",
		},
		{
			name:    "invalid ID",
			id:      "short-id",
			address: "127.0.0.1:9000",
		},
		{
			name:    "empty address",
			id:      testPeerID,
			address: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(test.id, test.address)
			if err == nil {
				t.Fatal("New() should return an error")
			}
		})
	}
}

func TestSetState(t *testing.T) {
	p, err := New(testPeerID, "127.0.0.1:9000")
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	states := []State{
		StateConnecting,
		StateHandshaking,
		StateConnected,
		StateDisconnected,
		StateClosed,
	}

	for _, state := range states {
		if err := p.SetState(state); err != nil {
			t.Fatalf("SetState(%q) returned error: %v", state, err)
		}

		if p.State != state {
			t.Fatalf("expected state %q, got %q", state, p.State)
		}
	}
}

func TestSetStateRejectsInvalidState(t *testing.T) {
	p, err := New(testPeerID, "127.0.0.1:9000")
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	err = p.SetState(State("invalid"))
	if err == nil {
		t.Fatal("SetState() should reject invalid state")
	}
}

func TestMarkConnected(t *testing.T) {
	before := time.Now().UTC()

	p, err := New(testPeerID, "127.0.0.1:9000")
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	p.MarkConnected()

	after := time.Now().UTC()

	if p.State != StateConnected {
		t.Fatalf("expected connected state, got %s", p.State)
	}

	if p.LastConnected.Before(before) || p.LastConnected.After(after) {
		t.Fatal("LastConnected was not set correctly")
	}

	if p.LastSeen.Before(before) || p.LastSeen.After(after) {
		t.Fatal("LastSeen was not set correctly")
	}
}

func TestMarkSeen(t *testing.T) {
	p, err := New(testPeerID, "127.0.0.1:9000")
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	p.LastSeen = time.Now().UTC().Add(-time.Hour)
	oldLastSeen := p.LastSeen

	p.MarkSeen()

	if !p.LastSeen.After(oldLastSeen) {
		t.Fatal("MarkSeen() should update LastSeen")
	}
}

func TestPeerIDLength(t *testing.T) {
	if len(testPeerID) != 64 {
		t.Fatalf("test peer ID should contain 64 characters, got %d", len(testPeerID))
	}

	if strings.TrimSpace(testPeerID) != testPeerID {
		t.Fatal("test peer ID should not contain whitespace")
	}
}

func TestSessionManagement(t *testing.T) {
	p, err := New(testPeerID, "127.0.0.1:9000")
	if err != nil {
		t.Fatalf("unexpected error creating peer: %v", err)
	}

	localID, err := identity.Generate()
	if err != nil {
		t.Fatalf("unexpected error generating identity: %v", err)
	}

	rawConn1, rawConn2 := net.Pipe()
	defer rawConn1.Close()
	defer rawConn2.Close()

	conn := transport.NewConnection(rawConn1)

	s, err := session.New(localID, conn)
	if err != nil {
		t.Fatalf("unexpected error creating session: %v", err)
	}

	if p.HasSession() {
		t.Fatal("peer should not have a session initially")
	}

	if err := p.AttachSession(s); err != nil {
		t.Fatalf("unexpected AttachSession() error: %v", err)
	}

	if !p.HasSession() {
		t.Fatal("peer should have a session after attaching")
	}

	got, err := p.GetSession()
	if err != nil {
		t.Fatalf("unexpected GetSession() error: %v", err)
	}

	if got != s {
		t.Fatal("GetSession() returned a different session")
	}

	if err := p.AttachSession(s); err != ErrSessionAlreadyAttached {
		t.Fatalf("expected ErrSessionAlreadyAttached, got %v", err)
	}

	if err := p.DetachSession(); err != nil {
		t.Fatalf("unexpected DetachSession() error: %v", err)
	}

	if p.HasSession() {
		t.Fatal("peer should not have a session after detaching")
	}

	if _, err := p.GetSession(); err != ErrSessionNotAttached {
		t.Fatalf("expected ErrSessionNotAttached, got %v", err)
	}

	if err := p.DetachSession(); err != ErrSessionNotAttached {
		t.Fatalf("expected ErrSessionNotAttached, got %v", err)
	}
}

func TestAttachSessionRejectsNil(t *testing.T) {
	p, err := New(testPeerID, "127.0.0.1:9000")
	if err != nil {
		t.Fatalf("unexpected error creating peer: %v", err)
	}

	if err := p.AttachSession(nil); err == nil {
		t.Fatal("AttachSession() should reject nil session")
	}
}
