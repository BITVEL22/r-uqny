package identity

import (
	"crypto/ed25519"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerate(t *testing.T) {
	identity, err := Generate()
	if err != nil {
		t.Fatalf("failed to generate identity: %v", err)
	}

	if err := identity.Validate(); err != nil {
		t.Fatalf("generated identity is invalid: %v", err)
	}

	if len(identity.NodeID) != NodeIDLength {
		t.Fatalf(
			"expected node ID length %d, got %d",
			NodeIDLength,
			len(identity.NodeID),
		)
	}
}

func TestGeneratedIdentitiesAreDifferent(t *testing.T) {
	first, err := Generate()
	if err != nil {
		t.Fatalf("failed to generate first identity: %v", err)
	}

	second, err := Generate()
	if err != nil {
		t.Fatalf("failed to generate second identity: %v", err)
	}

	if first.NodeID == second.NodeID {
		t.Fatal("expected generated identities to have different node IDs")
	}
}

func TestKeySizes(t *testing.T) {
	identity, err := Generate()
	if err != nil {
		t.Fatalf("failed to generate identity: %v", err)
	}

	if len(identity.PublicKey) != ed25519.PublicKeySize {
		t.Fatalf("unexpected public key size: %d", len(identity.PublicKey))
	}

	if len(identity.PrivateKey) != ed25519.PrivateKeySize {
		t.Fatalf("unexpected private key size: %d", len(identity.PrivateKey))
	}
}

func TestNodeIDMatchesPublicKey(t *testing.T) {
	identity, err := Generate()
	if err != nil {
		t.Fatalf("failed to generate identity: %v", err)
	}

	originalID := identity.NodeID
	identity.NodeID = "invalid"

	if err := identity.Validate(); err == nil {
		t.Fatal("expected validation error for mismatched node ID")
	}

	identity.NodeID = originalID

	if err := identity.Validate(); err != nil {
		t.Fatalf("identity should be valid after restoring node ID: %v", err)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "identity.json")

	expected, err := Generate()
	if err != nil {
		t.Fatalf("failed to generate identity: %v", err)
	}

	if err := expected.Save(path); err != nil {
		t.Fatalf("failed to save identity: %v", err)
	}

	actual, err := Load(path)
	if err != nil {
		t.Fatalf("failed to load identity: %v", err)
	}

	if actual.NodeID != expected.NodeID {
		t.Fatalf("node ID changed after load")
	}

	if string(actual.PublicKey) != string(expected.PublicKey) {
		t.Fatalf("public key changed after load")
	}

	if string(actual.PrivateKey) != string(expected.PrivateKey) {
		t.Fatalf("private key changed after load")
	}
}

func TestSaveUsesRestrictedPermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "identity.json")

	identity, err := Generate()
	if err != nil {
		t.Fatalf("failed to generate identity: %v", err)
	}

	if err := identity.Save(path); err != nil {
		t.Fatalf("failed to save identity: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat identity file: %v", err)
	}

	if info.Mode().Perm() != 0600 {
		t.Fatalf(
			"expected identity file permissions 0600, got %o",
			info.Mode().Perm(),
		)
	}
}
