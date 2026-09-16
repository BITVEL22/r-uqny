package protocol

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/BITVEL22/r-uqny/internal/identity"
)

const HandshakeVersion = 1

type Handshake struct {
	Version   int       `json:"version"`
	NodeID    string    `json:"node_id"`
	PublicKey []byte    `json:"public_key"`
	Timestamp time.Time `json:"timestamp"`
	Signature []byte    `json:"signature"`
}

func NewHandshake(id identity.Identity) (Handshake, error) {
	if err := id.Validate(); err != nil {
		return Handshake{}, err
	}

	handshake := Handshake{
		Version:   HandshakeVersion,
		NodeID:    id.NodeID,
		PublicKey: append([]byte(nil), id.PublicKey...),
		Timestamp: time.Now().UTC(),
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

func (h Handshake) Encode() ([]byte, error) {
	if err := h.Validate(); err != nil {
		return nil, err
	}

	return json.Marshal(h)
}

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

func (h Handshake) signingData() ([]byte, error) {
	unsigned := struct {
		Version   int       `json:"version"`
		NodeID    string    `json:"node_id"`
		PublicKey []byte    `json:"public_key"`
		Timestamp time.Time `json:"timestamp"`
	}{
		Version:   h.Version,
		NodeID:    h.NodeID,
		PublicKey: h.PublicKey,
		Timestamp: h.Timestamp,
	}

	return json.Marshal(unsigned)
}
