package protocol

import (
	"testing"

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
