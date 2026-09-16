package identity

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
)

const NodeIDLength = 64

// Identity represents the cryptographic identity of a r/uqny node.
type Identity struct {
	NodeID     string `json:"node_id"`
	PublicKey  []byte `json:"public_key"`
	PrivateKey []byte `json:"private_key"`
}

// Generate creates a new Ed25519 identity.
func Generate() (Identity, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return Identity{}, err
	}

	return Identity{
		NodeID:     deriveNodeID(publicKey),
		PublicKey:  append([]byte(nil), publicKey...),
		PrivateKey: append([]byte(nil), privateKey...),
	}, nil
}

// Validate checks whether the identity is structurally valid.
func (i Identity) Validate() error {
	if len(i.PublicKey) != ed25519.PublicKeySize {
		return errors.New("invalid public key length")
	}

	if len(i.PrivateKey) != ed25519.PrivateKeySize {
		return errors.New("invalid private key length")
	}

	expectedNodeID := deriveNodeID(ed25519.PublicKey(i.PublicKey))
	if i.NodeID != expectedNodeID {
		return errors.New("node ID does not match public key")
	}

	return nil
}

// Save writes the identity to a JSON file.
func (i Identity) Save(path string) error {
	if err := i.Validate(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err
	}

	data = append(data, '\n')

	return os.WriteFile(path, data, 0600)
}

// Load reads an identity from a JSON file.
func Load(path string) (Identity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Identity{}, err
	}

	var identity Identity

	if err := json.Unmarshal(data, &identity); err != nil {
		return Identity{}, err
	}

	if err := identity.Validate(); err != nil {
		return Identity{}, err
	}

	return identity, nil
}

func deriveNodeID(publicKey ed25519.PublicKey) string {
	hash := sha256.Sum256(publicKey)
	return hex.EncodeToString(hash[:])
}
