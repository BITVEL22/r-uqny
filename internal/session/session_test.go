package session

import (
	"net"
	"testing"

	"github.com/BITVEL22/r-uqny/internal/identity"
	"github.com/BITVEL22/r-uqny/internal/protocol"
	"github.com/BITVEL22/r-uqny/internal/transport"
)

func createTestConnection(t *testing.T) (*transport.Connection, *transport.Connection) {
	t.Helper()

	left, right := net.Pipe()

	return transport.NewConnection(left), transport.NewConnection(right)
}

func TestNewSession(t *testing.T) {
	id, err := identity.Generate()
	if err != nil {
		t.Fatalf("unexpected identity generation error: %v", err)
	}

	conn, peer := createTestConnection(t)
	defer peer.Close()

	s, err := New(id, conn)
	if err != nil {
		t.Fatalf("unexpected session creation error: %v", err)
	}

	if s.State() != StateNew {
		t.Fatalf("expected state %q, got %q", StateNew, s.State())
	}

	if s.RemoteID() != "" {
		t.Fatalf("expected empty remote ID, got %q", s.RemoteID())
	}
}

func TestSessionEstablish(t *testing.T) {
	localID, err := identity.Generate()
	if err != nil {
		t.Fatalf("unexpected identity generation error: %v", err)
	}

	remoteID, err := identity.Generate()
	if err != nil {
		t.Fatalf("unexpected identity generation error: %v", err)
	}

	conn, peer := createTestConnection(t)
	defer peer.Close()

	s, err := New(localID, conn)
	if err != nil {
		t.Fatalf("unexpected session creation error: %v", err)
	}

	if err := s.Establish(remoteID.NodeID); err != nil {
		t.Fatalf("unexpected establish error: %v", err)
	}

	if s.State() != StateEstablished {
		t.Fatalf("expected state %q, got %q", StateEstablished, s.State())
	}

	if s.RemoteID() != remoteID.NodeID {
		t.Fatalf("expected remote ID %q, got %q", remoteID.NodeID, s.RemoteID())
	}

	if err := s.Establish(remoteID.NodeID); err != ErrAlreadyEstablished {
		t.Fatalf("expected ErrAlreadyEstablished, got %v", err)
	}
}

func TestSessionSendReceive(t *testing.T) {
	localID, err := identity.Generate()
	if err != nil {
		t.Fatalf("unexpected identity generation error: %v", err)
	}

	remoteID, err := identity.Generate()
	if err != nil {
		t.Fatalf("unexpected identity generation error: %v", err)
	}

	conn, peer := createTestConnection(t)
	defer peer.Close()

	s, err := New(localID, conn)
	if err != nil {
		t.Fatalf("unexpected session creation error: %v", err)
	}

	if err := s.Establish(remoteID.NodeID); err != nil {
		t.Fatalf("unexpected establish error: %v", err)
	}

	message := protocol.NewTextMessage(
		localID.NodeID,
		"test-message-1",
		"hello from session",
	)

	receiveDone := make(chan error, 1)

	go func() {
		received, err := peer.Receive()
		if err != nil {
			receiveDone <- err
			return
		}

		if received.MessageID != message.MessageID {
			receiveDone <- errUnexpectedMessageID
			return
		}

		if received.Payload != message.Payload {
			receiveDone <- errUnexpectedPayload
			return
		}

		receiveDone <- nil
	}()

	if err := s.Send(message); err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}

	if err := <-receiveDone; err != nil {
		t.Fatalf("unexpected receive error: %v", err)
	}
}

func TestSessionRejectsSendBeforeEstablish(t *testing.T) {
	localID, err := identity.Generate()
	if err != nil {
		t.Fatalf("unexpected identity generation error: %v", err)
	}

	conn, peer := createTestConnection(t)
	defer peer.Close()

	s, err := New(localID, conn)
	if err != nil {
		t.Fatalf("unexpected session creation error: %v", err)
	}

	message := protocol.NewTextMessage(
		localID.NodeID,
		"test-message-2",
		"should fail",
	)

	if err := s.Send(message); err != ErrNotEstablished {
		t.Fatalf("expected ErrNotEstablished, got %v", err)
	}
}

func TestSessionClose(t *testing.T) {
	localID, err := identity.Generate()
	if err != nil {
		t.Fatalf("unexpected identity generation error: %v", err)
	}

	conn, peer := createTestConnection(t)
	defer peer.Close()

	s, err := New(localID, conn)
	if err != nil {
		t.Fatalf("unexpected session creation error: %v", err)
	}

	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}

	if s.State() != StateClosed {
		t.Fatalf("expected state %q, got %q", StateClosed, s.State())
	}

	if err := s.Close(); err != nil {
		t.Fatalf("expected repeated close to succeed, got %v", err)
	}
}
