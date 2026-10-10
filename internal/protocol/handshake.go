package protocol

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/BITVEL22/r-uqny/internal/identity"
)

// HandshakeVersion defines the current supported protocol handshake version.
const HandshakeVersion = 1

// ChallengeSize specifies the required size in bytes for the random challenge nonce used in authentication.
const ChallengeSize = 32

// MaxHandshakeAge defines the maximum allowed time skew or age for a handshake timestamp to prevent replay attacks.
const MaxHandshakeAge = 5 * time.Minute

// Handshake represents an authenticated protocol handshake payload exchanged between nodes during session setup.
// It contains cryptographic identity proof, a fresh random challenge nonce for challenge-response authentication,
// and a timestamp for replay protection.
type Handshake struct {
	Version   int       `json:"version"`
	NodeID    string    `json:"node_id"`
	PublicKey []byte    `json:"public_key"`
	Timestamp time.Time `json:"timestamp"`
	Challenge []byte    `json:"challenge"`
	Signature []byte    `json:"signature"`
}

// NewHandshake creates and signs a new Handshake for the given identity.
// It generates a fresh random challenge nonce using crypto/rand to enforce challenge-response authentication
// and sets a UTC timestamp to provide replay protection.
func NewHandshake(id identity.Identity) (Handshake, error) {
	if err := id.Validate(); err != nil {
		return Handshake{}, err
	}

	challenge := make([]byte, ChallengeSize)
	if _, err := io.ReadFull(rand.Reader, challenge); err != nil {
		return Handshake{}, err
	}

	handshake := Handshake{
		Version:   HandshakeVersion,
		NodeID:    id.NodeID,
		PublicKey: append([]byte(nil), id.PublicKey...),
		Timestamp: time.Now().UTC(),
		Challenge: challenge,
	}

	data, err := handshake.signingData()
	if err != nil {
		return Handshake{}, err
	}

	handshake.Signature = ed25519.Sign(
		ed25519.PrivateKey(id.PrivateKey),
		data,
	)

	return handshake, nil
}

// Validate checks the structural validity, challenge size, timestamp age (replay protection),
// node ID key derivation, and Ed25519 cryptographic signature of the handshake.
func (h Handshake) Validate() error {
	if h.Version != HandshakeVersion {
		return errors.New("unsupported handshake version")
	}

	if len(h.PublicKey) != ed25519.PublicKeySize {
		return errors.New("invalid public key length")
	}

	if h.NodeID == "" {
		return errors.New("node ID cannot be empty")
	}

	if len(h.NodeID) != identity.NodeIDLength {
		return errors.New("invalid node ID length")
	}

	if h.Timestamp.IsZero() {
		return errors.New("timestamp cannot be zero")
	}

	// Replay protection: verify timestamp is within the acceptable time window.
	now := time.Now().UTC()
	age := now.Sub(h.Timestamp)
	if age < 0 {
		age = -age
	}
	if age > MaxHandshakeAge {
		return errors.New("handshake timestamp expired or in future (replay protection trigger)")
	}

	// Challenge-response authentication: verify challenge nonce size.
	if len(h.Challenge) != ChallengeSize {
		return errors.New("invalid challenge nonce length")
	}

	if len(h.Signature) != ed25519.SignatureSize {
		return errors.New("invalid signature length")
	}

	hash := sha256.Sum256(h.PublicKey)
	expectedNodeID := hex.EncodeToString(hash[:])

	if h.NodeID != expectedNodeID {
		return errors.New("node ID does not match public key")
	}

	data, err := h.signingData()
	if err != nil {
		return err
	}

	if !ed25519.Verify(
		ed25519.PublicKey(h.PublicKey),
		data,
		h.Signature,
	) {
		return errors.New("invalid handshake signature")
	}

	return nil
}

// Encode serializes the handshake to JSON format after validating its fields.
func (h Handshake) Encode() ([]byte, error) {
	if err := h.Validate(); err != nil {
		return nil, err
	}

	return json.Marshal(h)
}

// DecodeHandshake deserializes JSON data into a Handshake struct and validates it.
func DecodeHandshake(data []byte) (Handshake, error) {
	var handshake Handshake

	if err := json.Unmarshal(data, &handshake); err != nil {
		return Handshake{}, err
	}

	if err := handshake.Validate(); err != nil {
		return Handshake{}, err
	}

	return handshake, nil
}

// signingData creates a canonical JSON representation of the unsigned handshake fields,
// including the version, node ID, public key, timestamp, and challenge nonce.
func (h Handshake) signingData() ([]byte, error) {
	unsigned := struct {
		Version   int       `json:"version"`
		NodeID    string    `json:"node_id"`
		PublicKey []byte    `json:"public_key"`
		Timestamp time.Time `json:"timestamp"`
		Challenge []byte    `json:"challenge"`
	}{
		Version:   h.Version,
		NodeID:    h.NodeID,
		PublicKey: h.PublicKey,
		Timestamp: h.Timestamp,
		Challenge: h.Challenge,
	}

	return json.Marshal(unsigned)
}
