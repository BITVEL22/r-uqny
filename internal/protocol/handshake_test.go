package protocol

import (
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/BITVEL22/r-uqny/internal/identity"
)

func TestNewHandshake(t *testing.T) {
	id, err := identity.Generate()
	if err != nil {
		t.Fatalf("unexpected error generating identity: %v", err)
	}

	handshake, err := NewHandshake(id)
	if err != nil {
		t.Fatalf("unexpected error creating handshake: %v", err)
	}

	if err := handshake.Validate(); err != nil {
		t.Fatalf("unexpected handshake validation error: %v", err)
	}

	if handshake.NodeID != id.NodeID {
		t.Fatalf("expected NodeID %q, got %q", id.NodeID, handshake.NodeID)
	}

	if len(handshake.Challenge) != ChallengeSize {
		t.Fatalf("expected challenge size %d, got %d", ChallengeSize, len(handshake.Challenge))
	}
}

func TestHandshakeEncodeDecode(t *testing.T) {
	id, err := identity.Generate()
	if err != nil {
		t.Fatalf("unexpected error generating identity: %v", err)
	}

	original, err := NewHandshake(id)
	if err != nil {
		t.Fatalf("unexpected error creating handshake: %v", err)
	}

	data, err := original.Encode()
	if err != nil {
		t.Fatalf("unexpected error encoding handshake: %v", err)
	}

	decoded, err := DecodeHandshake(data)
	if err != nil {
		t.Fatalf("unexpected error decoding handshake: %v", err)
	}

	if decoded.NodeID != original.NodeID {
		t.Fatalf("expected NodeID %q, got %q", original.NodeID, decoded.NodeID)
	}

	if string(decoded.PublicKey) != string(original.PublicKey) {
		t.Fatal("decoded public key does not match original")
	}

	if string(decoded.Signature) != string(original.Signature) {
		t.Fatal("decoded signature does not match original")
	}

	if string(decoded.Challenge) != string(original.Challenge) {
		t.Fatal("decoded challenge does not match original")
	}
}

func TestHandshakeRejectsModifiedNodeID(t *testing.T) {
	id, err := identity.Generate()
	if err != nil {
		t.Fatalf("unexpected error generating identity: %v", err)
	}

	handshake, err := NewHandshake(id)
	if err != nil {
		t.Fatalf("unexpected error creating handshake: %v", err)
	}

	handshake.NodeID = "0000000000000000000000000000000000000000000000000000000000000000"

	if err := handshake.Validate(); err == nil {
		t.Fatal("expected modified NodeID to be rejected")
	}
}

func TestHandshakeRejectsModifiedSignature(t *testing.T) {
	id, err := identity.Generate()
	if err != nil {
		t.Fatalf("unexpected error generating identity: %v", err)
	}

	handshake, err := NewHandshake(id)
	if err != nil {
		t.Fatalf("unexpected error creating handshake: %v", err)
	}

	handshake.Signature[0] ^= 0xff

	if err := handshake.Validate(); err == nil {
		t.Fatal("expected modified signature to be rejected")
	}
}

func TestHandshakeRejectsInvalidChallenge(t *testing.T) {
	id, err := identity.Generate()
	if err != nil {
		t.Fatalf("unexpected error generating identity: %v", err)
	}

	handshake, err := NewHandshake(id)
	if err != nil {
		t.Fatalf("unexpected error creating handshake: %v", err)
	}

	// Tamper with challenge length
	handshake.Challenge = handshake.Challenge[:10]

	if err := handshake.Validate(); err == nil {
		t.Fatal("expected handshake with invalid challenge length to be rejected")
	}
}

func TestHandshakeReplayProtection(t *testing.T) {
	id, err := identity.Generate()
	if err != nil {
		t.Fatalf("unexpected error generating identity: %v", err)
	}

	handshake, err := NewHandshake(id)
	if err != nil {
		t.Fatalf("unexpected error creating handshake: %v", err)
	}

	// Set timestamp outside MaxHandshakeAge (e.g., 10 minutes ago)
	handshake.Timestamp = time.Now().UTC().Add(-10 * time.Minute)

	// Re-sign with expired timestamp
	data, err := handshake.signingData()
	if err != nil {
		t.Fatalf("failed to calculate signing data: %v", err)
	}

	handshake.Signature = ed25519.Sign(
		ed25519.PrivateKey(id.PrivateKey),
		data,
	)

	if err := handshake.Validate(); err == nil {
		t.Fatal("expected handshake with expired timestamp (replay attempt) to be rejected")
	}
}

